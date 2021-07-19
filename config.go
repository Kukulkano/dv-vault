package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ConnectionString string `json:"connectionString"`
	MaxConnections   int    `json:"maxConnections"`
	IP               string `json:"listenIP"`
	Port             int    `json:"listenPort"`
	LetsEncrypt      int    `json:"useLetsEncrypt"`
	Domain           string `json:"domain"`
	DebugMode        int    `json:"debugMode"`
	IPExtractor      string `json:"IPExtractor"`
	DisableIPCheck   int    `json:"disableIPCheck"`
	CORSDomains      string `json:"CORSDomains"`
	RunAs            string `json:"runAs"`
	CertFolder       string `json:"certFolder"`
}

var cfg Configuration

func loadConfig() {
	fmt.Print("Open config.json… ")
	file, err := os.Open("config.json")
	if err != nil {
		panic("Missing config.json")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	cfg = Configuration{}
	err = decoder.Decode(&cfg)
	if err != nil {
		panic("Invalid config.json")
	}
	fmt.Println("Done")
}
