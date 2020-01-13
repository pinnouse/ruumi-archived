package main

//GOGO
type GOGOEpisode struct {
	EpNum       int    `json:"epNum"`
	SrcURL      string `json:"srcURL"`
	Category    string `json:"category"`
	LastUpdated int64  `json:"lastUpdated"`
}

type GOGOCategory struct {
	Name     string `json:"name"`
	CatURL   string `json:"catURL"`
	Episodes int    `json:"episodes"`
}

type GOGOSearchResults struct {
	SearchTerm  string         `json:"searchTerm"`
	Results     []GOGOCategory `json:"results"`
	LastUpdated int64          `json:"lastUpdated"`
}
