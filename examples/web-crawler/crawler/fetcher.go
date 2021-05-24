package crawler

// API inspired from https://gist.github.com/harryhare/6a4979aa7f8b90db6cbc74400d0beb49
// Implementation inspiration from: https://github.com/JackDanger/gocrawler/blob/master/crawl.go

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// TODO fetcher implementation here
