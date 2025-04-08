package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"url_shortener/datastore"
	"url_shortener/http_server"
)

func main() {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM)

	// TEST
	mongoDB := datastore.NewMongoDB()
	if err := mongoDB.Connect(""); err != nil {
		log.Fatalf("Failed to connect to Mongo DB: %v\n", err)
	}
	defer mongoDB.Disconnect()

	server := http_server.NewServer(mongoDB)
	go server.Start(8080) // TODO: move port to config

	<-stopChannel // Wait for the termination signals

	server.Stop()
}
