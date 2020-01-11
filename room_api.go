package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
)

func addRoom(rooms map[string]Room, w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		room := Room{}
		err = json.Unmarshal(body, &room)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing json", http.StatusInternalServerError)
		}
		room.Id = uuid.Must(uuid.NewV4()).String()[:8]
		room.hub = newHub()
		go room.hub.run()
		rooms[room.Id] = room
		rjs, err := json.Marshal(room)
		if err != nil {
			http.Error(w, "Error converting room to json", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, string(rjs))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getRoom(rooms map[string]Room, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if room, ok := rooms[params["roomId"]]; ok {
		rjs, err := json.Marshal(room)
		if err != nil {
			http.Error(w, "Error converting room to json", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, string(rjs))
	} else {
		http.Error(w, "Room not found", http.StatusNotFound)
	}
}

func delRoom(rooms map[string]Room, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if room, ok := rooms[params["roomId"]]; ok {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
			}
			user := User{}
			err = json.Unmarshal(body, &user)
			isOwner := cmp.Equal(room.Owner, user)
			if !isOwner {
				http.Error(w, "You are not the owner of this room!", http.StatusForbidden)
			} else {
				room.hub.close()
				delete(rooms, room.Id)
			}
		} else {
			http.Error(w, "Wrong request type", http.StatusMethodNotAllowed)
		}
	}
}
