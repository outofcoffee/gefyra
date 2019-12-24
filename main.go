package main

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
)

type side struct {
	Config sideConfig
	Client *redis.Client
}

type bridge struct {
	Upstreams   []side
	Downstreams []side
}

var bridges []bridge

func main() {
	config := loadConfig()

	upstreams, downstreams := connect(config)
	b := bridge{
		Upstreams:   upstreams,
		Downstreams: downstreams,
	}
	bridges = append(bridges, b)
	defer cleanup(b)

	var wg = sync.WaitGroup{}
	bridgeChannels(&wg, b)
	wg.Wait()
}

func cleanup(b bridge) {
	for _, upstream := range b.Upstreams {
		upstream.Client.Close()
	}
	for _, downstream := range b.Downstreams {
		downstream.Client.Close()
	}
}

func connect(config bridgeConfig) (upstreams []side, downstreams []side) {
	log.Print("connecting to upstreams")
	for _, c := range config.Upstreams {
		upstreams = append(upstreams, side{
			Config: c,
			Client: connectSide(c),
		})
	}

	log.Print("connecting to downstreams")
	for _, c := range config.Downstreams {
		downstreams = append(downstreams, side{
			Config: c,
			Client: connectSide(c),
		})
	}

	return upstreams, downstreams
}

func fatalIfError(err interface{}, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}
