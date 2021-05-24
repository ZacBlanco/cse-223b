package main

import "crawler"

// Does the work of the web crawler
func Main(args map[string]interface{}, state interface{}) map[string]interface{} {
	seeds := crawler.GetSeedsFromArgs(args)
	return crawler.StartWebCrawler(seeds)
}
