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

	//Configuration parameters from command line
	afc := conf.NewAgentFlagConfiguration()
	//Configuration parameters from environment
	aec := conf.NewAgentEnvConfiguration()

	aec.UpdateNotGiven(afc)

	//http.Client cover
	client := &agent.AgentClient{Client: &http.Client{}, Retries: 10, RetryPause: time.Second * 2}

	agent.NewAgent(aec, client).Run(ctx)

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
