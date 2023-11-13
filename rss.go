package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type ChannelItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Enclosure   string `xml:"enclosure"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

type RssChannel struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description string        `xml:"description"`
	Language    string        `xml:"language"`
	Copyright   string        `xml:"copyright"`
	PubDate     string        `xml:"pubDate"`
	WebMaster   string        `xml:"webMaster"`
	Items       []ChannelItem `xml:"item"`
}

type RSSFeed struct {
	Channels []RssChannel `xml:"channel"`
}

func UrlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	var rss RSSFeed
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return RSSFeed{}, err
	}

	return rss, err
}
