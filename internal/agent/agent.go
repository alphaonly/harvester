package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/alphaonly/harvester/internal/agent/workerpool"
	"github.com/alphaonly/harvester/internal/schema"
	"github.com/alphaonly/harvester/internal/server/compression"
	sign "github.com/alphaonly/harvester/internal/signchecker"
	"github.com/go-resty/resty/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"strconv"
	"time"

	"math/rand"
	"net/url"

	conf "github.com/alphaonly/harvester/internal/configuration"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	RandomValue   Gauge

	PollCount Counter

	//GOPSUtil
	TotalMemory     Gauge
	FreeMemory      Gauge
	CPUutilization1 Gauge
}

type Agent struct {
	Configuration *conf.AgentConfiguration
	baseURL       url.URL
	Client        *resty.Client
	Signer        sign.Signer
	UpdateLocker  *sync.RWMutex
}

func NewAgent(c *conf.AgentConfiguration, client *resty.Client) Agent {


	return Agent{
		Configuration: c,
		baseURL: url.URL{

			Scheme: c.Scheme,
			Host:   c.Address,
		},
		Client:       client,
		Signer:       sign.NewSHA256(c.Key),
		UpdateLocker: new(sync.RWMutex),
	}
}

func AddCounterData(common sendData, val Counter, name string, data map[*sendData]bool) {
	URL := common.url.
		JoinPath("counter").
		JoinPath(name).
		JoinPath(strconv.FormatUint(uint64(val), 10)) //value float

	sd := sendData{
		url:  URL,
		keys: common.keys,
	}
	data[&sd] = true

}
func AddGaugeData(common sendData, val Gauge, name string, data map[*sendData]bool) {

	URL := common.url.
		JoinPath("gauge").
		JoinPath(name).
		JoinPath(strconv.FormatFloat(float64(val), 'E', -1, 64)) //value float

	// empty := bytes.NewBufferString(URL.String()).Bytes()

	sd := sendData{
		url:  URL,
		keys: common.keys,
		// body: &empty, //need to tranfser something
	}
	data[&sd] = true

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func AddGaugeDataJSONToBatch(common *sendData, val Gauge, name string) {
	if common.JSONBatchBody == nil {
		common.JSONBatchBody = new([]schema.MetricsJSON)
	}

	v := float64(val)

	mj := schema.MetricsJSON{
		ID:    name,
		MType: "gauge",
		Value: &v,
	}
	//Вычисляем hash и помещаем в mj.Hash
	err := common.signer.Sign(&mj)
	logFatal(err)

	*common.JSONBatchBody = append(*common.JSONBatchBody, mj)
}

func AddCounterDataJSONToBatch(common *sendData, val Counter, name string) {

	if common.JSONBatchBody == nil {
		common.JSONBatchBody = new([]schema.MetricsJSON)
	}
	v := int64(val)

	mj := schema.MetricsJSON{
		ID:    name,
		MType: "counter",
		Delta: &v,
	}
	//Вычисляем hash и помещаем в mj.Hash
	err := common.signer.Sign(&mj)
	logFatal(err)

	*common.JSONBatchBody = append(*common.JSONBatchBody, mj)
}

func AddGaugeDataJSON(common sendData, val Gauge, name string, data map[*sendData]bool) {
	v := float64(val)

	mj := schema.MetricsJSON{
		ID:    name,
		MType: "gauge",
		Value: &v,
	}
	//Вычисляем hash и помещаем в mj.Hash
	err := common.signer.Sign(&mj)
	logFatal(err)

	sd := sendData{
		url:      common.url,
		keys:     common.keys,
		JSONBody: &mj,
	}
	data[&sd] = true

}
func AddCounterDataJSON(common sendData, val Counter, name string, data map[*sendData]bool) {
	v := int64(val)

	mj := schema.MetricsJSON{
		ID:    name,
		MType: "counter",
		Delta: &v,
	}
	if val == -1 {
		//API /value check
		mj = schema.MetricsJSON{
			ID:    name,
			MType: "counter",
		}
	}

	//Вычисляем hash и помещаем в mj.Hash
	err := common.signer.Sign(&mj)
	logFatal(err)

	sd := sendData{
		url:      common.url,
		keys:     common.keys,
		JSONBody: &mj,
	}
	data[&sd] = true

}

type HeaderKeys map[string]string
type sendData struct {
	url            *url.URL
	keys           HeaderKeys
	JSONBody       *schema.MetricsJSON
	JSONBatchBody  *[]schema.MetricsJSON
	compressedBody *[]byte
	signer         sign.Signer
}

func (sd sendData) SendData(client *resty.Client) error {

	//a resty attempt

	r := client.R().
		SetHeaders(sd.keys)

	if sd.JSONBody != nil {
		r.SetBody(sd.JSONBody)
	} else if sd.JSONBatchBody != nil {
		r.SetBody(sd.JSONBatchBody)
	} else {
		return errors.New("both bodies is nil")
	}

	resp, err := r.
		Post(sd.url.String())
	if err != nil {
		log.Fatalf("send new request error:%v", err)
	}
	log.Println("agent:response status from server:" + resp.Status())
	log.Printf("agent:response body from server:%v", string(resp.Body()))
	log.Printf("Content-Encoding:%v", resp.Header().Get("Content-Encoding"))


	return err
}

func (a Agent) Update(ctx context.Context, metrics *Metrics) {
	var m runtime.MemStats

	ticker := time.NewTicker(time.Duration(a.Configuration.PollInterval))

	defer ticker.Stop()
repeatAgain:
	select {
	case <-ticker.C:
		{
			runtime.ReadMemStats(&m)

			a.UpdateLocker.Lock()

			metrics.Alloc = Gauge(m.Alloc)
			metrics.BuckHashSys = Gauge(m.BuckHashSys)
			metrics.Frees = Gauge(m.Frees)
			metrics.GCCPUFraction = Gauge(m.GCCPUFraction)
			metrics.GCSys = Gauge(m.GCSys)
			metrics.HeapAlloc = Gauge(m.HeapAlloc)
			metrics.HeapIdle = Gauge(m.HeapIdle)
			metrics.HeapInuse = Gauge(m.HeapInuse)
			metrics.HeapObjects = Gauge(m.HeapObjects)
			metrics.HeapReleased = Gauge(m.HeapReleased)
			metrics.HeapSys = Gauge(m.HeapSys)
			metrics.LastGC = Gauge(m.LastGC)
			metrics.Lookups = Gauge(m.Lookups)
			metrics.MCacheInuse = Gauge(m.MCacheInuse)
			metrics.MCacheSys = Gauge(m.MCacheSys)
			metrics.MSpanInuse = Gauge(m.MSpanInuse)
			metrics.MSpanSys = Gauge(m.MSpanSys)
			metrics.Mallocs = Gauge(m.Mallocs)
			metrics.NextGC = Gauge(m.NextGC)
			metrics.NumForcedGC = Gauge(m.NumForcedGC)
			metrics.NumGC = Gauge(m.NumGC)
			metrics.OtherSys = Gauge(m.OtherSys)
			metrics.PauseTotalNs = Gauge(m.PauseTotalNs)
			metrics.StackInuse = Gauge(m.StackInuse)
			metrics.StackSys = Gauge(m.StackSys)
			metrics.Sys = Gauge(m.Sys)
			metrics.TotalAlloc = Gauge(m.TotalAlloc)
			metrics.RandomValue = Gauge(rand.Int63())
			metrics.PollCount++


			a.UpdateLocker.Unlock()
			goto repeatAgain
		}
	case <-ctx.Done():
		{
			log.Println("Metrics reading cancelled by context")
			return
		}
	}

}

func (a Agent) UpdateGOPS(ctx context.Context, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration(a.Configuration.PollInterval))
	defer ticker.Stop()
repeatAgain:
	select {
	case <-ticker.C:
		{
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Println(err)
				ctx.Done()
			}
			i, err := cpu.Times(true)
			if err != nil {
				log.Println(err)
				ctx.Done()
			}
			a.UpdateLocker.Lock()
			metrics.TotalMemory = Gauge(v.Total)
			metrics.FreeMemory = Gauge(v.Free)
			metrics.CPUutilization1 = Gauge(i[0].System)
			a.UpdateLocker.Unlock()

			goto repeatAgain
		}
	case <-ctx.Done():
		{
			log.Println("context is done received, metrics reading cancelled by app context")
			return
		}
	}

}

func (a Agent) CompressData(data map[*sendData]bool) map[*sendData]bool {

	switch a.Configuration.CompressType {
	// case "deflate":
	// 	{
	// 		for k := range data {

	// 			var b bytes.Buffer

	// 			w, err := flate.NewWriter(&b, flate.BestCompression)
	// 			if err != nil {
	// 				log.Fatalf("failed init compress writer: %v", err)
	// 			}
	// 			_, err = w.Write(*k.JSONBody)
	// 			if err != nil {
	// 				log.Fatalf("failed write data to compress temporary buffer: %v", err)
	// 			}

	// 			err = w.Close()
	// 			if err != nil {
	// 				log.Fatalf("failed compress data: %v", err)
	// 			}
	// 			body := b.Bytes()
	// 			k.JSONBody = &body
	// 		}
	// 	}
	case "gzip":
		{
			var body any
			for k := range data {
				if k.JSONBody != nil {
					body = *k.JSONBody
				} else if k.JSONBatchBody != nil {
					body = *k.JSONBatchBody
				} else {
					logFatal(errors.New("agent:nothing to marshal as sendData bodies are nil"))
				}

				b, err := json.Marshal(body)
				logFatal(err)

				k.compressedBody, err = compression.GzipCompress(b)
				logFatal(err)
			}
		}
	}

	return data
}
func (a Agent) prepareData(metrics *Metrics) map[*sendData]bool {
	m := make(map[*sendData]bool)
	keys := make(HeaderKeys)

	switch a.Configuration.CompressType {
	case "deflate":
		{
			keys["Accept-Encoding"] = "deflate"
			keys["Content-Encoding"] = "deflate"
		}
	case "gzip":
		{
			keys["Accept-Encoding"] = "gzip"
			keys["Content-Encoding"] = "gzip"
		}
	}

	switch a.Configuration.UseJSON {

	case 2: //JSON Batch
		{
			keys["Content-Type"] = "application/json"
			keys["Accept"] = "application/json"

			data := &sendData{
				url:    a.baseURL.JoinPath("updates"),
				keys:   keys,
				signer: a.Signer,
			}

			AddGaugeDataJSONToBatch(data, metrics.Alloc, "Alloc")
			AddGaugeDataJSONToBatch(data, metrics.GCCPUFraction, "GCCPUFraction")
			AddGaugeDataJSONToBatch(data, metrics.GCSys, "GCSys")
			AddGaugeDataJSONToBatch(data, metrics.HeapAlloc, "HeapAlloc")
			AddGaugeDataJSONToBatch(data, metrics.HeapIdle, "HeapIdle")
			AddGaugeDataJSONToBatch(data, metrics.HeapInuse, "HeapInuse")
			AddGaugeDataJSONToBatch(data, metrics.HeapObjects, "HeapObjects")
			AddGaugeDataJSONToBatch(data, metrics.HeapReleased, "HeapReleased")
			AddGaugeDataJSONToBatch(data, metrics.HeapSys, "HeapSys")
			AddGaugeDataJSONToBatch(data, metrics.LastGC, "LastGC")
			AddGaugeDataJSONToBatch(data, metrics.Lookups, "Lookups")
			AddGaugeDataJSONToBatch(data, metrics.MCacheSys, "MCacheSys")
			AddGaugeDataJSONToBatch(data, metrics.MSpanInuse, "MSpanInuse")
			AddGaugeDataJSONToBatch(data, metrics.MSpanSys, "MSpanSys")
			AddGaugeDataJSONToBatch(data, metrics.Mallocs, "Mallocs")
			AddGaugeDataJSONToBatch(data, metrics.NextGC, "NextGC")
			AddGaugeDataJSONToBatch(data, metrics.NumForcedGC, "NumForcedGC")
			AddGaugeDataJSONToBatch(data, metrics.NumGC, "NumGC")
			AddGaugeDataJSONToBatch(data, metrics.OtherSys, "OtherSys")
			AddGaugeDataJSONToBatch(data, metrics.PauseTotalNs, "PauseTotalNs")
			AddGaugeDataJSONToBatch(data, metrics.StackInuse, "StackInuse")
			AddGaugeDataJSONToBatch(data, metrics.StackSys, "StackSys")
			AddGaugeDataJSONToBatch(data, metrics.Sys, "Sys")
			AddGaugeDataJSONToBatch(data, metrics.TotalAlloc, "TotalAlloc")
			AddGaugeDataJSONToBatch(data, metrics.RandomValue, "RandomValue")
			AddGaugeDataJSONToBatch(data, metrics.Frees, "Frees")
			AddCounterDataJSONToBatch(data, metrics.PollCount, "PollCount")

			AddGaugeDataJSONToBatch(data, metrics.TotalMemory, "TotalMemory")
			AddGaugeDataJSONToBatch(data, metrics.FreeMemory, "FreeMemory")
			AddGaugeDataJSONToBatch(data, metrics.CPUutilization1, "CPUutilization1")

			m[data] = true
		}
	case 1: //JSON
		{
			keys["Content-Type"] = "application/json"
			keys["Accept"] = "application/json"

			data := sendData{
				url:    a.baseURL.JoinPath("update"),
				keys:   keys,
				signer: a.Signer,
			}
			AddGaugeDataJSON(data, metrics.Alloc, "Alloc", m)
			AddGaugeDataJSON(data, metrics.Frees, "Frees", m)
			AddGaugeDataJSON(data, metrics.GCCPUFraction, "GCCPUFraction", m)
			AddGaugeDataJSON(data, metrics.GCSys, "GCSys", m)
			AddGaugeDataJSON(data, metrics.HeapAlloc, "HeapAlloc", m)
			AddGaugeDataJSON(data, metrics.HeapIdle, "HeapIdle", m)
			AddGaugeDataJSON(data, metrics.HeapInuse, "HeapInuse", m)
			AddGaugeDataJSON(data, metrics.HeapObjects, "HeapObjects", m)
			AddGaugeDataJSON(data, metrics.HeapReleased, "HeapReleased", m)
			AddGaugeDataJSON(data, metrics.HeapSys, "HeapSys", m)
			AddGaugeDataJSON(data, metrics.LastGC, "LastGC", m)
			AddGaugeDataJSON(data, metrics.Lookups, "Lookups", m)
			AddGaugeDataJSON(data, metrics.MCacheSys, "MCacheSys", m)
			AddGaugeDataJSON(data, metrics.MSpanInuse, "MSpanInuse", m)
			AddGaugeDataJSON(data, metrics.MSpanSys, "MSpanSys", m)
			AddGaugeDataJSON(data, metrics.Mallocs, "Mallocs", m)
			AddGaugeDataJSON(data, metrics.NextGC, "NextGC", m)
			AddGaugeDataJSON(data, metrics.NumForcedGC, "NumForcedGC", m)
			AddGaugeDataJSON(data, metrics.NumGC, "NumGC", m)
			AddGaugeDataJSON(data, metrics.OtherSys, "OtherSys", m)
			AddGaugeDataJSON(data, metrics.PauseTotalNs, "PauseTotalNs", m)
			AddGaugeDataJSON(data, metrics.StackInuse, "StackInuse", m)
			AddGaugeDataJSON(data, metrics.StackSys, "StackSys", m)
			AddGaugeDataJSON(data, metrics.Sys, "Sys", m)
			AddGaugeDataJSON(data, metrics.TotalAlloc, "TotalAlloc", m)
			AddGaugeDataJSON(data, metrics.RandomValue, "RandomValue", m)
			AddCounterDataJSON(data, metrics.PollCount, "PollCount", m)

			//GOPSUtil
			AddGaugeDataJSON(data, metrics.TotalMemory, "TotalMemory", m)
			AddGaugeDataJSON(data, metrics.FreeMemory, "FreeMemory", m)
			AddGaugeDataJSON(data, metrics.CPUutilization1, "CPUutilization1", m)

			//// value1, value2 := int64(rand.Int31()), int64(rand.Int31())
			////// var value0 int64
			////// value1, value2 := int64(1), int64(2)
			////// // //check api no value POST with expected response
			//baseURL := url.URL{Scheme: a.Configuration.Scheme, Host: a.Configuration.Address}
			//dataAPI := sendData{
			//	url:    baseURL.JoinPath("value"),
			//	keys:   keys,
			//	signer: a.Signer,
			//}
			////AddCounterDataJSON(dataAPI, Counter(value1), "SetGet12344", m)
			////AddCounterDataJSON(dataAPI, Counter(value2), "SetGet12344", m)
			//AddCounterDataJSON(dataAPI, -1, "PollCount", m)
			////log.Printf("sum:%v", value1+value2+value0)
		}
	case 0:
		{

			keys["Content-Type"] = "plain/text"
			keys["Accept"] = "text/html"

			data := sendData{
				url:  a.baseURL.JoinPath("update"),
				keys: keys,
			}

			AddGaugeData(data, metrics.Alloc, "Alloc", m)
			AddGaugeData(data, metrics.Frees, "Frees", m)
			AddGaugeData(data, metrics.GCCPUFraction, "GCCPUFraction", m)
			AddGaugeData(data, metrics.GCSys, "GCSys", m)
			AddGaugeData(data, metrics.HeapAlloc, "HeapAlloc", m)
			AddGaugeData(data, metrics.HeapIdle, "HeapIdle", m)
			AddGaugeData(data, metrics.HeapInuse, "HeapInuse", m)
			AddGaugeData(data, metrics.HeapObjects, "HeapObjects", m)
			AddGaugeData(data, metrics.HeapReleased, "HeapReleased", m)
			AddGaugeData(data, metrics.HeapSys, "HeapSys", m)
			AddGaugeData(data, metrics.LastGC, "LastGC", m)
			AddGaugeData(data, metrics.Lookups, "Lookups", m)
			AddGaugeData(data, metrics.MCacheSys, "MCacheSys", m)
			AddGaugeData(data, metrics.MSpanInuse, "MSpanInuse", m)
			AddGaugeData(data, metrics.MSpanSys, "MSpanSys", m)
			AddGaugeData(data, metrics.Mallocs, "Mallocs", m)
			AddGaugeData(data, metrics.NextGC, "NextGC", m)
			AddGaugeData(data, metrics.NumForcedGC, "NumForcedGC", m)
			AddGaugeData(data, metrics.NumGC, "NumGC", m)
			AddGaugeData(data, metrics.OtherSys, "OtherSys", m)
			AddGaugeData(data, metrics.PauseTotalNs, "PauseTotalNs", m)
			AddGaugeData(data, metrics.StackInuse, "StackInuse", m)
			AddGaugeData(data, metrics.StackSys, "StackSys", m)
			AddGaugeData(data, metrics.Sys, "Sys", m)
			AddGaugeData(data, metrics.TotalAlloc, "TotalAlloc", m)
			AddGaugeData(data, metrics.RandomValue, "RandomValue", m)
			AddCounterData(data, metrics.PollCount, "PollCount", m)
			//GOPSUtil
			AddGaugeData(data, metrics.TotalMemory, "TotalMemory", m)
			AddGaugeData(data, metrics.FreeMemory, "FreeMemory", m)
			AddGaugeData(data, metrics.CPUutilization1, "CPUutilization1", m)
		}

	}

	return m
}

func (a Agent) Send(ctx context.Context, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration(a.Configuration.ReportInterval))

	defer ticker.Stop()

	//Worker pool
	var f workerpool.TypicalJobFunction[*sendData]
	var wp workerpool.WorkerPool[*sendData]

	if a.Configuration.RateLimit > 0 {
		f = func(key *sendData) workerpool.JobResult {
			err := key.SendData(a.Client)
			if err != nil {
				log.Println(err)
				return workerpool.JobResult{Result: err.Error()}
			}

			return workerpool.JobResult{Result: "OK"}
		}
		wp = workerpool.NewWorkerPool[*sendData](a.Configuration.RateLimit)
		// worker pool function
		wp.Start(ctx)
	}
repeatAgain:
	select {
	case <-ticker.C:
		{
			dataPackage := a.prepareData(metrics)
			dataPackage = a.CompressData(dataPackage)
			i := 0
			for key := range dataPackage {
				i++
				name := fmt.Sprintf("Send metric job %v", i)
				switch a.Configuration.RateLimit {
				case 0,1:
					{
						err := key.SendData(a.Client)
						if err != nil {
							log.Println(err)
							return
						}
					}
				default:
					{
						job:=workerpool.Job[*sendData]{Name: name, Data: key, Func: f}
						wp.SendJob(ctx,job)
					}
				}
			}
			goto repeatAgain
		}
	case <-ctx.Done():
		{
			break
		}
	}

}

func (a Agent) Run(ctx context.Context) {


	metrics := Metrics{}
	go a.Update(ctx, &metrics)
	go a.UpdateGOPS(ctx, &metrics)
	go a.Send(ctx, &metrics)

}
