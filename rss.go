package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	XMLName  xml.Name `xml:"rss"`
	Channel  Channel  `xml:"channel"`
	Version  string   `xml:"version,attr"`
	DC       string   `xml:"xmlns:dc,attr"`
	Content  string   `xml:"xmlns:content,attr"`
	Atom     string   `xml:"xmlns:atom,attr"`
	Hashnode string   `xml:"xmlns:hashnode,attr"`
}

type Channel struct {
	Title         string     `xml:"title"`
	Description   string     `xml:"description"`
	Link          string     `xml:"link"`
	Generator     string     `xml:"generator"`
	LastBuildDate string     `xml:"lastBuildDate"`
	AtomLinks     []AtomLink `xml:"atom:link"`
	Language      string     `xml:"language"`
	TTL           int        `xml:"ttl"`
	Items         []Item     `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr,omitempty"`
}

type Item struct {
	Title          string   `xml:"title"`
	Description    string   `xml:"description"`
	Link           string   `xml:"link"`
	GUID           GUID     `xml:"guid"`
	Categories     []string `xml:"category"`
	Creator        string   `xml:"dc:creator"`
	PubDate        string   `xml:"pubDate"`
	ContentEncoded string   `xml:"content:encoded"`
	CoverImage     string   `xml:"hashnode:coverImage"`
}

type GUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

func urlToFeedURL(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

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
