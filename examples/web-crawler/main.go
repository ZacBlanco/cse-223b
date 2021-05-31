package main

import (
	"fmt"
	"hash"
	"hash/fnv"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	SEEDS_KEY              = "seeds"
	STATE_KEY              = "state"
	MAX_CHANNEL_SIZE       = 100
	NUM_ACTORS             = 1
	MAX_NUM_PAGES_TO_VISIT = 4
	WSK_HOST               = "http://172.17.0.1:3233"
)

type Actor struct {
	Id        int
	State     WebCrawlerState
	MyFetcher Fetcher // retrieves web pages
	Hasher    hash.Hash64
}

type WebCrawlerState struct {
	Id   int
	Seen map[string]bool // seeds already visited
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
}

type Interface interface{}

func Main(args map[string]interface{}, state *interface{}) map[string]interface{} {
	// single actor.
	child_ret := Mainchild1(args, state)
	(*state) = child_ret[STATE_KEY]
	return child_ret
}

func Mainchild1(args map[string]interface{}, state *interface{}) map[string]interface{} {
	const actorId = 0
	seeds := GetSeedsFromArgs(args)
	rand.Shuffle(len(seeds), func(i, j int) { seeds[i], seeds[j] = seeds[j], seeds[i] })
	actor := NewActorWithId(actorId)
	actor.UseLatestState(state)
	return actor.StartWebCrawlerAndReturnWebCrawlerState(seeds)
}

func ParseActorStates(state *interface{}) []WebCrawlerState {
	if state == nil {
		return make([]WebCrawlerState, 0)
	}
	if states, ok := (*state).([]WebCrawlerState); ok {
		// order the states by ID
		return states
	}
	return make([]WebCrawlerState, 0)
}

func (actor *Actor) UseLatestState(state *interface{}) {
	if state == nil {
		return
	}
	if newState, ok := (*state).(WebCrawlerState); ok {
		actor.State = newState
	}
}

func (actor *Actor) StartWebCrawlerAndReturnWebCrawlerState(seeds []string) map[string]interface{} {
	actor.crawl(seeds)
	return map[string]interface{}{
		STATE_KEY: actor.State,
	}
}

/* === Actor implementations === */

func (actor *Actor) crawl(seeds []string) {
	// url references
	urlRefsToVisit := make(chan string, MAX_CHANNEL_SIZE)

	fmt.Println("Populating seeds:", actor.Id, seeds)
	actor.populateUrlRefsToVisit(urlRefsToVisit, seeds)

	fmt.Println("Start crawling", actor.Id)
	done := false
	numVisitedPages := 0
	for !done {
		select {
		case urlRef := <-urlRefsToVisit:
			if !actor.State.Seen[urlRef] {
				numVisitedPages += 1
			}
			actor.fetchAndUpdate(urlRef, urlRefsToVisit)
			if numVisitedPages >= MAX_NUM_PAGES_TO_VISIT {
				done = true
			}
		default:
			fmt.Println("No more urls to visit")
			done = true // no more urls to visit
		}
	}
	fmt.Println("Done", actor.Id)
}

func (actor *Actor) populateUrlRefsToVisit(urlRefsToVisit chan string, urlRefs []string) {
	for _, urlRef := range urlRefs {
		if actor.shouldVisitUrl(urlRef) {
			select {
			case urlRefsToVisit <- urlRef: // best effor insert into channel.
			default: // channel is full. just skip this entry
			}
		}
	}
}

func (actor *Actor) shouldVisitUrl(urlRef string) bool {
	if actor.State.Seen[urlRef] {
		fmt.Println("Already seen URL:", urlRef)
		return false
	}

	if myUrl, e := url.Parse(urlRef); e != nil {
		return false
	} else {
		return actor.hashTheNameAndReturnOwnerId(myUrl.Hostname()) == actor.Id
	}
}

func (actor *Actor) hashTheNameAndReturnOwnerId(name string) int {
	actor.Hasher.Reset()
	actor.Hasher.Sum([]byte(name))
	return int(actor.Hasher.Sum64() % uint64(NUM_ACTORS))
}

func (actor *Actor) fetchAndUpdate(urlRef string, urlRefsToVisit chan string) {
	fmt.Println("fetchAndUpdate:", actor.State.Id, urlRef)
	if actor.MyFetcher == nil || actor.State.Seen[urlRef] {
		return
	}
	if urlRefs, e := actor.MyFetcher.Fetch(urlRef); e == nil {
		actor.State.Seen[urlRef] = true
		actor.populateUrlRefsToVisit(urlRefsToVisit, urlRefs)
	}
}

/* === helper functions === */

func GetSeedsFromArgs(args map[string]interface{}) []string {
	// Get the seed and try to cast to []string.
	// Return empty []string if SEED_KEY doesn't exist or error with cast,
	fmt.Println("Try to get seeds:", args[SEEDS_KEY])

	// TODO this cast can be error prone.
	s := make([]string, 0)
	for _, v := range args[SEEDS_KEY].([]interface{}) {
		s = append(s, fmt.Sprint(v))
	}
	return s
}

func NewWebCrawlerStateWithId(id int) WebCrawlerState {
	return WebCrawlerState{Id: id, Seen: make(map[string]bool)}
}

func NewActorWithId(actorId int) Actor {
	return Actor{
		Id:        actorId,
		State:     NewWebCrawlerStateWithId(actorId),
		MyFetcher: &WebFetcher{},
		Hasher:    fnv.New64(),
	}
}

/* === Fetcher Implementation */

type WebFetcher struct{}

// https://www.devdungeon.com/content/web-scraping-go
func (fetcher *WebFetcher) Fetch(url string) ([]string, error) {
	fmt.Println("Fetch page:", url)
	// Make HTTP request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Unable to connect:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Error loading HTTP response body:", err)
		return nil, err
	}

	urls := make([]string, 0)

	processElement :=
		func(index int, element *goquery.Selection) {
			// See if the href attribute exists on the element
			href, exists := element.Attr("href")
			if exists {
				urls = append(urls, href)
			}
		}

	document.Find("a").Each(processElement)
	return urls, nil
}
