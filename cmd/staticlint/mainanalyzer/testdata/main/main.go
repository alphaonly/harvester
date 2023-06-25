package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {

	channel := make(chan os.Signal, 1)
	//Graceful shutdown
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown by os signal")

}
