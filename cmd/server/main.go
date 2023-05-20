package main

import (
	"context"
	"log"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"

	"github.com/alphaonly/harvester/internal/signchecker"

)

func main() {

	configuration := conf.NewServerConfiguration()
	configuration.UpdateFromEnvironment()
	configuration.UpdateFromFlags()
	fileStorage := fileStor.FileArchive{StoreFile: configuration.StoreFile}

	handlers := &handlers.Handlers{
		MemKeeper: mapstorage.New(),
		Signer:    signchecker.NewSHA256(configuration.Key),
	}
	server := server.New(configuration, fileStorage, handlers)


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
