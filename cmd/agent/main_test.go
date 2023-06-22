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
			name:  "test#1 - Positive:  there are values",
			value: agent.Metrics{},
			want:  true,
		},
		{
			name:  "test#2 - Negative: there are not values",
			want:  false,
		},
	}

	agentConf:=C.NewAgentConf(C.UpdateACFromEnvironment,C.UpdateACFromFlags)

	client := resty.New().SetRetryCount(10)
	a := agent.NewAgent(agentConf, client,nil,nil)

	for i, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			if i!=1{
				go a.Update(ctx, &tt.value)
			}	
			
			time.Sleep(time.Second * 3)
			a.UpdateLocker.Lock()
			fmt.Println(tt.value.PollCount)
			a.UpdateLocker.Unlock()
			if !assert.Equal(t, tt.want, tt.value.PollCount > 0) {
				t.Error("UpdateMemStatsMetrics is not received form runtime values")
			}
		})
	}

}
