package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Link struct {
	Href string `xml:"href,attr"`
}

type RSSFeed struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Author      struct {
		Name string `xml:"name"`
	} `xml:"author"`
	Links []Link     `xml:"link"`
	Entry []RSSEntry `xml:"entry"`
}

type RSSEntry struct {
	Title     string `xml:"title"`
	Link      Link   `xml:"link"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	ID        string `xml:"id"`
	Summary   string `xml:"summary"`
	Content   string `xml:"content"`
}

func urlToFeed(url string) (RSSFeed, error) {

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	fmt.Println("Fetching feed", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}
