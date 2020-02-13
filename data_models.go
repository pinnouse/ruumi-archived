package main

//Anime
type Episode struct {
	EpNum uint8  `json:"epNum"`
	Key   string `json:"key"`
}

type Anime struct {
	Id        string    `json:"id" bson:"_id"`
	Title     string    `json:"title"`
	AltTitles []string  `json:"altTitles"`
	Poster    string    `json:"poster"`
	Episodes  []Episode `json:"episodes"`
}
