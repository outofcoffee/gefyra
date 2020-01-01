package main

import (
	"github.com/go-redis/redis"
	"sync"
)

type serverConnection struct {
	Client *redis.Client
}

type side struct {
	Config     sideConfig
	Connection *serverConnection
}

type bridge struct {
	Upstreams   []side
	Downstreams []side
}

func main() {
	bridgeConfigs := loadConfig()
	startDisconnectionWorker()

	var wg = sync.WaitGroup{}
	var bridges []bridge
	for bridgeName, bc := range bridgeConfigs {
		b := initBridge(&wg, bridgeName, bc)
		bridges = append(bridges, b)
	}

	defer disconnectBridges(bridges)
	wg.Wait()
}
