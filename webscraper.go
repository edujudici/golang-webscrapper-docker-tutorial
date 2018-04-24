package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func hasClass(attribs []html.Attribute, className string) bool {
	for _, attr := range attribs {
		if attr.Key == "class" && strings.Contains(attr.Val, className) {
			return true
		}
	}
	return false
}

func getFirstTextNode(htmlParsed *html.Node) *html.Node {
	for m := htmlParsed.FirstChild; m != nil; m = m.NextSibling {
		if m.Type == html.TextNode {
			return m
		}
		r := getFirstTextNode(m)
		if r != nil {
			return r
		}
	}
	return nil
}

func getFirstElementByClass(htmlParsed *html.Node, elm, className string) *html.Node {
	for m := htmlParsed.FirstChild; m != nil; m = m.NextSibling {
		if m.Data == elm && hasClass(m.Attr, className) {
			return m
		}
		r := getFirstElementByClass(m, elm, className)
		if r != nil {
			return r
		}
	}
	return nil
}

func scrap(url string) (r Result) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: It can't scrap '", url, "'")
	}
	// Close body when function ends
	defer resp.Body.Close()
	body := resp.Body
	htmlParsed, err := html.Parse(body)
	if err != nil {
		fmt.Println("ERROR: It can't parse html '", url, "'")
	}
	header := getFirstElementByClass(htmlParsed, "header", "")

	a := getFirstElementByClass(header, "a", "ds-link--styleSubtle")
	r.userName = getFirstTextNode(a).Data

	div := getFirstElementByClass(htmlParsed, "div", "section-content")
	h1 := getFirstElementByClass(div, "h1", "graf--title")
	r.title = getFirstTextNode(h1).Data

	footer := getFirstElementByClass(htmlParsed, "footer", "u-paddingTop10")
	buttonLikes := getFirstElementByClass(footer, "button", "js-multirecommendCountButton")
	r.likes = getFirstTextNode(buttonLikes).Data

	return
}
