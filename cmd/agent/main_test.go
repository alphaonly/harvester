package main

import (
	"context"
	"github.com/alphaonly/harvester/internal/agent"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestUpdateMemStatsMetrics(t *testing.T) {

	tests := []struct {
		name  string
		value agent.Metrics
		want  bool
	}{
		{
			name: "test#1 - Positive: are there values?",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			ctx := context.Background()
			//ctxMetrics, cancel := context.WithCancel(ctx)
			ctxMetrics, cancel := context.WithTimeout(ctx, time.Second*3)
			defer cancel()
			go agent.Agent{}.Update(ctxMetrics, &tt.value)

			time.Sleep(time.Second * 4)

			if !assert.Equal(t, tt.want, reflect.ValueOf(tt.value).IsZero()) {
				t.Error("UpdateMemStatsMetrics is not received form runtime values")
			}
		})
	}

}
