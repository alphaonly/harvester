package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alphaonly/harvester/internal/server/handlers"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	handlers      *handlers.Handlers
	configuration *Configuration
}

func NewConfiguration(serverPort string) *Configuration {
	return &Configuration{serverPort: ":" + serverPort}
}

func New(handlers *handlers.Handlers, configuration *Configuration) (server Server) {
	return Server{
		handlers:      handlers,
		configuration: configuration,
	}
}

func (s Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику

	server := http.Server{
		Addr: s.configuration.serverPort,
	}

	go func() {
		err := http.ListenAndServe(s.configuration.serverPort, s.handlers.NewRouter())
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Setting up signal capturing
	channelInt := make(chan os.Signal, 1)
	signal.Notify(channelInt, os.Interrupt)

	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*5)
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

	err := server.Shutdown(ctx2)
	log.Println("Server shutdown")
	return err
}
