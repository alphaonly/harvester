package main

import (
	"context"
	s "github.com/alphaonly/harvester/internal/server"
	h "github.com/alphaonly/harvester/internal/server/handlers"
	m "github.com/alphaonly/harvester/internal/server/storage/implementaions/MapStorage"
	"github.com/alphaonly/harvester/internal/server/storage/interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
				var (
					mapStorage          interfaces.Storage = m.New()
					handlers                               = h.New(&mapStorage)
					serverConfiguration                    = s.NewConfiguration("8080")
					server                                 = s.New(handlers, serverConfiguration)
				)
				err := server.Run(ctx)
				if err != nil {
					return
				}
			}()

			if !assert.Equal(t, tt.want, err) {
				t.Error("Server doesn't run")
			}

		})
	}

}
