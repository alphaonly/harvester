package main

import (
	"context"
	"testing"
	"time"

	"github.com/alphaonly/harvester/internal/agent"
	C "github.com/alphaonly/harvester/internal/configuration"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {

	tests := []struct {
		name  string
		value agent.Metrics
		want  bool
	}{
		{
			name:  "test#1 - Positive: are there values?",
			value: agent.Metrics{},
			want:  false,
		},
	}

	ac := C.NewAgentEnvConfiguration()
	(*ac).Update()

	a := agent.NewAgent(ac)

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			ctx := context.Background()
			ctxMetrics, cancel := context.WithTimeout(ctx, time.Second*3)
			defer cancel()
			go a.Update(ctxMetrics, &tt.value)

			time.Sleep(time.Second * 4)

			if !assert.Equal(t, tt.want, tt.value.PollCount > 0) {
				t.Error("UpdateMemStatsMetrics is not received form runtime values")
			}
		})
	}

}
