package main

import (
	"log"
	"sync"
)

func bridgeChannel(wg *sync.WaitGroup, upstream side, b bridge) {
	wg.Add(1)
	pubSub := upstream.Client.Subscribe(upstream.Config.Name)

	go func() {
		defer wg.Done()
		for {
			message, err := pubSub.ReceiveMessage()
			if err != nil {
				log.Printf("failed to receive message from %v", describeSide(upstream.Config))
				continue
			}
			log.Printf("message received from upstream %v: %v", describeSide(upstream.Config), message.Payload)
			forward(message.Payload, b.Downstreams)
		}
	}()
}

func forwardChannel(message string, downstream side) {
	downstream.Client.Publish(downstream.Config.Name, message)
}
