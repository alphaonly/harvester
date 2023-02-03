package main_test

import (
	"context"
	"testing"
	"time"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	mapStor "github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name string

		want error
	}{
		{
			name: "test#1 - Positive: server works",
			want: nil,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(tst *testing.T) {

			var err error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			go func() {

				configuration := conf.NewServerEnvConfiguration()
				mapStorage := mapStor.New()
				archive := filestorage.New(configuration)
				handlers := handlers.New(mapStorage)
				server := server.New(configuration, mapStorage, archive, handlers)

				err := server.Run(ctx)
				if err != nil {
					return
				}
			}()

			time.Sleep(time.Second * 10)

			if !assert.Equal(t, tt.want, err) {
				t.Error("Server doesn't run")
			}

		})
	}

}
