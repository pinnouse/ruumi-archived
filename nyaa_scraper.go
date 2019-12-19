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

func request(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
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

func nyaaSearch(query string, page string) (nyaaResponse ENDsearch, err error) {
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
	nyaaResp := ENDsearch{
		[]NyaaTorrent{},
		true,
	}
	if err != nil {
		fmt.Printf("Failed to crawl page %s\n", nyaaUrl)
		return nyaaResp, err
	}
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)
	var context uint8 = 0
	/*  Context
	OFFSET VALUE
	0: ROW
	1: VALID COLUMN
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
				if context&3 < 3 || getAttr(t, "class") == "comments" {
					continue
				}
				nyaaResp.Torrents = append(nyaaResp.Torrents, NyaaTorrent{
					extractNumRE.FindString(getAttr(t, "href")),
					getAttr(t, "title"),
				})
				context = 0
			case t.Data == "li":
				if checkAttr(t, "class", "next") {
					nyaaResp.LastPage = false
				}
			case t.Data == "tr":
				context = 1
			case t.Data == "td":
				if checkAttr(t, "colspan", "2") {
					context = context | 2
				}
			}
		case tt == html.EndTagToken:
			if context&1 == 0 || z.Token().Data != "tr" {
				continue
			}
			context = 0
		}
	}
	return nyaaResp, nil
}
