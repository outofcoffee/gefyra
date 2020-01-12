package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type serverConfig struct {
	Host               string
	Port               int
	Password           string
	ConnectionAttempts int `yaml:"connection_attempts"`
}

var disconnectionRequests = make(chan side)

func connect(control chan int, config bridgeConfig) (upstreams []side, downstreams []side) {
	log.Print("connecting to upstreams")
	for _, c := range config.Upstreams {
		upstreams = append(upstreams, connectSide(c, control))
	}

	log.Print("connecting to downstreams")
	for _, c := range config.Downstreams {
		downstreams = append(downstreams, connectSide(c, control))
	}

	return upstreams, downstreams
}

func connectSide(config sideConfig, control chan int) side {
	log.Printf("initialising connection to %v", describeClient(config))
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Server.Host, config.Server.Port),
		Password: config.Server.Password,
		DB:       0, // use default DB
	})

	s := side{
		Config: config,
		Connection: &serverConnection{
			Client: client,
		},
	}
	err := ping(s, 1*time.Second, config.Server.ConnectionAttempts, true)
	fatalIfError(err, "unable to connect")

	monitorLiveness(control, s)
	return s
}

func monitorLiveness(control chan int, s side) {
	go func() {
		defer signalFailure(control)
		for {
			time.Sleep(60 * time.Second)
			if s.Connection.Client == nil {
				// The client connection is no longer in use, so
				// this healthcheck is no longer required.
				break
			}
			err := ping(s, 0, 1, false)
			if nil != err {
				log.Printf("liveness check failed (%v)", err)
				break
			}
		}
	}()
}

func ping(s side, checkInterval time.Duration, maxAttempts int, logAttempt bool) error {
	attempt := 0
	for {
		attempt++
		if logAttempt {
			log.Printf("checking connection [attempt #%v]...", attempt)
		}
		pong, err := s.Connection.Client.Ping().Result()
		if err != nil || "PONG" != pong {
			msg := fmt.Sprintf("could not ping server at %v", describeClient(s.Config))
			if maxAttempts == 0 || attempt < maxAttempts {
				log.Printf("%v - will retry in %v", msg, checkInterval)
				time.Sleep(checkInterval)
			} else {
				return errors.New(fmt.Sprintf("%v after %v attempts: %v", msg, attempt, err))
			}
		} else {
			if logAttempt {
				log.Printf("connected to %v", describeClient(s.Config))
			}
			break
		}
	}
	return nil
}

func disconnectBridges(bridges []bridge) {
	for _, b := range bridges {
		disconnectBridge(b)
	}
}

func disconnectBridge(b bridge) {
	for _, upstream := range b.Upstreams {
		if upstream.Connection.Client != nil {
			disconnectSide(upstream)
		}
	}
	for _, downstream := range b.Downstreams {
		if downstream.Connection.Client != nil {
			disconnectSide(downstream)
		}
	}
}

func disconnectSide(s side) {
	disconnectionRequests <- s
}

// Ensures ordered access when disconnecting clients
func startDisconnectionWorker() {
	go func() {
		for s := range disconnectionRequests {
			if s.Connection.Client == nil {
				continue
			}
			log.Printf("disconnecting %v", describeSide(s.Config))
			_ = s.Connection.Client.Close()
			s.Connection.Client = nil
		}
	}()
}
