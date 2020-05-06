package main

//Anime
type Anime struct {
	Id        string   `json:"id" bson:"_id"`
	Title     string   `json:"title"`
	AltTitles []string `json:"altTitles"`
	Poster    string   `json:"poster"`
	Episodes  uint16   `json:"episodes"`
	Key       string   `json:"key"`
}
