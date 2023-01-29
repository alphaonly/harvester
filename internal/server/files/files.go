package files

import (
	"bufio"
	"encoding/json"
	"fmt"
	M "github.com/alphaonly/harvester/internal/server/metricvalue"
	"log"
	"os"

	metricsjson "github.com/alphaonly/harvester/internal/server/metricsJSON"
)

type Producer interface {
	Write(*map[string]M.MetricValue) error
	Close() error
}

type Consumer interface {
	Read() (*map[string]M.MetricValue, error)
	Close() error
}

type MetricsProducer struct {
	file *os.File
	buf  *bufio.Writer
}

func NewProducer(filename string) (*MetricsProducer, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewWriter(file)

	return &MetricsProducer{
		file: file,
		buf:  buf,
	}, nil
}

func (mp *MetricsProducer) Write(metricsData *map[string]M.MetricValue) error {

	mapJSONBuf, err := json.Marshal(metricsData)
	if err != nil {
		log.Println("error:cannot marshal data")
		return err
	}

	_, err = mp.buf.Write(mapJSONBuf)
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
func (mp *MetricsProducer) Close() error {

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

func (mp *metricsConsumer) Read() (metricsData *map[string]M.MetricValue, err error) {

	scanner := bufio.NewScanner(mp.buf)
	var mapJSON []byte
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		mapJSON = scanner.Bytes()
	}

	*metricsData = make(metricsjson.MetricsMapType)
	err = json.Unmarshal(mapJSON, metricsData)
	if err != nil {
		log.Println("error:cannot unmarshal data:" + err.Error())

		return nil, err
	}

	return metricsData, nil
}
func (mp *metricsConsumer) Close() error {

	return mp.file.Close()
}

var producer Producer = &MetricsProducer{}
var consumer Consumer = &metricsConsumer{}
