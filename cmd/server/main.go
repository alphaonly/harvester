package main

import (
	"context"
	"log"

	c "github.com/alphaonly/harvester/internal/configuration"
	s "github.com/alphaonly/harvester/internal/server"
	h "github.com/alphaonly/harvester/internal/server/handlers"
	f "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	m "github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
)

func main() {

	var (
		configuration = (*c.NewServerConfiguration()).Update()
		mapStorage    = m.New()
		fileStorage   = f.New(configuration)
		handlers      = h.New(mapStorage)
		server        = s.New(configuration, mapStorage, fileStorage, handlers)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
