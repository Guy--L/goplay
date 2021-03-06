package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func doCrawl(url string, fetcher Fetcher, results chan []string) {
    body, urls, err := fetcher.Fetch(url)
    results <- urls

    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("found: %s %q\n", url, body)
    }
}

func Crawl(url string, depth int, fetcher Fetcher) {
    results := make(chan []string)
    crawled := make(map[string]bool)
    go doCrawl(url, fetcher, results)
    // counter for unfinished crawling goroutines
    toWait := 1

    for urls := range results {
        toWait--

        for _, u := range urls {
            if !crawled[u] {
                crawled[u] = true
                go doCrawl(u, fetcher, results)
                toWait++
            }
        }

        if toWait == 0 {
            break
        }
    }
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

