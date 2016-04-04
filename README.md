# `rss-reader`

[![GoDoc](https://godoc.org/github.com/patrickdappollonio/rss-reader?status.svg)](https://godoc.org/github.com/patrickdappollonio/rss-reader)

An RSS reader that, for the moment, only works with Tumblr blogs. It also fetches
all the images within the content and pick the most appropiate one based on a minimum
width or height.

As a recommendation, it's better to use this package with caching enabled, since if there's
any RSS with several images, you'll ended up using a bit of bandwidth.

### Usage

```go
articles, err := rssreader.Setup(rssreader.Config{
	RSSURL:        "http://your.weblog.site/rss",
	MaxItems:      3,
	MinImageWidth: 200,
	UseCache:      true,
}).ReadFeed()

if err != nil {
	fmt.Println(err.Error())
	return
}

for _, v := range articles {
	fmt.Println("Title:", v.Title)
	fmt.Println("Content:", v.Content)
	fmt.Println("PreviewImage:", v.PreviewImage)
	fmt.Println("Published at:", v.Published)
	fmt.Println()
}
```

### Benchmarks with cache and no cache

For cache, the package uses an in-memory cache implementation with a single record and a
default expiration of 24 hours —meaning the RSS items are going to be cached for 24 hours—.

```
PASS
BenchmarkWithCache-8  	20000000	        74.5 ns/op
BenchmarkWithNoCache-8	       5	 289574604 ns/op
ok  	github.com/patrickdappollonio/rss-reader	5.272s
```
