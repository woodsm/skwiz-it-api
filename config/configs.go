package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfig() Configuration {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Config File Missing: ", err)
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	return config
}
