package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alphaonly/harvester/internal/agent"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ac := agent.Configuration{
		PollInterval:   2,
		ReportInterval: 3, //10
		ServerHost:     "127.0.0.1",
		ServerPort:     "8080",
		UseJSON:        false,
	}

	agent.NewAgent(&ac).Run(ctx, &http.Client{})

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
