package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chevron-weather-sensor-simulator/internal/server"
)

func main() {
	serverShutdown := make(chan bool)
	serverDone := make(chan bool)
	go func() {
		log.Println("Starting http server...")
		server.RunWeather(serverShutdown)
		serverDone <- true
	}()

	// Wait for a signal before exiting
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	serverShutdown <- true
	<-serverDone

	log.Println("Done!")
}
