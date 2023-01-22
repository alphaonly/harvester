package main

import (
	"context"
	s "github.com/alphaonly/harvester/internal/server"
	h "github.com/alphaonly/harvester/internal/server/handlers"
	m "github.com/alphaonly/harvester/internal/server/storage/implementaions/MapStorage"
	"github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"log"
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
