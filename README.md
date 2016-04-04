# `rss-reader`

Still under development. A RSS reader that, for the moment, only works with Tumblr blogs.

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

```
PASS
BenchmarkWithCache-8  	20000000	        74.5 ns/op
BenchmarkWithNoCache-8	       5	 289574604 ns/op
ok  	github.com/patrickdappollonio/rss-reader	5.272s
```
