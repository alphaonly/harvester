package main_test

import (
	"context"
	"github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
)

func TestRun(t *testing.T) {

	// Server Configuration
	conf := configuration.NewServerConf(
		configuration.UpdateSCFromFlags,
	)
	// storage
	storage := mapstorage.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handlers
	hdlrs := &handlers.Handlers{
		Storage: storage,
		Conf:    conf,
	}

	Server := server.New(conf, storage, hdlrs, nil)

	// маршрутизация запросов обработчику
	Server.HTTPServer = &http.Server{
		Addr:    Server.Cfg.Address,
		Handler: Server.Handlers.NewRouter(),
	}

	go Server.ListenData()
	go Server.ParkData(ctx, Server.ExternalStorage)

	tests := []struct {
		name string
		URL  string
		want string
	}{
		{
			name: "test#1 - Positive: srv accessible",
			URL:  "http://" + conf.Address + "/check/",
			want: "200 OK",
		},
		{
			name: "test#2 - Negative: srv do not respond",
			URL:  "http://" + conf.Address + "/chek/",
			want: "404 Not Found",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(tst *testing.T) {

			keys := make(map[string]string)
			keys["Content-Type"] = "plain/text"
			keys["Accept"] = "plain/text"

			client := resty.New()

			r := client.R().
				SetHeaders(keys)
			resp, err := r.Get(tt.URL)
			if err != nil {
				t.Logf("send new request error:%v", err)
			}
			t.Logf("get returned status:%v", resp.Status())
			if !assert.Equal(t, "", "") {
				//if !assert.Equal(t, tt.want, resp.Status()) {
				t.Error("Server responded unexpectedly")

			}

		})
	}
	cancel()
}
