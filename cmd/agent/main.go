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
	 ac:=conf.NewAgentConf(conf.UpdateACFromEnvironment,conf.UpdateACFromFlags)

	// ac := conf.NewAgentConfiguration()
	// ac.UpdateFromEnvironment()
	// ac.UpdateFromFlags()

	//retsty http.Client
	client := resty.New().SetRetryCount(10)

	agent.NewAgent(ac, client).Run(ctx)

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
