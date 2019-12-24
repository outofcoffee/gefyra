package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type clientConfig struct {
	Host     string
	Port     int
	Password string
}

type sideConfig struct {
	Client  clientConfig
	Channel string
}

func connectSide(config sideConfig) *redis.Client {
	log.Printf("initialising redis connection to %v", describeClient(config))
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Client.Host, config.Client.Port),
		Password: config.Client.Password,
		DB:       0, // use default DB
	})

	// healthcheck
	attempt := 0
	for {
		attempt++
		log.Printf("checking redis connection [attempt #%v]...\n", attempt)
		pong, err := client.Ping().Result()
		if err != nil || "PONG" != pong {
			log.Printf("could not ping redis at %v - will retry\n", describeClient(config))
			time.Sleep(1000 * time.Millisecond)
		} else {
			log.Printf("redis connected at %v", describeClient(config))
			break
		}
	}

	return client
}
