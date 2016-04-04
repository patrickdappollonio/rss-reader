package rssreader

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/daryl/cash"
	"github.com/kennygrant/sanitize"
	"github.com/patrickdappollonio/image-extractor"
)

// Article is a single article which contains all the
// information of an RSS article, as well as the content
// in HTML and in raw (with HTML tags stripped)
type Article struct {
	Title        string
	ContentHTML  template.HTML
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

type RSS struct {
	config   Config
	cache    *cash.Cash
	cachekey string
}

type Config struct {
	RSSURL         string
	MaxItems       int
	MinImageWidth  int
	MinImageHeight int
	UseCache       bool
}

const defaultExpiration = 24 * time.Hour

// Setup creates a single instance of an RSS reader with
// a specific configuration given.
func Setup(c Config) *RSS {
	var cache *cash.Cash
	var cachekey string

	if c.UseCache {
		cache = cash.New(cash.Conf{
			defaultExpiration,
			defaultExpiration,
		})

		cachekey = fmt.Sprintf(
			"%v-%v-%v-%v",
			c.RSSURL,
			c.MaxItems,
			c.MinImageWidth,
			c.MinImageHeight,
		)
	}

	return &RSS{
		config:   c,
		cache:    cache,
		cachekey: cachekey,
	}
}

func (r *RSS) ReadFeed() ([]*Article, error) {
	// Check if the info is in cache already
	if r.cachekey != "" {
		if data, ok := r.cache.Get(r.cachekey); ok {
			return data.([]*Article), nil
		}
	}

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
	var selectedOne, articleURL *url.URL
	var parsedDate time.Time

	// Iterate over the items
	for _, v := range xmlresponse.Channel.Items {
		if len(list) < r.config.MaxItems {
			// Select the best image
			selectedOne = selectBestImage(
				r.config.MinImageWidth,
				r.config.MinImageHeight,
				extractor.ImageExtractor{Content: v.Description}.GetAll(),
			)

			// Parse date, if is not correctly parsed, a
			// zero-value is set
			parsedDate, _ = time.Parse(time.RFC1123Z, v.PubDate)

			// Parse the article URL, trusting the address is right
			articleURL, _ = url.Parse(v.Link)

			// Append to the list
			list = append(list, &Article{
				Title:        v.Title,
				ContentHTML:  template.HTML(v.Description),
				Content:      strings.TrimSpace(sanitize.HTML(v.Description)),
				PreviewImage: selectedOne,
				URL:          articleURL,
				Published:    parsedDate,
			})
		}
	}

	// Save to cache
	if r.cachekey != "" {
		r.cache.Set(r.cachekey, list, cash.Default)
	}

	return list, nil
}
