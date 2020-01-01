package main

import (
	"errors"
	"fmt"
	"log"
)

func bridgeList(control chan int, upstream side, b bridge) {
	go func() {
		defer signalFailure(control)
		for {
			if upstream.Connection.Client == nil {
				log.Printf("no connection to %v", describeSide(upstream.Config))
				break
			}

			// timeout=0 means infinite wait
			entries := upstream.Connection.Client.BLPop(0, upstream.Config.Name)
			result, err := entries.Result()
			if err != nil {
				log.Printf("failed to receive message from %v", describeSide(upstream.Config))
				break
			}
			message := result[1]
			log.Printf("message received from upstream %v: %v", describeSide(upstream.Config), message)

			err = forward(message, b.Downstreams)
			if err != nil {
				log.Printf("failed to forward: %v", err)
				break
			}
		}
	}()
}

func forwardList(message string, downstream side) error {
	if downstream.Connection.Client == nil {
		return errors.New(fmt.Sprintf("no connection to %v", describeSide(downstream.Config)))
	}
	result := downstream.Connection.Client.RPush(downstream.Config.Name, message)
	return handleResult(result, downstream, message)
}
