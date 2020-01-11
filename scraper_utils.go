package main

import "golang.org/x/net/html"

func getAttr(token html.Token, attr string) string {
	for _, a := range token.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return ""
}

func checkAttr(token html.Token, attr string, val string) bool {
	return getAttr(token, attr) == val
}
