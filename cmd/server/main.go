package main

import (
	"context"
	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
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

	var (
		ExternalStorage stor.Storage
		InternalStorage stor.Storage
	)
	if configuration.DatabaseDsn == "" {
		ExternalStorage = fileStor.FileArchive{StoreFile: configuration.StoreFile}
		InternalStorage = mapstorage.New()
	} else {
		ExternalStorage = nil
		InternalStorage = db.NewDBStorage(context.Background(), configuration.DatabaseDsn)
	}

	_handlers := &handlers.Handlers{
		Storage: InternalStorage,
		Signer:  signchecker.NewSHA256(configuration.Key),
		Conf:    conf.ServerConfiguration{DatabaseDsn: configuration.DatabaseDsn},
	}

	_server := server.New(configuration, ExternalStorage, _handlers)


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := _server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
