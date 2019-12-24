package main

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

func bridgeChannels(wg *sync.WaitGroup, b bridge) {
	log.Printf("starting bridge %v", describeBridge(b))
	for _, upstream := range b.Upstreams {
		wg.Add(1)
		pubSub := upstream.Client.Subscribe(upstream.Config.Channel)

		go func() {
			defer wg.Done()
			for {
				message, err := pubSub.ReceiveMessage()
				if err != nil {
					log.Printf("failed to receive message from %v", describeSide(upstream.Config))
					continue
				}
				log.Printf("message received from upstream %v: %v", describeSide(upstream.Config), message.Payload)
				forward(message, b.Downstreams)
			}
		}()
	}
}

func forward(message *redis.Message, downstreams []side) {
	for _, downstream := range downstreams {
		log.Printf("forwarding message to downstream %v", describeSide(downstream.Config))
		downstream.Client.Publish(downstream.Config.Channel, message.Payload)
	}
}
