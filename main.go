package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("AUTH") == "" || r.Header.Get("Authorization") == os.Getenv("AUTH") {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not authorized.", http.StatusForbidden)
		}
	})
}

func main() {
	//Database setup first
	client := connectDB()
	defer client.Disconnect(context.Background())
	svc := newSession()

	r := mux.NewRouter()

	r.Use(authMiddleware)

	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/anime", func(w http.ResponseWriter, r *http.Request) {
		animeHandler(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/episode", func(w http.ResponseWriter, r *http.Request) {
		episodeHandler(w, r, client, svc)
	}).Methods("GET")
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			userHandler(w, r, client)
		} else if r.Method == "POST" {
			addUserHandler(w, r, client)
		} else {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		}
	}).Methods("GET", "POST")
	r.HandleFunc("/addAnime", func(w http.ResponseWriter, r *http.Request) {
		addAnimeHandler(w, r, client)
	}).Methods("POST")
	r.HandleFunc("/addEpisode", func(w http.ResponseWriter, r *http.Request) {
		addEpisodeHandler(w, r, client)
	}).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Server starting on 8000")
	log.Fatal(srv.ListenAndServe())
}
