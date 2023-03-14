package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/alphaonly/harvester/internal/agent"
	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/go-resty/resty/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Configuration parameters
	ac := conf.NewAgentConf(conf.UpdateACFromEnvironment, conf.UpdateACFromFlags)
	//resty client
	client := resty.New().SetRetryCount(10)
	//Run agent (with context)
	agent.NewAgent(ac, client).Run(ctx)
	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	select {
	case <-channel:
		log.Print("Agent shutdown by os signal")
	case <-ctx.Done():
		log.Print("Agent shutdown by cancelled context")
	}
}
