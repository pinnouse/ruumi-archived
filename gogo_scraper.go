package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const BASE_URL = "https://gogoanime.movie"

func fetchSrcFromEmbed(url string, src chan string) {
	resp, err := http.Get(url)
	if err != nil {
		src <- ""
	}
	b := resp.Body
	defer b.Close()
	allBytes, err := ioutil.ReadAll(b)
	if err != nil {
		panic(err)
	}
	srcExtracter := regexp.MustCompile(`sources:\[\{file: '([\w\:\/\.\?\=\&\-\_\,]+)`)
	matches := srcExtracter.FindStringSubmatch(string(allBytes))
	if len(matches) > 0 {
		src <- matches[1]
		return
	}
	src <- ""
}

func gogoFetchEpisode(categoryURL string, episode int, episodeSrc chan string) {
	url := fmt.Sprintf("%s%s-episode-%d", BASE_URL, categoryURL[len("/category"):], episode)
	resp, err := http.Get(url)
	if err != nil {
		episodeSrc <- ""
	}
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			episodeSrc <- ""
		case html.StartTagToken:
			token := z.Token()
			if token.Data == "a" {
				attr := getAttr(token, "data-video")
				if attr != "" {
					if attr[:2] == "//" {
						attr = "https:" + attr
					}
					go fetchSrcFromEmbed(attr, episodeSrc)
					<-episodeSrc
					return
				}
			}
		}
	}
}

func gogoFetchCategory(categoryURL string) (episodes int, err error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", BASE_URL, categoryURL))
	if err != nil {
		return
	}
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	inEpisodes := false
	lastEpisode := 0
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "ul" && checkAttr(t, "id", "episode_page") {
				inEpisodes = true
			}

			if inEpisodes {
				switch {
				case t.Data == "a":
					lastEpisode, err = strconv.Atoi(getAttr(t, "ep_end"))
					if lastEpisode > episodes {
						episodes = lastEpisode
					}
				}
			}
		case tt == html.EndTagToken:
			if z.Token().Data == "ul" && inEpisodes {
				return
			}
		}
	}
}

func gogoSearch(searchTerm string, page int, categories chan []GOGOCategory) {
	queries := url.Values{}
	queries.Set("keyword", searchTerm)
	queries.Set("page", strconv.Itoa(page))
	resp, err := http.Get(fmt.Sprintf("%s/search.html?%s", BASE_URL, queries.Encode()))
	cats := []GOGOCategory{}
	if err != nil {
		categories <- cats
		return
	}

	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	inItems := false
	added := false
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			categories <- cats
			return
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "ul" && checkAttr(t, "class", "items") {
				inItems = true
			}

			if inItems && t.Data == "a" && !added {
				cats = append(cats, GOGOCategory{
					Name:     getAttr(t, "title"),
					CatURL:   getAttr(t, "href"),
					Episodes: 0,
				})
				added = true
			}
		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "ul" && inItems {
				categories <- cats
				return
			} else if t.Data == "li" {
				added = false
			}
		}
	}
}
