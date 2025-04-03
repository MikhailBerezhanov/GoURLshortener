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
	if _, err := datastore.ConnectMongoDb(""); err != nil {
		log.Printf("Failed to connect to Mongo DB: %v\n", err)
	}

	dataStore := datastore.NewMemDB()

	go http_server.Start(8080, dataStore)

	<-stopChannel // Wait for the termination signals

	http_server.Stop()
}
