package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		decode([]byte("l5:hello5:worldee"))
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi!")
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.RawQuery))
	})

	db := connectServer() //Database
	defer db.Close()

	torrentD := make(chan string)
	go getTorrentNyaa(db, 1204164, torrentD)
	d := <-torrentD
	fmt.Println(decode([]byte(d)))

	log.Fatal(http.ListenAndServe(":9000", nil))
}
