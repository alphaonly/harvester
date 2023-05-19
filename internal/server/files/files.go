package files

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
	"github.com/alphaonly/harvester/internal/server/metricvalue"
)

type Producer interface {
	Write(*map[string]metricvalue.MetricValue) error
	Close() error
}

type Consumer interface {
	Read() (*map[string]metricvalue.MetricValue, error)
	Close() error
}

type metricsProducer struct {
	file *os.File
	buf  *bufio.Writer
}

func NewProducer(filename string) (*metricsProducer, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewWriter(file)

	return &metricsProducer{
		file: file,
		buf:  buf,
	}, nil
}

func (mp *metricsProducer) Write(metricsData *map[string]metricvalue.MetricValue) error {

	// mj:= []metricsjson.MetricsJSON{}

	mapJSON, err := json.Marshal(metricsData)
	if err != nil {
		log.Println("error:cannot marshal data")
		return err
	}

	_, err = mp.buf.Write(mapJSON)
	if err != nil {
		log.Println("error:cannot write file")
		return err
	}
	err = mp.buf.Flush()
	if err != nil {
		log.Println("error:cannot write(flush) file")
		return err
	}
	return nil
}
func (mp *metricsProducer) Close() error {

	return mp.file.Close()
}

type metricsConsumer struct {
	file *os.File
	buf  *bufio.Reader
}

func NewConsumer(filename string) (*metricsConsumer, error) {

	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_EXCL, 0777)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(file)

	return &metricsConsumer{
		file: file,
		buf:  buf,
	}, nil
}

func (mp *metricsConsumer) Read() (metricsData *map[string]metricvalue.MetricValue, err error) {

	mapJSON := make([]byte, 747)
	_, er := mp.buf.Read(mapJSON)
	if er != nil {
		log.Printf("warning:cannot read file %v", mp.file.Name())
		return nil, er
	}

	// if metricsData == nil {
	// 	m := make(map[string]interface{})
	// 	// metricsData = &m
	// }
	m := make(metricsjson.MetricsMapType)
	err = json.Unmarshal(mapJSON, &m)
	if err != nil {
		log.Println("error:cannot unmarshal data:" + err.Error())

		return nil, err
	}

	return metricsData, nil
}
func (mp *metricsConsumer) Close() error {

	return mp.file.Close()
}

// var producer Producer = &metricsProducer{}
// var consumer Consumer = &metricsConsumer{}
