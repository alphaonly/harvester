package configuration

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alphaonly/harvester/internal/schema"
)

const ServerDefaultJSON = `{"ADDRESS":"localhost:8080","STORE_INTERVAL": "300s","STORE_FILE":"/tmp/devops-metrics-db.json","RESTORE":true}`
const AgentDefaultJSON = `{"POLL_INTERVAL":"2s","REPORT_INTERVAL":"10s","ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":true}`

type AgentConfiguration struct {
	Address        string          `json:"ADDRESS,omitempty"`
	Scheme         string          `json:"SCHEME,omitempty"`
	CompressType   string          `json:"COMPRESS_TYPE,omitempty"`
	PollInterval   schema.Duration `json:"POLL_INTERVAL,omitempty"`
	ReportInterval schema.Duration `json:"REPORT_INTERVAL,omitempty"`
	UseJSON        bool            `json:"USE_JSON,omitempty"`
}

type ServerConfiguration struct {
	Address       string          `json:"ADDRESS,omitempty"`
	StoreInterval schema.Duration `json:"STORE_INTERVAL,omitempty"`
	StoreFile     string          `json:"STORE_FILE,omitempty"`
	Restore       bool            `json:"RESTORE,omitempty"`
	Port          string          `json:"PORT,omitempty"` //additionally for listen and serve func
}

type FileArchiveConfiguration struct {
	StoreFile string
}

// func getInterval(string) time.Duration

func UnMarshalServerDefaults(s string) ServerConfiguration {
	sc := ServerConfiguration{}
	err := json.Unmarshal([]byte(s), &sc)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return sc

}
func UnMarshalAgentDefaults(s string) AgentConfiguration {
	ac := AgentConfiguration{}
	err := json.Unmarshal([]byte(s), &ac)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return ac
}

func NewAgentConfiguration() *AgentConfiguration {
	c := UnMarshalAgentDefaults(AgentDefaultJSON)
	// ac.PORT = ":" + strings.Split(ac.ADDRESS, ":")[1]

	return &c

}
func NewServerConfiguration() *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.Port = ":" + strings.Split(c.Address, ":")[1]
	return &c

}

func (c *AgentConfiguration) UpdateFromEnvironment() {
	log.Printf("Have environmet: %v", os.Environ())
	c.Address = getEnv("ADDRESS", c.Address).(string)
	c.CompressType = getEnv("COMPRESS_TYPE", c.CompressType).(string)
	c.PollInterval = getEnv("POLL_INTERVAL", c.PollInterval).(schema.Duration)
	c.ReportInterval = getEnv("REPORT_INTERVAL", c.ReportInterval).(schema.Duration)
	c.Scheme = getEnv("SCHEME", c.Scheme).(string)
	c.UseJSON = getEnv("USE_JSON", c.UseJSON).(bool)
}

func (c *ServerConfiguration) UpdateFromEnvironment() {
	log.Printf("Have environmet: %v", os.Environ())
	c.Address = getEnv("ADDRESS", c.Address).(string)
	c.Restore = getEnv("RESTORE", c.Restore).(bool)
	c.StoreFile = getEnv("STORE_FILE", c.StoreFile).(string)
	c.StoreInterval = getEnv("STORE_INTERVAL", c.StoreInterval).(schema.Duration)

	//PORT is derived from ADDRESS
	c.Port = ":" + strings.Split(c.Address, ":")[1]
}

func (c *AgentConfiguration) UpdateFromFlags() {
	dc := NewAgentConfiguration()
	var (
		a = flag.String("a", dc.Address, "Domain name and :port")
		p = flag.Duration("p", time.Duration(dc.PollInterval), "Poll interval")
		r = flag.Duration("r", time.Duration(dc.ReportInterval), "Report interval")
		j = flag.Bool("j", dc.UseJSON, "Use JSON true/false")
		t = flag.String("t", dc.CompressType, "Compress type: \"deflate\" supported")
	)

	flag.Parse()

	//Если значение параметра из переменных окружения равно по умолчанию, то обновляем из флагов

	message := "variable %v  updated from flags, value %v"
	if c.Address == dc.Address {
		c.Address = *a
		log.Printf(message, "ADDRESS", c.Address)
	}
	if c.PollInterval == dc.PollInterval {
		c.PollInterval = schema.Duration(*p)
		log.Printf(message, "POLL_INTERVAL", c.PollInterval)
	}

	if c.ReportInterval == dc.ReportInterval {
		c.ReportInterval = schema.Duration(*r)
		log.Printf(message, "REPORT_INTERVAL", c.ReportInterval)
	}

	if c.UseJSON == dc.UseJSON {
		c.UseJSON = *j
		log.Printf(message, "USE_JSON", c.UseJSON)
	}

	if c.CompressType != dc.CompressType {
		c.CompressType = *t
		log.Printf(message, "COMPRESS_TYPE", c.CompressType)
	}

}

func (c *ServerConfiguration) UpdateFromFlags() {

	dc := NewServerConfiguration()

	var (
		a = flag.String("a", dc.Address, "Domain name and :port")
		i = flag.Duration("i", time.Duration(dc.StoreInterval), "Store interval")
		f = flag.String("f", dc.StoreFile, "Store file full path")
		r = flag.Bool("r", dc.Restore, "Restore from external storage:true/false")
	)
	flag.Parse()

	message := "variable %v  updated from flags, value %v"
	//Если значение из переменных равно значению по умолчанию, тогда берем из flags
	if c.Address == dc.Address {
		c.Address = *a
		c.Port = ":" + strings.Split(c.Address, ":")[1]
		log.Printf(message, "ADDRESS", c.Address)
		log.Printf(message, "PORT", c.Port)
	}
	if c.StoreInterval == dc.StoreInterval {
		c.StoreInterval = schema.Duration(*i)
		log.Printf(message, "STORE_INTERVAL", c.StoreInterval)
	}
	if c.StoreFile == dc.StoreFile {
		c.StoreFile = *f
		log.Printf(message, "STORE_FILE", c.StoreFile)
	}
	if c.Restore == dc.Restore {
		c.Restore = *r
		log.Printf(message, "RESTORE", c.Restore)
	}
}

func getEnv(variableName string, variableValue interface{}) (changedValue interface{}) {
	var stringVal string
	var err error

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}
	var exists bool
	stringVal, exists = os.LookupEnv(variableName)
	if !exists {
		log.Printf("variable "+variableName+" not presented in environment, remains default:%v", variableValue)
		return variableValue
	}
	switch variableValue.(type) {
	case string:
		changedValue = stringVal
	case int64:
		{
			changedValue, err = strconv.ParseInt(stringVal, 10, 64)
			if err != nil {
				log.Fatal("Int64 Parse error")
			}
		}
	case bool:
		{
			changedValue, err = strconv.ParseBool(stringVal)
			if err != nil {
				log.Fatal("Bool Parse error")
			}
		}
	case schema.Duration:
		{
			changedValue, err = time.ParseDuration(stringVal)
			if err != nil {
				log.Fatal("Duration Parse error")
			}
		}
	default:
		log.Fatal("unknown type getEnv")
	}
	if stringVal != "" {
		log.Println("variable " + variableName + " presented in environment, value: " + stringVal)
	}

	return changedValue
}
