package main

import "time"

//GOGO
type GOGOEpisode struct {
	EpID   string `json:"epID"`
	SrcURL string `json:"srcURL"`
}

type GOGOCategory struct {
	Name     string        `json:"name"`
	CatURL   string        `json:"catURL""`
	Episodes []GOGOEpisode `json:"episodes"`
}

type GOGOCategoryD struct {
	ID       string `db:"gid"`
	Name     string `db:"name"`
	CatURL   string `db:"caturl"`
	Episodes string `db:"episodes"`
}

//BITTORRENT
type BTFile struct {
	Length uint32
	Path   [][]byte
}

type BTInfo struct {
	Length      uint32
	Files       []BTFile
	Name        []byte
	PieceLength uint32
	Pieces      []byte
}

type BitTorrent struct {
	Announce     []byte     `mapstructure:"announce"`
	AnnounceList [][][]byte `mapstructure:"announce-list"`
	Comment      []byte     `mapstructure:"comment"`
	CreatedBy    []byte     `mapstructure:"created by"`
	CreationDate uint32     `mapstructure:"creation date"`
	Encoding     []byte     `mapstructure:"encoding"`
	Info         BTInfo     `mapstructure:"info"`
}

type NyaaTorrent struct {
	Id   string
	Name string
}

//ROOMS
type Message struct {
	Id      string
	Content string
	Sent    time.Time
	Author  User
}

type Video struct {
	Torrent BitTorrent
	Seek    time.Time
}

type User struct {
	Id    string
	Name  string
	Email string
}

type Room struct {
	Id       string
	Owner    User
	Video    Video
	Seek     uint32
	Users    []User
	Messages []Message
}

//ENDPOINT STRUCTS
type ENDsearch struct {
	Torrents []NyaaTorrent
	LastPage bool
}
