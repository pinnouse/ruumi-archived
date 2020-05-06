package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func searchHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "Resources not found, specify a query.", http.StatusNotAcceptable)
		return
	}
	anime, err := search(client, query)
	if err != nil {
		http.Error(w, "Anime not found, check logs for details.", http.StatusNotFound)
		return
	}
	js, err := json.Marshal(anime)
	if err != nil {
		http.Error(w, "Error parsing animes.", http.StatusInternalServerError)
		return
	}
	if len(js) == 0 {
		http.Error(w, "No results found.", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func listHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	amount, err := strconv.Atoi(r.URL.Query().Get("a"))
	if err != nil {
		http.Error(w, "Random amount no configured correctly.", http.StatusNotFound)
		return
	}
	random, err := getList(client, amount)
	if err != nil {
		http.Error(w, "Could not fetch random anime.", http.StatusNotFound)
		return
	}
	js, err := json.Marshal(random)
	if err != nil {
		http.Error(w, "Error parsing animes.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func animeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	anime, err := getAnime(client, r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Anime not found, check logs for details.", http.StatusNotFound)
		return
	}
	js, err := json.Marshal(anime)
	if err != nil {
		http.Error(w, "Error parsing anime.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func episodeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client, svc *s3.S3) {
	epNum, err := strconv.Atoi(r.URL.Query().Get("ep"))
	if err != nil {
		http.Error(w, "Incorrect episode number.", http.StatusNotFound)
		return
	}
	anime, err := getAnime(client, r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Anime not found, check logs for details.", http.StatusNotFound)
		return
	}
	key := fmt.Sprintf(anime.Key, epNum)
	log.Info("Got key", key)
	url, err := getObject(svc, key)
	if err != nil {
		http.Error(w, "Error retrieving the URL.", http.StatusInternalServerError)
		log.Error("Error retrieving url: ", err)
		return
	}
	if len(url) == 0 {
		http.Error(w, "No results found.", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"url\": \"%s\"}", url)
}

func addAnimeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body data.", http.StatusNotAcceptable)
		return
	}
	var anime Anime
	err = json.Unmarshal(body, &anime)
	if err != nil {
		http.Error(w, "Could not parse the anime.", http.StatusInternalServerError)
		return
	}
	err = addAnime(client, anime)
	if err != nil {
		http.Error(w, "Failed to add anime to DB.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, "{\"success\": true}")
}
