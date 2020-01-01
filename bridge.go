package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"sync"
)

func initBridge(wg *sync.WaitGroup, bridgeName string, bc bridgeConfig) bridge {
	log.Printf("bridge %v has %v upstream(s) and %v downstream(s)", bridgeName, len(bc.Upstreams), len(bc.Downstreams))

	wg.Add(1)
	b := bridge{}

	go func() {
		defer wg.Done()

		attempt := 0
		for {
			attempt++
			log.Printf("starting bridge: %v [attempt #%v]", describeBridgeConfig(bc), attempt)

			control := make(chan int)
			upstreams, downstreams := connect(bc, control)
			b.Upstreams = upstreams
			b.Downstreams = downstreams

			err := bridgeUp(control, b)
			fatalIfError(err, fmt.Sprintf("could not bring up bridge %v", describeBridge(b)))

			<-control
			log.Printf("bridge %v going down", describeBridge(b))
			disconnectBridge(b)
		}
	}()

	return b
}

func bridgeUp(control chan int, b bridge) error {
	log.Printf("starting bridge %v", describeBridge(b))
	for _, upstream := range b.Upstreams {
		log.Printf("listening to %v", describeSide(upstream.Config))
		switch upstream.Config.Type {
		case PubSub:
			bridgeChannel(control, upstream, b)
		case List:
			bridgeList(control, upstream, b)
		default:
			return errors.New(fmt.Sprintf("unsupported upstream type: %v", upstream.Config.Type))
		}
	}
	return nil
}

func forward(message string, downstreams []side) error {
	for _, downstream := range downstreams {
		log.Printf("forwarding message to downstream %v", describeSide(downstream.Config))
		switch downstream.Config.Type {
		case PubSub:
			return forwardChannel(message, downstream)
		case List:
			return forwardList(message, downstream)
		default:
			return errors.New(fmt.Sprintf("unsupported downstream type: %v", downstream.Config.Type))
		}
	}
	return nil
}

func handleResult(result *redis.IntCmd, downstream side, message string) error {
	if result.Err() != nil {
		return errors.New(fmt.Sprintf("failed to forward message to %v: %v: %v", describeSide(downstream.Config), message, result.Err()))
	}
	return nil
}
