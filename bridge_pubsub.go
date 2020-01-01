package main

import (
	"errors"
	"fmt"
	"log"
)

func bridgeChannel(control chan int, upstream side, b bridge) {
	go func() {
		defer signalFailure(control)
		pubSub := upstream.Connection.Client.Subscribe(upstream.Config.Name)
		for {
			if upstream.Connection.Client == nil {
				log.Printf("no connection to %v", describeSide(upstream.Config))
				break
			}

			message, err := pubSub.ReceiveMessage()
			if err != nil {
				log.Printf("failed to receive message from %v", describeSide(upstream.Config))
				break
			}
			log.Printf("message received from upstream %v: %v", describeSide(upstream.Config), message.Payload)

			err = forward(message.Payload, b.Downstreams)
			if err != nil {
				log.Printf("failed to forward: %v", err)
				break
			}
		}
	}()
}

func forwardChannel(message string, downstream side) error {
	if downstream.Connection.Client == nil {
		return errors.New(fmt.Sprintf("no connection to %v", describeSide(downstream.Config)))
	}
	result := downstream.Connection.Client.Publish(downstream.Config.Name, message)
	return handleResult(result, downstream, message)
}
