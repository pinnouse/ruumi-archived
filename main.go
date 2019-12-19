package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprintf(w, "There was an error with the search query: %s", err)
		return
	}
	searchQuery := values.Get("q")
	pageQuery := values.Get("p")
	pgq := 1
	if len(searchQuery) > 0 {
		if len(pageQuery) == 0 {
			pageQuery = "1"
		}

		pgq, err = strconv.Atoi(pageQuery)
		if err != nil || pgq == 0 {
			fmt.Println(err)
			pgq = 1
		} else if pgq < 0 {
			pgq = -pgq
		}
		endS, err := nyaaSearch(searchQuery, strconv.Itoa(pgq))
		if err != nil {
			fmt.Fprintf(w, "There was an error with query to Nyaa: %s", err)
			return
		}

		resJson, err := json.Marshal(endS)
		if err != nil {
			fmt.Fprintf(w, "Failed to marshall JSON: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, string(resJson))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Please provide a search term.")
	}
}

func main() {
	//Database setup first
	db := connectServer()
	defer db.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/nyaaSearch", searchHandler)

	http.HandleFunc("/nyaaAPI", func(w http.ResponseWriter, r *http.Request) {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			fmt.Fprintf(w, "Couldn't understand your query: %s", err)
			return
		}

		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			fmt.Fprintf(w, "ID not formatted properly, make sure it's an int: %s", err)
			return
		}

		torrentBytes := make(chan []byte)
		go getTorrentNyaa(db, int32(id), values.Get("name"), torrentBytes)
		bytes := <-torrentBytes
		jsonTorrent, err := json.Marshal(decode(bytes))
		if err != nil {
			fmt.Fprintf(w, "Issue parsing the torrent: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, string(jsonTorrent))
		//TODO: Convert the base64 encoded strings to UTF-8
		bt, err := decodeBT(bytes)
		jsonTorrent, err = json.Marshal(bt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jsonTorrent)
		fmt.Println(string(jsonTorrent))
	})

	/* Download a Bencoded file
	torrentD := make(chan []byte)
	go getTorrentNyaa(db, 1204164, "Boruto I think", torrentD)
	d := <-torrentD
	fmt.Println(decode(d))*/

	log.Fatal(http.ListenAndServe(":9000", nil))
}
