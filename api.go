package main

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func animeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	animeId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Incorrect ID format.", http.StatusNotAcceptable)
		return
	}
	anime, err := getAnime(client, int32(animeId))
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

func episodeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	animeId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Incorrect ID format.", http.StatusNotAcceptable)
		return
	}
	epNum, err := strconv.Atoi(r.URL.Query().Get("ep"))
	if err != nil {
		http.Error(w, "Incorrect episode number.", http.StatusNotFound)
		return
	}
	anime, err := getAnime(client, int32(animeId))
	if err != nil {
		http.Error(w, "Anime not found, check logs for details.", http.StatusNotFound)
		return
	}
	if epNum > len(anime.Episodes) || epNum < 1 {
		http.Error(w, "That episode could not be found.", http.StatusNotFound)
		return
	}
	js, err := json.Marshal(anime.Episodes[epNum])
	if err != nil {
		http.Error(w, "Error parsing the anime.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func userHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	userId := r.URL.Query().Get("id")
	_, err := strconv.ParseUint(userId, 10, 64)
	if len(userId) == 0 || err != nil {
		http.Error(w, "Incorrect ID format.", http.StatusNotAcceptable)
		return
	}
	user, err := getUser(client, userId)
	if err != nil {
		http.Error(w, "User not found, check logs for details.", http.StatusNotFound)
		return
	}
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error parsing the user.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, string(js))
}

func addUserHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body data.", http.StatusNotAcceptable)
		return
	}
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Could not parse the user.", http.StatusInternalServerError)
		return
	}
	err = addUser(client, user)
	if err != nil {
		http.Error(w, "Failed to add user to DB.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, "{success: true}")
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
	io.WriteString(w, "{success: true}")
}

func addEpisodeHandler(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body data.", http.StatusNotAcceptable)
		return
	}
	var response struct {
		animeId int32
		episode Episode
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		http.Error(w, "Could not parse the episode.", http.StatusInternalServerError)
		return
	}
	err = addEpisode(client, response.animeId, response.episode)
	if err != nil {
		http.Error(w, "Failed to add episode to DB.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	io.WriteString(w, "{success: true}")
}
