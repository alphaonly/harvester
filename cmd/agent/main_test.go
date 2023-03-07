package main_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alphaonly/harvester/internal/agent"
	C "github.com/alphaonly/harvester/internal/configuration"
	"github.com/go-resty/resty/v2"
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
			want:  true,
		},
	}

	agentConf:=C.NewAgentConf(C.UpdateACFromEnvironment,C.UpdateACFromFlags)
	// agentConf := C.NewAgentConfiguration()
	// agentConf.UpdateFromEnvironment()
	// agentConf.UpdateFromFlags()

	client := resty.New().SetRetryCount(10)
	a := agent.NewAgent(agentConf, client)

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			ctx := context.Background()
			ctxMetrics, cancel := context.WithTimeout(ctx, time.Second*2)
			defer cancel()
			go a.Update(ctxMetrics, &tt.value)

			time.Sleep(time.Second * 3)
			fmt.Println(tt.value.PollCount)
			if !assert.Equal(t, tt.want, tt.value.PollCount > 0) {
				t.Error("UpdateMemStatsMetrics is not received form runtime values")
			}
		})
	}

}
