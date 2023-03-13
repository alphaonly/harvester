package workerpool

import (
	"context"
	"log"
	"sync"
	"time"
)

type TypicalJobFunction func(data any) JobResult
type JobChannel chan Job
type JobChannelR <-chan Job
type ResultChannel chan JobResult

type Job struct {
	Name string
	Data any
	Func TypicalJobFunction
}

func (j Job) Execute() {
	if j.Func == nil || j.Data == nil {
		log.Fatal("Job execute parameters is nil")
	}
	log.Printf("Executing job %v", j.Name)
	result := j.Func(j.Data)
	log.Printf("Executed job %v with result:%v", j.Name, result)
}

type WorkerPool struct {
	context       context.Context
	Workers       int
	JobChannel    JobChannel
	ResultChannel ResultChannel
	WaitGroup     *sync.WaitGroup
}

func NewWorkerPool(workers int, ctx context.Context) WorkerPool {
	return WorkerPool{
		context:       ctx,
		Workers:       workers,
		JobChannel:    make(JobChannel, workers),
		ResultChannel: make(ResultChannel, workers),
		WaitGroup:     new(sync.WaitGroup),
	}
}

func (wp WorkerPool) Start() {
	wp.WaitGroup.Add(1)
	go func() {
		var w sync.WaitGroup
		for i := 0; i < wp.Workers; i++ {
			w.Add(1)
			go NewWorker(wp.context, wp.JobChannel, i, &w).Start()
		}
		w.Wait()
		wp.WaitGroup.Done()
		log.Println("Worker pool finished")
	}()
}

func (wp WorkerPool) SendJob(job Job) {
	ctxDeadLine, cancel := context.WithDeadline(wp.context, time.Now().Add(time.Second*5))
	defer cancel()

	select {

	case wp.JobChannel <- job:
		{
			log.Printf("Job %v is sent to the worker pool", job.Name)
		}
	case <-ctxDeadLine.Done():
		{
			log.Println("unable to send job to the worker pool in five seconds")
		}
	case <-wp.context.Done():
		{
			log.Println("send job context cancelled")
		}
	}

}

type Worker struct {
	number      int
	context     context.Context
	jobChannelR JobChannelR
	waitGroup   *sync.WaitGroup
	result      JobResult
}

func NewWorker(ctx context.Context, jobChannelR <-chan Job, number int, wg *sync.WaitGroup) Worker {
	return Worker{
		number:      number,
		context:     ctx,
		jobChannelR: jobChannelR,
		waitGroup:   wg,
	}
}

func (w Worker) Start() {
	log.Printf("worker #%v started ", w.number)

	var deadLine, cancel = context.WithDeadline(w.context, time.Now().Add(time.Second*20))

forLabel:
	for {
		select {
		case j := <-w.jobChannelR:
			{
				log.Printf("job \"%v\" received by worker: %v", j.Name, w.number)
				j.Execute()
				log.Printf("job \"%v\" is done by worker: %v", j.Name, w.number)
				deadLine, cancel = context.WithDeadline(w.context, time.Now().Add(time.Second*4))
			}
		case <-w.context.Done():
			{
				log.Printf("worker#%v:application context is done", w.number)
				break forLabel

			}
		case <-deadLine.Done():
			{
				log.Printf("worker#%v: finished as it has done nothing within 4 seconds", w.number)
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
