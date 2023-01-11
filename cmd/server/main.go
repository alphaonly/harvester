package main

import (
	"context"
	"log"

	"github.com/alphaonly/harvester/internal/server"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := server.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
