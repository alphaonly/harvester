package main_test

import (
	"context"
	"github.com/alphaonly/harvester/internal/server/storage/implementations/mapstorage"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/alphaonly/harvester/internal/server"
	"github.com/alphaonly/harvester/internal/server/handlers"
)

func TestRun(t *testing.T) {

	tests := []struct {
		name string
		URL  string
		want string
	}{
		{
			name: "test#1 - Positive: srv accessible",
			URL:  "http://localhost:8080/check/",
			want: "200 OK",
		},
		{
			name: "test#2 - Negative: srv do not respond",
			URL:  "http://localhost:8080/chek/",
			want: "404 Not Found",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var (
		cfg     = conf.NewServerConf()
		storage = mapstorage.NewStorage()
		hnd     = &handlers.Handlers{}
		srv     = server.New(cfg, storage, hnd, nil)
	)
	cfg.Address = "localhost:8080"
	cfg.Port = "8080"

	go func() {
		err := srv.Run(ctx)
		if err != nil {
			return
		}

	}()

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
			if !assert.Equal(t, tt.want, resp.Status()) {
				t.Error("Server responded unexpectedly")

			}

		})
	}
	cancel()
}
