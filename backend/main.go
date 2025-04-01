package main

import (
	"os"
	"os/signal"
	"syscall"

	mem_db "url_shortener/db"
	"url_shortener/http_server"
)

func main() {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM)

	dataStore := mem_db.NewMemDB()

	go http_server.Start(8080, dataStore)

	<-stopChannel // Wait for the termination signals

	http_server.Stop()
}
