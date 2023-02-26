package main_test

import (
	"context"
	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"testing"
	"time"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
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

				sc := conf.NewServerConfiguration()
				sc.UpdateFromEnvironment()
				sc.UpdateFromFlags()
				var storage stor.Storage
				if sc.DatabaseDsn == "" {
					storage = fileStor.FileArchive{StoreFile: sc.StoreFile}
				} else {
					storage = db.NewDBStorage(context.Background(), sc.DatabaseDsn)
				}
				handlers := &handlers.Handlers{}
				server := server.New(sc, storage, handlers)

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
