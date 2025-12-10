//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// my idea is to create a global channel that acts as a gate pass holder, this gate pass holder can obtain a pass every second
// and then pass this 'pass' to a Crawl goroutine (whoever Crawl goroutine that can first grab it wins, while other goroutines will be blocked from fetching)

var globalChannel = make(chan bool, 1)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	// block other go routines from executing except for 1 that can obtain the pass
	<- globalChannel

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg)
	}
}

func main() {
	var wg sync.WaitGroup

	go func(){
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// hand over a pass to global channel every sec
			globalChannel <- true
		}
	}()

	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg)
	wg.Wait()
}
