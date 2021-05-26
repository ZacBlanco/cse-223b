package main

import (
	"fmt"
	"hash"
	"hash/fnv"
	"net/http"
	"net/url"

	"github.com/apache/openwhisk-client-go/whisk"
)

const (
	SEEDS_KEY  = "seeds"
	STATE_KEY  = "state"
	NUM_ACTORS = 2
)

// Parent Actor
func Main(args map[string]interface{}, state *interface{}) map[string]interface{} {
	var unorderedChildActorStates []WebCrawlerState
	actorStates := make(chan WebCrawlerState)

	// Start each child actor
	for i := 0; i < NUM_ACTORS; i++ {
		go StartIthChildWebCrawlerAndGetState(i, args, actorStates)
	}

	for i := 0; i < NUM_ACTORS; i++ {
		unorderedChildActorStates = append(unorderedChildActorStates, <-actorStates)
	}

	return map[string]interface{}{
		"actorStates": unorderedChildActorStates,
	}
}

// Child Actor: crawls web pages
func MainChild1(args map[string]interface{}, state *interface{}) map[string]interface{} {
	const actorId = 0
	seeds := GetSeedsFromArgs(args)
	actor := NewActorWithId(actorId)
	actor.UseLatestState(state)
	return actor.StartWebCrawler(seeds)
}

// Child Actor: crawls web pages
func MainChild2(args map[string]interface{}, state *interface{}) map[string]interface{} {
	const actorId = 1
	seeds := GetSeedsFromArgs(args)
	actor := NewActorWithId(actorId)
	actor.UseLatestState(state)
	return actor.StartWebCrawler(seeds)
}

type Actor struct {
	Id        int
	State     WebCrawlerState
	MyFetcher Fetcher // retrieves web pages
	Hasher    hash.Hash64
}

type WebCrawlerState struct {
	Seen    map[string]bool // seeds already visited
	StateId int             // monotonic increasing counter. Larger StateIds are more recent.
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func (actor *Actor) UseLatestState(state *interface{}) {
	if newState, ok := (*state).(WebCrawlerState); ok {
		if newState.StateId > actor.State.StateId {
			actor.State = newState
		}
	}
}

func (actor *Actor) StartWebCrawler(seeds []string) map[string]interface{} {
	actor.crawl(seeds)
	return map[string]interface{}{
		"state": actor.State,
	}
}

/* === Actor implementations === */

func (actor *Actor) crawl(seeds []string) {
	// url references
	urlRefsToVisit := make(chan string)

	actor.populateUrlRefsToVisit(urlRefsToVisit, seeds)

	done := false
	for !done {
		select {
		case urlRef := <-urlRefsToVisit:
			actor.fetchAndUpdate(urlRef, urlRefsToVisit)
		default:
			done = true // no more urls to visit
		}
	}
}

func (actor *Actor) populateUrlRefsToVisit(urlRefsToVisit chan string, urlRefs []string) {
	for _, urlRef := range urlRefs {
		if actor.shouldVisitUrl(urlRef) {
			urlRefsToVisit <- urlRef
		}
	}
}

func (actor *Actor) shouldVisitUrl(urlRef string) bool {
	if actor.State.Seen[urlRef] {
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
	if _, urlRefs, e := actor.MyFetcher.Fetch(urlRef); e == nil {
		actor.State.Seen[urlRef] = true
		actor.populateUrlRefsToVisit(urlRefsToVisit, urlRefs)
	}
}

/* === helper functions === */

func StartIthChildWebCrawlerAndGetState(
	id int,
	args map[string]interface{},
	actorStates chan WebCrawlerState,
) {
	// https://pkg.go.dev/github.com/apache/openwhisk-client-go@v0.0.0-20210313152306-ea317ea2794c/whisk#section-documentation
	config := &whisk.Config{
		AuthToken: "23bc46b1-71f6-4ed5-8c54-816aa4f8c502:123zO3xZCLrMN6v2BKK1dXYFpXlPkccOFqm12CdAsMgRU4VrNZ9lyGVCGuMDGIwP",
		Host:      "http://localhost:3233",
		Insecure:  true,
	}

	client, err := whisk.NewClient(http.DefaultClient, config)
	if err != nil {
		actorStates <- WebCrawlerState{}
		return
	}

	res, _, err := client.Actions.Invoke(fmt.Sprintf("web-crawler-c%v", id), args, true, true)

	if err != nil {
		actorStates <- WebCrawlerState{}
		return
	}

	if state, ok := res[STATE_KEY].(WebCrawlerState); ok {
		actorStates <- state
	} else {
		actorStates <- WebCrawlerState{}
	}
	//actorStates <- WebCrawlerState{}
}

func GetSeedsFromArgs(args map[string]interface{}) []string {
	// Get the seed and try to cast to []string.
	// Return empty []string if SEED_KEY doesn't exist or error with cast,
	if seed, ok := args[SEEDS_KEY].([]string); seed == nil && ok {
		return make([]string, 0)
	} else {
		return seed
	}
}

func NewActorWithId(actorId int) Actor {
	// TODO add new Fetcher
	return Actor{
		Id:     actorId,
		Hasher: fnv.New64(),
	}
}
