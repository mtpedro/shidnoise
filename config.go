package main

import (
	"encoding/json"
	"log"
	"os"
)

/*
	Read the config.json file
	Write contents to cfg variable.
*/

func init() {
	GetConfig(&cfg);
}

type Config struct {
	Token string
	Prefix string
}; 

func GetConfig(cfg *Config) {
	data, err := os.ReadFile("config.json");
	if err != nil { log.Fatal(FailedConfigOpen) };

	json.Unmarshal([]byte(data), cfg); 
}