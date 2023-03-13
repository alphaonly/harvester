package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alphaonly/harvester/internal/agent/workerpool"
)

var f workerpool.TypicalJobFunction = func(data any) workerpool.JobResult {
	s := data.(string)
	log.Println("executing function:" + s)
	return workerpool.JobResult{Result: "executed well"}
}

func main() {
	ctx := context.Background()
	wp := workerpool.NewWorkerPool(90, ctx)
	wp.Start()

	for i := 0; i < 80; i++ {
		n := fmt.Sprintf("job #%v", i)

		j := workerpool.Job{Name: n, Func: f}
		time.Sleep(5 * time.Second)
		wp.SendJob(j)
	}
	wp.WaitGroup.Wait()

	log.Println("application is finished")
}
