package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type configWrapper struct {
	Bridges map[string]bridgeConfig
}

type bridgeConfig struct {
	Upstreams   []sideConfig
	Downstreams []sideConfig
}

type sideConfig struct {
	Client clientConfig
	Name   string
	Type   sideType
}

type sideType string

const (
	PubSub sideType = "pubsub"
	List            = "list"
)

func loadConfig() map[string]bridgeConfig {
	var configFile string
	if value, ok := os.LookupEnv("BRIDGE_CONFIG"); ok {
		configFile = value
	} else {
		log.Fatal("missing bridge config file")
	}
	log.Printf("loading config file: %v", configFile)

	wrapper := configWrapper{}
	yamlFile, err := ioutil.ReadFile(configFile)
	fatalIfError(err, "failed to read config file")
	err = yaml.Unmarshal([]byte(yamlFile), &wrapper)
	fatalIfError(err, "failed to parse config")

	bridges := wrapper.Bridges
	normalise(bridges)
	log.Printf("loaded %v bridge(s)", len(bridges))
	return bridges
}

func normalise(bridges map[string]bridgeConfig) {
	for _, config := range bridges {
		for _, c := range config.Upstreams {
			normaliseConfig(c)
		}
		for _, c := range config.Downstreams {
			normaliseConfig(c)
		}
	}
}

func normaliseConfig(c sideConfig) {
	if c.Client.Port == 0 {
		c.Client.Port = 6379
	}
	if len(c.Type) == 0 {
		panic("bridge side type not set")
	}
}

func describeClient(config sideConfig) string {
	return fmt.Sprintf("%v:%v", config.Client.Host, config.Client.Port)
}

func describeSide(config sideConfig) string {
	return fmt.Sprintf("%v/%v=%v", describeClient(config), config.Type, config.Name)
}

func describeBridge(b bridge) string {
	return fmt.Sprintf("%v->%v", describeSide(b.Upstreams[0].Config), describeSide(b.Downstreams[0].Config))
}
