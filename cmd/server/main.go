package main

import (
	"context"
	"github.com/alphaonly/harvester/cmd/server/handlers"
	"github.com/alphaonly/harvester/cmd/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const serverPort = ":8080"

type Server struct {
}

func NewRouter(ds *storage.DataServer) chi.Router {

	r := chi.NewRouter()
	h := handlers.Handlers{}
	h.SetDataServer(ds)
	//
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.HandleGetMetricFieldList)
		r.Get("/value/{TYPE}/{NAME}", h.HandleGetMetricValue)
		r.Post("/update/{TYPE}/{NAME}/{VALUE}", h.HandlePostMetric)

	})

	return r
}

func (s Server) run(ctx context.Context) error {

	dataServer := storage.DataServer{}.New()

	h := handlers.Handlers{}
	h.SetDataServer(dataServer)

	// маршрутизация запросов обработчику

	router := NewRouter(dataServer)

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
