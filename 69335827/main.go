package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const ConfigFileEnv = "ConfigFile"

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	var config Config
	var configfile string
	if configfile = os.Getenv(ConfigFileEnv); configfile == "" {
		configfile = "config.json"
	}
	if f, err := os.Open(configfile); err != nil {
		panic(fmt.Errorf("Couldn't open config file %s: %w", configfile, err))
	} else if err := json.NewDecoder(f).Decode(&config); err != nil {
		panic(fmt.Errorf("Could not decode config file %s as Json: %w", configfile, err))
	}
	fmt.Printf("%+v\n", config)
}
