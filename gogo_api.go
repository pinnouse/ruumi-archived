package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func gogoSearchHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}
	q := values.Get("q")
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 1
	}

	categories := dbGetSearch(client, strings.TrimSpace(strings.ToLower(q)), page)

	if err != nil {
		fmt.Fprintf(w, "There was an error searching")
		return
	}

	jsonCat, err := json.Marshal(categories.Results)
	if err != nil {
		fmt.Fprintf(w, "Internal error, whoops!")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonCat))
	for _, cat := range categories.Results {
		dbSetCategory(client, cat)
	}
}

func gogoCategoryHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}

	catURL := values.Get("category")
	cat, err := dbGetCategory(client, catURL)
	numEpisodes := 0
	if cat.Episodes == 0 {
		numEpisodes, err = gogoFetchCategory(catURL)
		cat.Episodes = numEpisodes
		dbSetCategory(client, cat)
	} else {
		numEpisodes = cat.Episodes
	}

	if err != nil {
		fmt.Fprintf(w, "Error fetching category")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"episodes\":%d}", numEpisodes)
}

func gogoEpisodeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}
	catURL := values.Get("category")
	epNum, err := strconv.Atoi(values.Get("episode"))
	if err != nil {
		epNum = 1
	}
	ep, err := dbGetEpisode(client, catURL, epNum)
	if err != nil {
		fmt.Fprintf(w, "Internal error whoops: %s", err)
		return
	}
	epJSON, err := json.Marshal(ep)
	if err != nil {
		fmt.Fprintf(w, "Internal error, my b: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(epJSON))
}
