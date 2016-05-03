// An RSS reader that, for the moment, only works with Tumblr blogs. It also fetches
// all the images within the content and pick the most appropiate one based on a minimum
// width or height.
//
// As a recommendation, it's better to use this package with caching enabled, since if there's
// any RSS with several images, you'll ended up using a bit of bandwidth.
//
// Usage:
//
//	 getter := Setup(Config{
//	 	RSSURL:        "http://your.blog.here/rss",
//	 	MaxItems:      3,
//	 	MinImageWidth: 200,
//	 	UseCache:      true,
//	 })
//
//	 articles, err := getter.ReadFeed()
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
	ErrNotInitialized    = errors.New("rssreader: not initialized")
	ErrNeededAbsURL      = errors.New("rssreader: relative url provided")
	ErrMalformedURL      = errors.New("rssreader: malformed url")
	ErrCantConnect       = errors.New("rssreader: can't connect to the given url")
	ErrCantParseResponse = errors.New("rssreader: can't parse xml format")
)

// RSS is a struct that holds the RSS reader config
// plus the caching configuration
type RSS struct {
	config   Config
	cache    *cash.Cash
	cachekey string
}

// Config is the struct that holds the Configuration,
// such as the RSS URL we'll fetch, the number of items
// we want to retrieve from the feed, the minimum image
// width or height and wether we use cache or not.
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

// ReadFeed performs the RSS reading and processing.
// The function will check first wether the process
// will use a cache, if so, it'll check if the articles
// are already in cache, and if not, it'll fetch them and
// store them for later use.
func (r *RSS) ReadFeed() ([]*Article, error) {
	if r == nil {
		return nil, ErrNotInitialized
	}

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

	// URL holder
	parsedURL, err := url.ParseRequestURI(web)

	// Try parse the URL
	if err != nil {
		return nil, ErrMalformedURL
	}

	// Check if the URL was absolute
	if !parsedURL.IsAbs() {
		return nil, ErrNeededAbsURL
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
