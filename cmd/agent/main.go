package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alphaonly/harvester/internal/agent"
	conf "github.com/alphaonly/harvester/internal/configuration"
)



func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Configuration parameters from command line
	afc := conf.NewAgentFlagConfiguration()
	//Configuration parameters from environment
	aec := (*conf.NewAgentEnvConfiguration()).Update()

	(*aec).UpdateNotGiven(afc)
	// client := AgentClient{client: &http.Client{}, retries: 10}

	agent.NewAgent(aec).Run(ctx, &http.Client{})

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
