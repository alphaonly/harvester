package main

import (
	"context"
	"log"

	e "github.com/alphaonly/harvester/internal/environment"
	s "github.com/alphaonly/harvester/internal/server"
	h "github.com/alphaonly/harvester/internal/server/handlers"
	m "github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
)

func main() {

	var (
		mapStorage          = m.New()
		handlers            = h.New(&mapStorage)
		serverConfiguration = (*e.NewServerConfiguration()).Update()
		server              = s.New(handlers, serverConfiguration)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
