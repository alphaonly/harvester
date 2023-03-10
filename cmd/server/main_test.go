package main_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	db "github.com/alphaonly/harvester/internal/server/storage/implementations/dbstorage"
	fileStor "github.com/alphaonly/harvester/internal/server/storage/implementations/filestorage"
	stor "github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name string

		want string
	}{
		{
			name: "test#1 - Positive: server works",
			want: "200 OK",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(tst *testing.T) {

			// var err error

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			sc := conf.NewServerConf(conf.UpdateSCFromEnvironment, conf.UpdateSCFromFlags)

			var storage stor.Storage
			if sc.DatabaseDsn == "" {
				storage = fileStor.FileArchive{StoreFile: sc.StoreFile}
			} else {
				storage = db.NewDBStorage(context.Background(), sc.DatabaseDsn)
			}
			handlers := &handlers.Handlers{}
			server := server.New(sc, storage, handlers)

			go func() {
				err := server.Run(ctx)
				if err != nil {
					return
				}
			}()

			time.Sleep(time.Second * 6)

			keys := make(map[string]string)
			keys["Content-Type"] = "plain/text"
			keys["Accept"] = "plain/text"

			client := resty.New()

			r := client.R().
				SetHeaders(keys)
			resp, err := r.Get("http://localhost:8080/check/")
			if err != nil {
				log.Fatalf("send new request error:%v", err)
			}
			fmt.Println(resp.Status())
			if !assert.Equal(t, tt.want, resp.Status()) {
				t.Error("Server didn't respond well")
			}

		})
	}

}
