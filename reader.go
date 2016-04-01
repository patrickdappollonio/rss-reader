package rssreader

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Article struct {
	Title        string
	Content      string
	Published    time.Time
	URL          *url.URL
	PreviewImage *url.URL
}

var (
	ErrNoURLProvided     = errors.New("rssreader: no url provided")
	ErrMalformedURL      = errors.New("rssreader: malformed url")
	ErrCantConnect       = errors.New("rssreader: can't connect to the given url")
	ErrCantParseResponse = errors.New("rssreader: can't parse xml format")
)

type rss struct {
	config Config
}

type Config struct {
	RSSURL   string
	MaxItems int
}

func Setup(c Config) *rss {
	return &rss{config: c}
}

func (r *rss) ReadFeed() ([]*Article, error) {
	// Shorthand
	web := r.config.RSSURL

	// Check first if there is an URL incoming
	if strings.TrimSpace(web) == "" {
		return nil, ErrNoURLProvided
	}

	// Try parse the URL
	if _, err := url.Parse(web); err != nil {
		return nil, ErrMalformedURL
	}

	// Fetch the RSS content
	resp, err := http.Get(web)

	// Check if possible to connect
	if err != nil {
		return nil, ErrCantConnect
	}

	// Close body once done
	defer resp.Body.Close()

	// Create a struct on-the-fly
	var xmlresponse struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			Items []struct {
				Title       string `xml:"title"`
				Description string `xml:"description"`
				Link        string `xml:"link"`
				PubDate     string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	// Try decoding the response
	if err := xml.NewDecoder(resp.Body).Decode(&xmlresponse); err != nil {
		return nil, ErrCantParseResponse
	}

	// Create an empty list
	var list []*Article

	// Iterate over the items
	for _, v := range xmlresponse.Channel.Items {
		if len(list) < r.config.MaxItems {
			list = append(list, &Article{
				Title:   v.Title,
				Content: v.Description,
			})
		}
	}

	return list, nil
}
