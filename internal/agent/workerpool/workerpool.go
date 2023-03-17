package workerpool

import (
	"context"
	"log"
	"sync"
	"time"
)

type Function[T any] struct {
	Value T
	F func(data T) JobResult
}
func (f *Function[T]) NewFunctionData(value T){
	f.Value=value
}

type TypicalJobFunction[T any] func(data T) JobResult
type JobChannel[T any] chan Job[T]
type JobChannelR[T any] <-chan Job[T]
type ResultChannel chan JobResult

type Job[T any] struct {
	Name string
	Data T
	Func TypicalJobFunction[T]
}

func (j Job[T]) Execute() {
	if j.Func == nil {
		//|| j.Data == nil

		log.Fatal("Job execute parameters is nil")
	}
	log.Printf("Executing job %v", j.Name)
	result := j.Func(j.Data)
	log.Printf("Executed job %v with result:%v", j.Name, result)
}

type WorkerPool[T any] struct {
	Workers       int
	JobChannel    JobChannel[T]
	ResultChannel ResultChannel
	WaitGroup     *sync.WaitGroup
}

func NewWorkerPool[T any](workers int) WorkerPool[T] {
	return WorkerPool[T]{
		Workers:       workers,
		JobChannel:    make(JobChannel[T], workers),
		ResultChannel: make(ResultChannel, workers),
		WaitGroup:     new(sync.WaitGroup),
	}
}

func (wp WorkerPool[T]) Start(ctx context.Context) {
	wp.WaitGroup.Add(1)
	go func() {
		var w sync.WaitGroup
		for i := 0; i < wp.Workers; i++ {
			w.Add(1)
			go NewWorker(wp.JobChannel, i, &w).Start(ctx)
		}
		w.Wait()
		wp.WaitGroup.Done()
		log.Println("Worker pool finished")
	}()
}

func (wp WorkerPool[T]) SendJob(ctx context.Context, job Job[T]) {
	ctxDeadLine, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()

	select {
	case wp.JobChannel <- job:
		{
			log.Printf("Job %v has been sent to the worker pool", job.Name)
		}
	case <-ctxDeadLine.Done():
		{
			log.Println("send job cancelled as it has been unable to send job to the worker pool for five seconds")
		}
	case <-ctx.Done():
		{
			log.Println("send job cancelled by application context")
		}
	}
}

type Worker[T any] struct {
	number      int
	jobChannelR JobChannelR[T]
	waitGroup   *sync.WaitGroup
	result      JobResult
}

func NewWorker[T any](jobChannelR <-chan Job[T], number int, wg *sync.WaitGroup) Worker[T] {
	return Worker[T]{
		number:      number,
		jobChannelR: jobChannelR,
		waitGroup:   wg,
	}
}

func (w Worker[T]) Start(ctx context.Context) {
	log.Printf("worker #%v started ", w.number)

	var deadLine, cancel = context.WithDeadline(ctx, time.Now().Add(time.Second*20))

forLabel:
	for {
		select {
		case j := <-w.jobChannelR:
			{
				log.Printf("job \"%v\" received by worker: %v", j.Name, w.number)
				j.Execute()
				log.Printf("job \"%v\" is done by worker: %v", j.Name, w.number)
				deadLine, cancel = context.WithDeadline(ctx, time.Now().Add(time.Second*4))
			}
		case <-ctx.Done():
			{
				log.Printf("worker#%v:application context is done", w.number)
				break forLabel
			}
		case <-deadLine.Done():
			{
				log.Printf("worker#%v: finished as it has done nothing within 20 seconds", w.number)
				break forLabel
			}
		}
	}
	log.Printf("worker #%v finished", w.number)
	w.waitGroup.Done()
	cancel()
}

type JobResult struct {
	Result string
}
