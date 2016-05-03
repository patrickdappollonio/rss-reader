package rssreader

import (
	"testing"
)

const BLOG_URL = "http://blog.largentfuels.com/rss"

func generateReader(items, minWidth int, cache bool) *RSS {
	return Setup(Config{
		RSSURL:        BLOG_URL,
		MaxItems:      items,
		MinImageWidth: minWidth,
		UseCache:      cache,
	})
}

func generateReaderWithURL(url string, items, minWidth int, cache bool) *RSS {
	return Setup(Config{
		RSSURL:        url,
		MaxItems:      items,
		MinImageWidth: minWidth,
		UseCache:      cache,
	})
}

func BenchmarkWithCache(b *testing.B) {
	getter := generateReader(10, 200, true)

	for n := 0; n < b.N; n++ {
		getter.ReadFeed()
	}
}

func BenchmarkWithNoCache(b *testing.B) {
	getter := generateReader(10, 200, false)

	for n := 0; n < b.N; n++ {
		getter.ReadFeed()
	}
}

func TestReader(t *testing.T) {
	expected := "Visítanos en nuestra página de Facebook!…. Te esperamos"
	news, err := generateReader(1, 200, true).ReadFeed()

	if err != nil {
		t.Fatalf("Error reading feed: %v", err.Error())
	}

	if len(news) != 1 {
		t.Fatalf("Expected 1 elements in reader, got %v", len(news))
	}

	if news[0].Title != expected {
		t.Errorf(`Expected title "%v", got "%v"`, news[0].Title)
	}
}

func TestReaderNoRSS(t *testing.T) {
	_, err := generateReaderWithURL("", 1, 200, true).ReadFeed()

	if err != ErrNoURLProvided {
		t.Fatalf(`Sent empty url, received error: %s`, err.Error())
	}
}

func TestReaderRelativeURL(t *testing.T) {
	_, err := generateReaderWithURL("/rss/", 1, 200, true).ReadFeed()

	if err != ErrNeededAbsURL {
		t.Fatalf(`Sent relative url, received error: %s`, err.Error())
	}
}

func TestReaderMalformedURL(t *testing.T) {
	_, err := generateReaderWithURL("http~://google.com", 1, 200, true).ReadFeed()

	if err != ErrMalformedURL {
		t.Fatalf(`Sent malformed url, received error: %s`, err.Error())
	}
}

func TestReaderCantConnectURL(t *testing.T) {
	_, err := generateReaderWithURL("http://localhost/", 1, 200, true).ReadFeed()

	if err != ErrCantConnect {
		t.Fatalf(`Sent localhost url so it can't connect, received error: %s`, err.Error())
	}
}

func TestReaderNonRSSURL(t *testing.T) {
	_, err := generateReaderWithURL("http://example.com/", 1, 200, true).ReadFeed()

	if err != ErrCantParseResponse {
		t.Fatalf(`Sent a non-xml URL, received error: %s`, err.Error())
	}
}

func TestReaderNotInitialized(t *testing.T) {
	var reader *RSS = nil
	_, err := reader.ReadFeed()

	if err != ErrNotInitialized {
		t.Fatalf("Sent an uninitialized reader, got %v", err.Error())
	}
}
