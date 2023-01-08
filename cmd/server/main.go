package main

import (
	"context"
	"github.com/alphaonly/harvester/cmd/server/handlers"
	"github.com/alphaonly/harvester/cmd/server/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const serverPort = ":8080"

type Server struct {
}

func (s Server) run(ctx context.Context) error {

	dataServer := storage.DataServer{}.New()

	h := handlers.Handlers{}
	h.SetDataServer(dataServer)

	// маршрутизация запросов обработчику

	router := handlers.NewRouter(dataServer)

	//http.HandleFunc("/update/", h.HandlePostMetric)

	server := http.Server{
		Addr: serverPort,
	}
	var err error

	go func() {
		err = http.ListenAndServe(serverPort, router)
	}()

	// Setting up signal capturing
	channelInt := make(chan os.Signal, 1)
	signal.Notify(channelInt, os.Interrupt)

	//<-channelInt

	var (
		ctx2   context.Context
		cancel context.CancelFunc
	)

	ctx2, cancel = context.WithTimeout(context.Background(), time.Second*5)
	select {
	case <-channelInt:
		{

			cancel()
		}
	case <-ctx.Done():
		{

			cancel()
		}
	}

	err = server.Shutdown(ctx2)

	return err
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := Server{}.run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
