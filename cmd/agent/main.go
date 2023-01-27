package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alphaonly/harvester/internal/agent"
	c "github.com/alphaonly/harvester/internal/configuration"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ac := c.NewAgentConfiguration()

	(*ac).Update()

	agent.NewAgent(ac).Run(ctx, &http.Client{})

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
