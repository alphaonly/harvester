package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/handlers"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	configuration *conf.ServerEnvConfiguration
	memKeeper     *stor.Storage
	archive       *stor.Storage
	handlers      *handlers.Handlers
}

func NewConfiguration(serverPort string) *Configuration {
	return &Configuration{serverPort: ":" + serverPort}
}

func New(configuration *conf.ServerEnvConfiguration, memKeeper *stor.Storage, archive *stor.Storage, handlers *handlers.Handlers) (server Server) {
	return Server{
		configuration: configuration,
		memKeeper:     memKeeper,
		archive:       archive,
		handlers:      handlers,
	}
}

func (s Server) ListenData(ctx context.Context) {
	err := http.ListenAndServe(s.configuration.Get("ADDRESS"), s.handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}

func (s Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику

	server := http.Server{
		Addr: s.configuration.Get("ADDRESS"),
	}

	s.restoreData(ctx, s.archive)

	go s.ListenData(ctx)
	go s.ParkData(ctx, s.archive)

	// Setting up signal capturing
	channelInt := make(chan os.Signal, 1)
	signal.Notify(channelInt, os.Interrupt)

	select {
	case <-channelInt:
		{
			goto shutdown
		}
	case <-ctx.Done():
		{
			goto shutdown
		}
	}

shutdown:
	time.Sleep(time.Second * 2)
	err := server.Shutdown(ctx)
	log.Println("Server shutdown")
	return err
}

func (s Server) restoreData(ctx context.Context, storageFrom *stor.Storage) {

	if s.configuration.GetBool("RESTORE") {
		mvList, err := (*storageFrom).GetAllMetrics(ctx)
		if err != nil {
			log.Println("cannot initially read metrics from file storage")
			return
		}
		if len((*mvList)) == 0 {
			log.Println("file storage is empty, nothing to recover")
			return
		}

		err = (*s.memKeeper).SaveAllMetrics(ctx, mvList)
		if err != nil {
			log.Fatal("cannot save metrics to internal storage")
		}

	}

}

func (s Server) ParkData(ctx context.Context, storageTo *stor.Storage) {

	if s.handlers.MemKeeper == storageTo {
		log.Fatal("a try to save to it is own")
		return
	}

	ticker := time.NewTicker(time.Duration(s.configuration.GetInt("STORE_INTERVAL")) * (time.Second))
	defer ticker.Stop()

DoitAgain:
	select {

	case <-ticker.C:
		{

			mvList, err := (*s.memKeeper).GetAllMetrics(ctx)
			if err != nil {
				log.Fatal("cannot read metrics from internal storage")
			}
			if mvList == nil {
				log.Println("read insufficient, internal storage empty")
			} else if len((*mvList)) == 0 {
				log.Println("internal storage is empty, nothing to save to file")
			} else {
				err = (*storageTo).SaveAllMetrics(ctx, mvList)
				if err != nil {
					log.Fatal("cannot write metrics to file storage")
				}
				log.Println("saved to file")
			}

		}
	case <-ctx.Done():
		return

	}
	goto DoitAgain
}
