package main

import (
	"log"
	"sync"
)

func bridgeUp(wg *sync.WaitGroup, b bridge) {
	log.Printf("starting bridge %v", describeBridge(b))
	for _, upstream := range b.Upstreams {
		log.Printf("listening to %v", describeSide(upstream.Config))
		switch upstream.Config.Type {
		case PubSub:
			bridgeChannel(wg, upstream, b)
		case List:
			bridgeList(wg, upstream, b)
		}
	}
}

func forward(message string, downstreams []side) {
	for _, downstream := range downstreams {
		log.Printf("forwarding message to downstream %v", describeSide(downstream.Config))
		switch downstream.Config.Type {
		case PubSub:
			forwardChannel(message, downstream)
		case List:
			forwardList(message, downstream)
		}
	}
}
