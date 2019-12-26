package main

import (
	"github.com/go-redis/redis"
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
	bridgeUp(&wg, b)
	wg.Wait()
}
