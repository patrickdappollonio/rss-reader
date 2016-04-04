package rssreader

import "testing"

func BenchmarkWithCache(b *testing.B) {
	getter := Setup(Config{
		RSSURL:        "http://blog.largentfuels.com/rss",
		MaxItems:      3,
		MinImageWidth: 200,
		UseCache:      true,
	})

	for n := 0; n < b.N; n++ {
		getter.ReadFeed()
	}
}

func BenchmarkWithNoCache(b *testing.B) {
	getter := Setup(Config{
		RSSURL:        "http://blog.largentfuels.com/rss",
		MaxItems:      3,
		MinImageWidth: 200,
		UseCache:      false,
	})

	for n := 0; n < b.N; n++ {
		getter.ReadFeed()
	}
}
