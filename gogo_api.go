package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/url"
	"strconv"
)

func gogoSearchHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}
	keyword := values.Get("keyword")
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 1
	}
	categories, err := gogoSearch(keyword, page)
	if err != nil {
		fmt.Fprintf(w, "There was an error searching")
		return
	}

	jsonCat, err := json.Marshal(categories)
	if err != nil {
		fmt.Fprintf(w, "Internal error, whoops!")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for _, cat := range categories {
		addCategory(db, cat.Name, cat.CatURL)
	}
	fmt.Fprintf(w, string(jsonCat))
}

func gogoCategoryHandler(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}
	catURL := values.Get("category")
	numEpisodes, err := gogoFetchCategory(catURL)
	if err != nil {
		fmt.Fprintf(w, "Error fetching category")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "{\"episodes\":%d}", numEpisodes)
}

func gogoEpisodeHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "Couldn't understand your query: %s", err)
		return
	}
	catName := values.Get("category")
	epNum, err := strconv.Atoi(values.Get("episode"))
	if err != nil {
		epNum = 1
	}
	episode := make(chan GOGOEpisode)
	go getEpisode(db, catName, epNum, episode)
	ep := <-episode
	epJSON, err := json.Marshal(ep)
	if err != nil {
		fmt.Fprintf(w, "Internal error, my b: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, string(epJSON))
}
