package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const URL = "https://nyaa.si"

func request(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		//TODO: Log this to error and don't panic
		panic(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Println("HTML:\n\n", string(bytes))

	return bytes
}

func getAttr(t html.Token, attr string) (val string) {
	for _, a := range t.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}

func checkAttr(t html.Token, attrName string, attrVal string) (has bool) {
	for _, a := range t.Attr {
		if a.Key == attrName {
			if a.Val == attrVal {
				return true
			}
			return false
		}
	}
	return false
}

func nyaaSearch(query string, page string) (torrents []Torrent, err error) {
	extractNumRE, _ := regexp.Compile("[0-9]+")
	//Initial scrape
	qValues := url.Values{}
	qValues.Set("c", "1_2")
	if len(query) > 0 {
		qValues.Set("q", query)
	}
	if len(page) > 0 {
		qValues.Set("p", page)
	}
	nyaaUrl := fmt.Sprintf("%s?%s", URL, qValues.Encode())
	resp, err := http.Get(nyaaUrl)
	nyaaResp := []Torrent{}
	if err != nil {
		fmt.Printf("Failed to crawl page %s", nyaaUrl)
		return nyaaResp, err
	}
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	row := false
	context := 0
	/*
		Valid contexts:
		0: NONE
		1: SET TORRENT
	*/
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nyaaResp, nil
		case tt == html.StartTagToken:
			t := z.Token()
			switch {
			case t.Data == "a":
				if !row || context == 0 {
					continue
				}
				nyaaResp = append(nyaaResp, Torrent{
					extractNumRE.FindString(getAttr(t, "href")),
					getAttr(t, "title"),
				})
				context = 0
			case t.Data == "tr":
				row = true
			case t.Data == "td":
				if checkAttr(t, "colspan", "2") {
					context = 1
				}
			}
		case tt == html.EndTagToken:
			if !row || z.Token().Data != "tr" {
				continue
			}
			row = false
			context = 0
		}
	}
	return nyaaResp, nil
}
