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

	// mj := agentjson.NewMetricJSON("PollCount", "counter", nil)

	// baseUrl := url.URL{
	// 	Scheme: "http",
	// 	Host:   "127.0.0.1:8080",
	// }
	// rj := mj.GetMetricJSON(&baseUrl, "PollCount", "counter")
	// fmt.Println(rj)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Configuration parameters from command line
	afc := c.NewAgentFlagConfiguration()
	//Configuration parameters from environment
	aec := (*c.NewAgentEnvConfiguration()).Update()

	(*aec).UpdateNotGiven(afc)

	agent.NewAgent(aec).Run(ctx, &http.Client{})

	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	<-channel
	log.Print("Agent shutdown")

}
