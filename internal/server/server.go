package server

import (
	"bufio"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
	pb "github.com/alphaonly/harvester/internal/common/grpc/proto"
	"github.com/alphaonly/harvester/internal/common/logging"
	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/crypto"
	"github.com/alphaonly/harvester/internal/server/handlers"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"google.golang.org/grpc"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	cfg             *conf.ServerConfiguration
	InternalStorage stor.Storage
	ExternalStorage stor.Storage
	handlers        *handlers.Handlers
	httpServer      *http.Server
	crypto          cryptoCommon.ServerCertificateManager
	grpcService     pb.ServiceServer
	grpcServer      *grpc.Server
}

func New(
	configuration *conf.ServerConfiguration,
	ExStorage stor.Storage,
	handlers *handlers.Handlers,
	certificate cryptoCommon.ServerCertificateManager,
	grpcService pb.ServiceServer) (server Server) {

	return Server{
		cfg:             configuration,
		InternalStorage: handlers.Storage,
		ExternalStorage: ExStorage,
		handlers:        handlers,
		crypto:          certificate,
		grpcService:     grpcService,
	}
}

// ListenGRPC - starts listening gRPC server
func (s *Server) ListenGRPC(ctx context.Context) {

	//check necessary data
	if s.cfg.GRPCPort == "" || s.grpcService != nil {
		return
	}
	//listener configuration
	listener, err := net.Listen("tcp", ":"+s.cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	// create grpc
	s.grpcServer = grpc.NewServer()
	// register service
	pb.RegisterServiceServer(s.grpcServer, s.grpcService)
	log.Println("Start gRPC server")
	if err := s.grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (s Server) listenHTTP(ctx context.Context) {

	err := s.httpServer.ListenAndServe()
	logging.LogFatal(err)
}

func (s *Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику
	s.httpServer = &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.handlers.NewRouter(),
	}

	s.restoreData(ctx, s.ExternalStorage)

	go s.listenHTTP(ctx)
	go s.ListenGRPC(ctx)

	go s.parkData(ctx, s.ExternalStorage)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	shutdown := func(s *Server) error {
		//Graceful gRPC shutdown
		s.grpcServer.GracefulStop()
		//Graceful http shutdown
		err := s.Shutdown(ctx)
		return err
	}
	var err error
	select {
	case <-osSignal:
		err = shutdown(s)
	case <-ctx.Done():
		err = shutdown(s)
	}
	return err
}

// shutdown
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
	if s.cfg.Restore {
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

func (s Server) parkData(ctx context.Context, storageTo stor.Storage) {

	if storageTo == nil {
		return
	}
	if s.handlers.Storage == storageTo {
		log.Fatal("a try to save to it is own")
		return
	}

	ticker := time.NewTicker(time.Duration(s.cfg.StoreInterval))

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

func (s *Server) CheckCertificateFile(dataType cryptoCommon.KeyType) error {

	if s.cfg.CryptoKey == "" {
		//generate certificates in test folder
		crypto.MakeCryptoFiles("rsa", s.cfg, nil)
		log.Println("path to rsa files is not defined, new rsa files were generated in /rsa/ folder")
		return nil
	}

	//Reading file with rsa key from os
	file, err := os.OpenFile(s.cfg.CryptoKey, os.O_RDONLY, 0777)
	if err != nil {
		log.Printf("error:file %v  is not read", file)
		return err
	}
	//put data to read buffer
	buf := bufio.NewReader(file)
	rs := crypto.RSA{}
	rs.Receive(dataType, buf)
	if rs.Error() != nil {
		log.Println("error:private rsa is not read")
		return rs.Error()
	}
	return nil

}
