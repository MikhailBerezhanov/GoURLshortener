package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"url_shortener/datastore"
	"url_shortener/http_server"
	"url_shortener/url"
)

func main() {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM)

	// TEST
	mongoDB := datastore.NewMongoDB()
	if err := mongoDB.Connect(""); err != nil {
		log.Printf("Failed to connect to Mongo DB: %v\n", err)
	}
	defer mongoDB.Disconnect()

	var dataStore url.RecordStore = mongoDB //datastore.NewMemDB()

	go http_server.Start(8080, dataStore)

	<-stopChannel // Wait for the termination signals

	http_server.Stop()
}
