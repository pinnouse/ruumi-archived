package main

//Anime
type Episode struct {
	EpNum  uint8  `json:"epNum"`
	Source string `json:"source"`
}

type Anime struct {
	Id       int32     `json:"id"`
	Title    string    `json:"title"`
	Episodes []Episode `json:"episodes"`
}

//User
type User struct {
	Id       string `json:"id"` //Technically a uint64
	Name     string `json:"name"`
	UserType uint8  `json:"userType"` //0 is no premium
}
