package extractor

import (
	"net/url"
	"regexp"
)

var reg = regexp.MustCompile(`https?://[^/\s]+/\S+\.(jpg|png|gif)`)

// ImageExtractor is the struct who contains the string of text
// where the extractor is going to run. It is a struct because
// if so, you can just create a single instance and call it any
// times you need.
type ImageExtractor struct {
	Content string
}

// GetAll gets a slice of all images found in text. Returns an empty
// slice (not nil) if no image was found.
func (d ImageExtractor) GetAll() []*url.URL {
	return d.extract(-1)
}

// GetNumber gets n images from the content. Returns an empty
// slice (not nil) if no image was found.
func (d ImageExtractor) GetNumber(quantity int) []*url.URL {
	return d.extract(quantity)
}

// GetFirst returns the first image found on the content. Returns nil
// if no image was found. The search stops when at least one image is
// found.
func (d ImageExtractor) GetFirst() *url.URL {
	urls := d.extract(1)

	if len(urls) >= 1 {
		return urls[0]
	}

	return nil
}

func (d ImageExtractor) extract(amount int) []*url.URL {
	urlsStr := reg.FindAllString(d.Content, amount)
	urls := make([]*url.URL, 0)

	for _, currentUrl := range urlsStr {
		if urlP, err := url.Parse(currentUrl); err == nil {
			urls = append(urls, urlP)
		}
	}

	return urls
}
