package main

import (
	"log"
	"sync"
)

func bridgeList(wg *sync.WaitGroup, upstream side, b bridge) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// timeout=0 means infinite wait
			entries := upstream.Client.BLPop(0, upstream.Config.Name)
			result, err := entries.Result()
			if err != nil {
				log.Printf("failed to receive message from %v", describeSide(upstream.Config))
				continue
			}
			message := result[1]
			log.Printf("message received from upstream %v: %v", describeSide(upstream.Config), message)
			forward(message, b.Downstreams)
		}
	}()
}

func forwardList(message string, downstream side) {
	downstream.Client.RPush(downstream.Config.Name, message)
}
