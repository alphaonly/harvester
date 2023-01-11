package main

import (
	"context"

	"github.com/alphaonly/harvester/internal/agent"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent.Run(ctx)
}
