package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	//Database setup first
	client := connectDB()

	r := mux.NewRouter()

	r.HandleFunc("/gogo", func(w http.ResponseWriter, r *http.Request) {
		gogoSearchHandler(w, r, client)
	})
	r.HandleFunc("/gogoCategory", func(w http.ResponseWriter, r *http.Request) {
		gogoCategoryHandler(w, r, client)
	})
	r.HandleFunc("/gogoEpisode", func(w http.ResponseWriter, r *http.Request) {
		gogoEpisodeHandler(w, r, client)
	})

	rooms := make(map[string]Room)

	r.HandleFunc("/addRoom", func(w http.ResponseWriter, r *http.Request) {
		addRoom(rooms, w, r)
		log.Println(rooms)
	})

	r.HandleFunc("/room/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		getRoom(rooms, w, r)
	})

	r.HandleFunc("/delRoom/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		delRoom(rooms, w, r)
	})

	r.HandleFunc("/chat/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		serveWs(rooms[params["roomId"]].hub, w, r)
	})

	/* Download a Bencoded file
	torrentD := make(chan []byte)
	go getTorrentNyaa(db, 1204164, "Boruto I think", torrentD)
	d := <-torrentD
	fmt.Println(decode(d))*/

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
