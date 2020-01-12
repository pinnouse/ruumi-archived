package main

//GOGO
type GOGOEpisode struct {
	EpNum    int    `json:"epNum"`
	SrcURL   string `json:"srcURL"`
	Category string `json:"category"`
}

type GOGOCategory struct {
	Name     string `json:"name"`
	CatURL   string `json:"catURL"`
	Episodes int    `json:"episodes"`
}
