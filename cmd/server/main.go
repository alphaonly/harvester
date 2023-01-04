package main

import (
	"fmt"
	"github.com/alphaonly/harvester/cmd/server/handlers"
	"github.com/alphaonly/harvester/cmd/server/storage"
	"log"
	"net/http"
)

const (
	// pollInterval   = 2
	// reportInterval = 3 //10

	serverHost = "127.0.0.1"
	serverPort = ":8080"
)

func main() {

	dataServer := storage.DataServer{}.New()

	h := handlers.Handlers{}
	h.SetDataServer(dataServer)

	// маршрутизация запросов обработчику
	http.HandleFunc("/update/", h.HandleMetric)

	log.Fatal(http.ListenAndServe(serverPort, nil))
	fmt.Println(log.Default())
}
