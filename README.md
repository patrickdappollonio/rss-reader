# `rss-reader`

Still under development. A RSS reader that, for the moment, only works with Tumblr blogs.

## Usage

```go
articles, err := rssreader.Setup(rssreader.Config{
	RSSURL:   "http://blog.largentfuels.com/rss",
	MaxItems: 3,
}).ReadFeed()

if err != nil {
	fmt.Println(err.Error())
	return
}

for _, v := range articles {
	fmt.Println("Title:", v.Title)
	fmt.Println("Content:", v.Content)
	fmt.Println()
}
```

## TODO

* [ ] Parse the rest of the content into the correspondent structs or pointers
* [ ] Check the `Content` attribute and fetch n number of pictures, then verify if they're larger than a specific number
* [ ] Add a flag to fetch wether if we need to find `PreviewImage`s or not
* [ ] Use goroutines in any case to "read" the images in parallel
