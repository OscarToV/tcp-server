package main

import (
	"log"

	"github.com/OscarToV/tcp-server/configs"
	"github.com/OscarToV/tcp-server/internal/tcpserver"
)

func main() {
	config, err := configs.LoadConfig("../../configs/config.json")
	if err != nil {
		log.Fatalf("Unable to load configuration: %v", err)
	}

	server := tcpserver.New(config.ServerAddress, 5)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
