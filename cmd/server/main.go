package main

import (
	"context"
	"log"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"

)

func main() {

	configuration := conf.NewServerConfiguration()
	configuration.UpdateFromEnvironment()
	configuration.UpdateFromFlags()

	fileStorage := fileStor.FileArchive{StoreFile: configuration.StoreFile}
	handlers := &handlers.Handlers{MemKeeper: mapstorage.New()}
	server := server.New(configuration, fileStorage, handlers)


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
