package main

import (
	"context"
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
				err = Server{}.run(ctx)
			}()

			if !assert.Equal(t, tt.want, err) {
				t.Error("Server doesn't run")
			}

		})
	}

}
