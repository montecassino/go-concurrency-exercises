//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// main idea is to pipe all tweets that we're receiving to the tweet channel and so that they can
// slide towards the consumer
func producer(tweetChan chan *Tweet, stream Stream) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			// close channel to prevent leakage
			close(tweetChan)
			return
		}

		tweetChan <- tweet
	}
}

func consumer(tweetChan chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tweetChan {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweetChan := make(chan *Tweet)

	// Producer
	go producer(tweetChan, stream)

	var wg sync.WaitGroup
	wg.Add(1)
	go consumer(tweetChan, &wg)

	wg.Wait() // ensure all shit has been processed

	fmt.Printf("Process took %s\n", time.Since(start))
}
