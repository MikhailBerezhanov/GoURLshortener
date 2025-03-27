package main

import (
	"os"
	"os/signal"
	"syscall"

	"url_shortener/http_server"
)

func main() {
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM)

	go http_server.Start(8080)

	<-stopChannel // Wair for the termination signal

	http_server.Stop()
}
