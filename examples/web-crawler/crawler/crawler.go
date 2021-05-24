package crawler

import (
	"hash"
	"hash/fnv"
	"net/url"
	"sync"
)

const (
	SEEDS = "seeds"
)

type WebCrawlerState struct {
	Id        int             //
	NumActors int             // total number of actors.
	Seen      map[string]bool // seeds already visited
	StateId   int             // monotonically increase state counter. larger StateId's are newer
}

func StartWebCrawler(seeds []string) map[string]interface{} {
	var wg sync.WaitGroup
	// Holds functions of child actors.
	actorFuncs := []func(event map[string]interface{}, state interface{}) map[string]interface{}{
		ChildMain1, ChildMain2}

	// Start each child actor
	for i, fn := range actorFuncs {
		var state = WebCrawlerState{
			Id:        i,
			NumActors: len(actorFuncs),
			Seen:      make(map[string]bool),
			StateId:   0,
		}
		event := make(map[string]interface{})
		event[SEEDS] = seeds
		go addFunctionToWaitGroupAndExecute(&wg, fn, event, state)
	}
	wg.Wait() // wait for all the child actors to complete.
	return nil
}

func ChildMain1(event map[string]interface{}, state interface{}) map[string]interface{} {
	actor := NewChildActor()
	actor.useLatestState(&state)
	actor.crawlAndCheckpoint(GetSeedsFromArgs(event))
	return nil
}

func ChildMain2(event map[string]interface{}, state interface{}) map[string]interface{} {
	actor := NewChildActor()
	actor.useLatestState(&state)
	actor.crawlAndCheckpoint(GetSeedsFromArgs(event))
	return nil
}

/* === ChildActor implementations === */

type childActor struct {
	State   WebCrawlerState
	Fetcher Fetcher
	Hasher  hash.Hash64
}

// Take the largest StateId of newState, previousState, and current state
func (actor *childActor) useLatestState(state *interface{}) {

	previousState := WebCrawlerState{} // TODO should fetch from checkpoint states.

	if previousState.StateId > actor.State.StateId {
		actor.State = previousState
	}

	if newState, ok := (*state).(WebCrawlerState); ok {
		if newState.StateId > actor.State.StateId {
			actor.State = newState
		}
	}
}

func (actor *childActor) crawlAndCheckpoint(seeds []string) {
	// url references
	urlRefsToVisit := make(chan string)

	actor.populateUrlRefsToVisit(urlRefsToVisit, seeds)

	done := false
	for !done {
		select {
		case urlRef := <-urlRefsToVisit:
			actor.crawl(urlRef, urlRefsToVisit)
			actor.checkpoint()
		default:
			done = true // no more urls to visit
		}
	}
}

func (actor *childActor) checkpoint() {
	// TODO
	// Increment the current stateId and save this new state.
}

func (actor *childActor) shouldVisitUrl(urlRef string) bool {
	if actor.State.Seen[urlRef] {
		return false
	}

	if myUrl, e := url.Parse(urlRef); e != nil {
		return false
	} else {
		return actor.hashTheNameAndReturnOwnerId(myUrl.Hostname()) == actor.State.Id
	}
}

func (actor *childActor) hashTheNameAndReturnOwnerId(name string) int {
	actor.Hasher.Reset()
	actor.Hasher.Sum([]byte(name))
	return int(actor.Hasher.Sum64() % uint64(actor.State.NumActors))
}

func (actor *childActor) populateUrlRefsToVisit(urlRefsToVisit chan string, urlRefs []string) {
	for _, urlRef := range urlRefs {
		if actor.shouldVisitUrl(urlRef) {
			urlRefsToVisit <- urlRef
		}
	}
}

func (actor *childActor) crawl(urlRef string, urlRefsToVisit chan string) {
	// optionally save the downloaded page here
	if _, urlRefs, e := actor.Fetcher.Fetch(urlRef); e == nil {
		actor.State.Seen[urlRef] = true
		actor.populateUrlRefsToVisit(urlRefsToVisit, urlRefs)
	}
}

/* === helper functions === */

func addFunctionToWaitGroupAndExecute(wg *sync.WaitGroup,
	fn func(map[string]interface{}, interface{}) map[string]interface{},
	event map[string]interface{},
	state interface{}) {

	wg.Add(1)        // increment wait group counter
	fn(event, state) // execute job
	wg.Done()        // notify completion
}

func GetSeedsFromArgs(args map[string]interface{}) []string {
	// Get the seed and try to cast to []string.
	// If key SEED doesn't exist or cannot cast return empty []string.
	if seed, ok := args[SEEDS].([]string); seed == nil && ok {
		return make([]string, 0)
	} else {
		return seed
	}
}

func NewChildActor() childActor {
	// TODO add new Fetcher
	return childActor{
		Hasher: fnv.New64(),
	}
}
