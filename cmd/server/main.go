package main

import (
	"context"
	"log"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	mapStor "github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
)

func main() {

	configuration := conf.NewServerEnvConfiguration()
	configuration.UpdateNotGiven(conf.NewServerFlagConfiguration())

	mapStorage := mapStor.New()
	fileStorage := fileStor.New(configuration)
	handlers := handlers.New(mapStorage)
	server := server.New(configuration, mapStorage, fileStorage, handlers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
