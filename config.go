package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type configWrapper struct {
	Bridge bridgeConfig
}

type bridgeConfig struct {
	Upstreams   []sideConfig
	Downstreams []sideConfig
}

func loadConfig() bridgeConfig {
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

	config := wrapper.Bridge
	normalise(config)
	log.Printf("loaded %v upstream(s) and %v downstream(s)", len(config.Upstreams), len(config.Downstreams))
	return config
}

func normalise(config bridgeConfig) {
	for _, c := range config.Upstreams {
		if c.Client.Port == 0 {
			c.Client.Port = 6379
		}
	}
	for _, c := range config.Downstreams {
		if c.Client.Port == 0 {
			c.Client.Port = 6379
		}
	}
}

func describeClient(config sideConfig) string {
	return fmt.Sprintf("%v:%v", config.Client.Host, config.Client.Port)
}

func describeSide(config sideConfig) string {
	return fmt.Sprintf("%v/%v", describeClient(config), config.Channel)
}

func describeBridge(b bridge) string {
	return fmt.Sprintf("%v->%v", describeSide(b.Upstreams[0].Config), describeSide(b.Downstreams[0].Config))
}
