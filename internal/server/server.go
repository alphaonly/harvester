package server

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/crypto"
	"github.com/alphaonly/harvester/internal/server/handlers"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	configuration   *conf.ServerConfiguration
	InternalStorage stor.Storage
	ExternalStorage stor.Storage
	handlers        *handlers.Handlers
	httpServer      *http.Server
	crypto          cryptoCommon.CertificateManager
}

func NewConfiguration(serverPort string) *Configuration {
	return &Configuration{serverPort: ":" + serverPort}
}

func New(
	configuration *conf.ServerConfiguration,
	ExStorage stor.Storage,
	handlers *handlers.Handlers,
	certificate cryptoCommon.CertificateManager) (server Server) {

	return Server{
		configuration:   configuration,
		InternalStorage: handlers.Storage,
		ExternalStorage: ExStorage,
		handlers:        handlers,
		crypto:          certificate,
	}
}

func (s Server) ListenData(ctx context.Context) {

	var err error
	if s.configuration.CryptoKey == "" {
		err = s.httpServer.ListenAndServe()
	} else if s.configuration.EnableHTTPS {
		//"/home/asus/goProjects/harvester/cmd/server/rsa/private/cert.rsa"
		err = s.httpServer.ListenAndServeTLS(s.configuration.CryptoKey, s.configuration.CryptoKey)
	}

	if err != nil {
		log.Println(err)
	}
}

func (s *Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику
	s.httpServer = &http.Server{
		Addr:    s.configuration.Address,
		Handler: s.handlers.NewRouter(),
	}

	s.restoreData(ctx, s.ExternalStorage)

	s.CheckCertificateFile(cryptoCommon.Private())

	go s.ListenData(ctx)
	go s.ParkData(ctx, s.ExternalStorage)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := s.Shutdown(ctx)

	return err
}
func (s Server) Shutdown(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	err := s.httpServer.Shutdown(ctx)
	log.Println("Server shutdown")
	return err
}

func (s Server) restoreData(ctx context.Context, storageFrom stor.Storage) {
	if storageFrom == nil {
		log.Println("external storage  not initiated ")
		return
	}
	if s.configuration.Restore {
		mvList, err := storageFrom.GetAllMetrics(ctx)

		if err != nil {
			log.Println("cannot initially read metrics from file storage")
			return
		}
		if len(*mvList) == 0 {
			log.Println("file storage is empty, nothing to recover")
			return
		}

		err = s.InternalStorage.SaveAllMetrics(ctx, mvList)
		if err != nil {
			log.Fatal("cannot save metrics to internal storage")
		}

	}

}

func (s Server) ParkData(ctx context.Context, storageTo stor.Storage) {

	if storageTo == nil {
		return
	}
	if s.handlers.Storage == storageTo {
		log.Fatal("a try to save to it is own")
		return
	}

	ticker := time.NewTicker(time.Duration(s.configuration.StoreInterval))

	defer ticker.Stop()

DoItAgain:
	select {

	case <-ticker.C:
		{
			mvList, err := s.InternalStorage.GetAllMetrics(ctx)

			if err != nil {
				log.Fatal("cannot read metrics from internal storage")
			}
			if mvList == nil {
				log.Println("read insufficient, internal storage empty")
			} else if len(*mvList) == 0 {
				log.Println("internal storage is empty, nothing to save to file")
			} else {
				err = storageTo.SaveAllMetrics(ctx, mvList)
				if err != nil {
					log.Fatal("cannot write metrics to file storage:" + err.Error())

				}
				log.Println("saved to file")
			}

		}
	case <-ctx.Done():
		return

	}
	goto DoItAgain
}

func (s *Server) CheckCertificateFile(dataType cryptoCommon.DataType) error {
	if !s.configuration.EnableHTTPS {
		return nil
	}
	if s.configuration.CryptoKey == "" {
		//generate certificates in test folder
		genCryptoFiles()
		log.Println("path to rsa files is not defined, new rsa files were generated in /rsa/ folder")
		return nil
	}

	//Reading file with rsa key from os
	file, err := os.OpenFile(s.configuration.CryptoKey, os.O_RDONLY, 0777)
	if err != nil {
		log.Printf("error:file %v  is not read", file)
		return err
	}
	//put data to read buffer
	buf := bufio.NewReader(file)
	rsa := crypto.RSA{}
	_, err = rsa.Receive(dataType, buf)
	if err != nil {
		log.Println("error:private rsa is not read")
		return err
	}
	return nil

}

func genCryptoFiles() {
	prefPath, _ := os.Getwd()
	fileMap := cryptoCommon.DataTypeMap{
		prefPath + "/rsa/cert.rsa":    cryptoCommon.Certificate,
		prefPath + "/rsa/public.rsa":  cryptoCommon.Public,
		prefPath + "/rsa/private.rsa": cryptoCommon.Private}

	rsa := crypto.NewRSA(6996, &conf.ServerConfiguration{EnableHTTPS: true})
	for filename, dataTypeFunc := range fileMap {
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
		fileWriter := bufio.NewWriter(file)
		rsa.Send(dataTypeFunc(), fileWriter)
		err = fileWriter.Flush()
		if err != nil {
			log.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

}
