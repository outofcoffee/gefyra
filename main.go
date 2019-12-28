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
	bridgeConfigs := loadConfig()

	for bridgeName, bridgeConfig := range bridgeConfigs {
		log.Printf("bridge %v has %v upstream(s) and %v downstream(s)", bridgeName, len(bridgeConfig.Upstreams), len(bridgeConfig.Downstreams))
		upstreams, downstreams := connect(bridgeConfig)
		b := bridge{
			Upstreams:   upstreams,
			Downstreams: downstreams,
		}
		bridges = append(bridges, b)
	}
	defer cleanup(bridges)

	var wg = sync.WaitGroup{}
	for _, b := range bridges {
		bridgeUp(&wg, b)
	}
	wg.Wait()
}
