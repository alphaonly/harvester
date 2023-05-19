package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	C "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/handlers"
	S "github.com/alphaonly/harvester/internal/server/storage/interfaces"
)

type Configuration struct {
	serverPort string
}

type Server struct {
	configuration *C.Configuration
	memKeeper     *S.Storage
	archive       *S.Storage
	handlers      *handlers.Handlers
}

func NewConfiguration(serverPort string) *Configuration {
	return &Configuration{serverPort: ":" + serverPort}
}

func New(configuration *C.Configuration, memKeeper *S.Storage, archive *S.Storage, handlers *handlers.Handlers) (server Server) {
	return Server{
		configuration: configuration,
		memKeeper:     memKeeper,
		archive:       archive,
		handlers:      handlers,
	}
}

func (s Server) ListenData(ctx context.Context) {
	err := http.ListenAndServe(":"+(*s.configuration).Get("S_PORT"), s.handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}

func (s Server) Run(ctx context.Context) error {

	// маршрутизация запросов обработчику

	server := http.Server{
		Addr: ":" + (*s.configuration).Get("S_PORT"),
	}

	s.restoreData(ctx, s.archive)

	go s.ListenData(ctx)
	go s.ParkData(ctx, s.archive)

	// Setting up signal capturing
	channelInt := make(chan os.Signal, 1)
	signal.Notify(channelInt, os.Interrupt)

	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*5)
	select {
	case <-channelInt:
		{
			cancel()
		}
	case <-ctx.Done():
		{
			cancel()
		}
	}

	err := server.Shutdown(ctx2)
	log.Println("Server shutdown")
	return err
}

func (s Server) restoreData(ctx context.Context, storageFrom *S.Storage) {

	if (*s.configuration).GetBool("RESTORE") {
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

func (s Server) ParkData(ctx context.Context, storageTo *S.Storage) {

	if s.handlers.MemKeeper == storageTo {
		log.Fatal("a try to save to it is own")
		return
	}

	ticker := time.NewTicker(time.Duration((*s.configuration).GetInt("STORE_INTERVAL")) * (time.Second))
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
