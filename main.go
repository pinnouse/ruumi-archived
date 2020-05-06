package main

import (
	"context"
	"fmt"
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
			log.Println(fmt.Sprintf("Some bad user is doing something: %s", r.RemoteAddr))
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
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/anime", func(w http.ResponseWriter, r *http.Request) {
		animeHandler(w, r, client)
	}).Methods("GET")
	r.HandleFunc("/episode", func(w http.ResponseWriter, r *http.Request) {
		episodeHandler(w, r, client, svc)
	}).Methods("GET")
	r.HandleFunc("/addAnime", func(w http.ResponseWriter, r *http.Request) {
		addAnimeHandler(w, r, client)
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
