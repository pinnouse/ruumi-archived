package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/pinnouse/ruumi/storage"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi!")
	})

	db := connect_server()

	log.Fatal(http.ListenAndServe(":9000", nil))
}
