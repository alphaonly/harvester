package main

import (
	"context"
	"log"

	s "github.com/alphaonly/harvester/internal/server"
	h "github.com/alphaonly/harvester/internal/server/handlers"
	m "github.com/alphaonly/harvester/internal/server/storage/implementations/MapStorage"
	"github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

func main() {

	var (
		mapStorage          interfaces.Storage = m.New()
		handlers                               = h.New(&mapStorage)
		serverConfiguration                    = s.NewConfiguration("8080")
		server                                 = s.New(handlers, serverConfiguration)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
