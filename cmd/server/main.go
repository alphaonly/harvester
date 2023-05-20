package main

import (
	"context"
	"log"

	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"

	"github.com/alphaonly/harvester/internal/signchecker"

)

func main() {

	configuration := conf.NewServerConf(conf.UpdateSCFromEnvironment, conf.UpdateSCFromFlags)
	// configuration.UpdateFromEnvironment()
	// configuration.UpdateFromFlags()

	var (
		externalStorage stor.Storage
		internalStorage stor.Storage
	)
	externalStorage = fileStor.FileArchive{StoreFile: configuration.StoreFile}
	internalStorage = mapstorage.New()

	if configuration.DatabaseDsn != "" {
		externalStorage = nil
		internalStorage = db.NewDBStorage(context.Background(), configuration.DatabaseDsn)
	}

	handlers := &handlers.Handlers{
		Storage: internalStorage,
		Signer:  signchecker.NewSHA256(configuration.Key),
		Conf:    conf.ServerConfiguration{DatabaseDsn: configuration.DatabaseDsn},
	}

	metricsServer := server.New(configuration, externalStorage, handlers)


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := metricsServer.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
