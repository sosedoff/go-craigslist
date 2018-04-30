package craigslist

import (
	"errors"
	"log"
	"sync"
)

func MultiSearch(sites []string, opts SearchOptions) (map[string]*SearchResults, error) {
	if len(sites) == 0 {
		return nil, errors.New("sites list can't be empty")
	}

	res := map[string]*SearchResults{}

	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(sites))

	for _, siteId := range sites {
		go func(id string) {
			defer wg.Done()

			results, err := Search(id, opts)
			if err != nil {
				log.Println("search error:", err)
				return
			}

			lock.Lock()
			res[id] = results
			lock.Unlock()
		}(siteId)
	}

	wg.Wait()

	return res, nil
}
