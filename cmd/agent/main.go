package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alphaonly/harvester/internal/agent"
	conf "github.com/alphaonly/harvester/internal/configuration"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Configuration parameters 
	ac := conf.NewAgentConfiguration()
	ac.UpdateFromEnvironment()
	ac.UpdateFromFlags()

	//http.Client cover
	client := &agent.AgentClient{Client: &http.Client{}, Retries: 10, RetryPause: time.Second * 2}

	agent.NewAgent(ac, client).Run(ctx)

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
