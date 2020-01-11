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

//ROOMS
type Message struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Sent    uint32 `json:"sent"`
	Author  User   `json:"author"`
}

type Video struct {
	Source string `json:"source"`
	Seek   uint32 `json:"seek"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	Id       string    `json:"id"`
	Owner    User      `json:"owner"`
	Video    Video     `json:"video"`
	Users    []User    `json:"users"`
	Messages []Message `json:"messages"`
	hub      *Hub
}
