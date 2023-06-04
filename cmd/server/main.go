package main

import (
	"context"
	"log"

	"github.com/alphaonly/harvester/internal/common"
	"github.com/alphaonly/harvester/internal/server/crypto"
	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"

	"github.com/alphaonly/harvester/internal/signchecker"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {

	//Build tags
	common.PrintBuildTags(buildVersion, buildDate, buildCommit)
	//Server Configuration
	configuration := conf.NewServerConf(
		conf.UpdateSCFromEnvironment,
		conf.UpdateSCFromFlags,
		conf.UpdateSCFromConfigFile)
	//Storages
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
	//Certificates for decryption
	certManager := crypto.NewRSA(9669, configuration)
	//Handlers
	handlers := &handlers.Handlers{
		Storage: internalStorage,
		Signer:  signchecker.NewSHA256(configuration.Key),
		Conf: conf.ServerConfiguration{
			DatabaseDsn: configuration.DatabaseDsn,
			CryptoKey:   configuration.CryptoKey},
		CertManager: certManager,
	}

	Server := server.New(configuration, externalStorage, handlers, certManager)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := Server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
