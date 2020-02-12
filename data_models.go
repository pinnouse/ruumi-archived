package main

//Anime
type Episode struct {
	EpNum uint8  `json:"epNum"`
	Key   string `json:"key"`
}

type Anime struct {
	Id        int32     `json:"id"`
	Title     string    `json:"title"`
	AltTitles []string  `json:"altTitles"`
	Poster    string    `json:"poster"`
	Episodes  []Episode `json:"episodes"`
}

//User
type User struct {
	Id       string `json:"id"` //Technically a uint64
	Name     string `json:"name"`
	UserType uint8  `json:"userType"` //0 is no premium
}
