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

const ServerDefaultJSON = `{"ADDRESS":"localhost:8080","STORE_INTERVAL": "300s","STORE_FILE":"/tmp/devops-metrics-db.json","RESTORE":true,"KEY":""}`
const AgentDefaultJSON = `{"POLL_INTERVAL":"2s","REPORT_INTERVAL":"10s","ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":true,"KEY":""}`

type AgentConfiguration struct {
	Address        string          `json:"ADDRESS,omitempty"`
	Scheme         string          `json:"SCHEME,omitempty"`
	CompressType   string          `json:"COMPRESS_TYPE,omitempty"`
	PollInterval   schema.Duration `json:"POLL_INTERVAL,omitempty"`
	ReportInterval schema.Duration `json:"REPORT_INTERVAL,omitempty"`
	UseJSON        bool            `json:"USE_JSON,omitempty"`
	Key            string          `json:"KEY,omitempty"`
	EnvChanged     map[string]bool
}

type ServerConfiguration struct {
	Address       string          `json:"ADDRESS,omitempty"`
	StoreInterval schema.Duration `json:"STORE_INTERVAL,omitempty"`
	StoreFile     string          `json:"STORE_FILE,omitempty"`
	Restore       bool            `json:"RESTORE,omitempty"`
	Port          string          `json:"PORT,omitempty"` //additionally for listen and serve func
	Key           string          `json:"KEY,omitempty"`
	EnvChanged    map[string]bool
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
	c.EnvChanged = make(map[string]bool)

	return &c

}
func NewServerConfiguration() *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.Port = ":" + strings.Split(c.Address, ":")[1]
	c.EnvChanged = make(map[string]bool)
	return &c

}

func (c *AgentConfiguration) UpdateFromEnvironment() {
	c.Address = getEnv("ADDRESS", c.Address, c.EnvChanged).(string)
	c.CompressType = getEnv("COMPRESS_TYPE", c.CompressType, c.EnvChanged).(string)
	c.PollInterval = getEnv("POLL_INTERVAL", c.PollInterval, c.EnvChanged).(schema.Duration)
	c.ReportInterval = getEnv("REPORT_INTERVAL", c.ReportInterval, c.EnvChanged).(schema.Duration)
	c.Scheme = getEnv("SCHEME", c.Scheme, c.EnvChanged).(string)
	c.UseJSON = getEnv("USE_JSON", c.UseJSON, c.EnvChanged).(bool)
	c.Key = getEnv("KEY", c.Key, c.EnvChanged).(string)
}

func (c *ServerConfiguration) UpdateFromEnvironment() {
	c.Address = getEnv("ADDRESS", c.Address, c.EnvChanged).(string)
	c.Restore = getEnv("RESTORE", c.Restore, c.EnvChanged).(bool)
	c.StoreFile = getEnv("STORE_FILE", c.StoreFile, c.EnvChanged).(string)
	c.StoreInterval = getEnv("STORE_INTERVAL", c.StoreInterval, c.EnvChanged).(schema.Duration)
	c.Key = getEnv("KEY", c.Key, c.EnvChanged).(string)
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
		k = flag.String("k", dc.Key, "string key for hash signing")
	)

	flag.Parse()

	//Если значение параметра из переменных окружения равно по умолчанию, то обновляем из флагов

	message := "variable %v  updated from flags, value %v"
	if !c.EnvChanged["ADDRESS"] {
		c.Address = *a
		log.Printf(message, "ADDRESS", c.Address)
	}
	if !c.EnvChanged["POLL_INTERVAL"] {
		c.PollInterval = schema.Duration(*p)
		log.Printf(message, "POLL_INTERVAL", c.PollInterval)
	}

	if !c.EnvChanged["REPORT_INTERVAL"] {
		c.ReportInterval = schema.Duration(*r)
		log.Printf(message, "REPORT_INTERVAL", c.ReportInterval)
	}

	if !c.EnvChanged["USE_JSON"] {
		c.UseJSON = *j
		log.Printf(message, "USE_JSON", c.UseJSON)
	}

	if !c.EnvChanged["COMPRESS_TYPE"] {
		c.CompressType = *t
		log.Printf(message, "COMPRESS_TYPE", c.CompressType)
	}
	if !c.EnvChanged["KEY"] {
		c.Key = *k
		log.Printf(message, "KEY", c.Key)
	}

}

func (c *ServerConfiguration) UpdateFromFlags() {

	dc := NewServerConfiguration()

	var (
		a = flag.String("a", dc.Address, "Domain name and :port")
		i = flag.Duration("i", time.Duration(dc.StoreInterval), "Store interval")
		f = flag.String("f", dc.StoreFile, "Store file full path")
		r = flag.Bool("r", dc.Restore, "Restore from external storage:true/false")
		k = flag.String("k", dc.Key, "string key for hash signing")
	)
	flag.Parse()

	message := "variable %v  updated from flags, value %v"
	//Если значение из переменных равно значению по умолчанию, тогда берем из flags
	if !c.EnvChanged["ADDRESS"] {
		c.Address = *a
		c.Port = ":" + strings.Split(c.Address, ":")[1]
		log.Printf(message, "ADDRESS", c.Address)
		log.Printf(message, "PORT", c.Port)
	}
	if !c.EnvChanged["STORE_INTERVAL"] {
		c.StoreInterval = schema.Duration(*i)
		log.Printf(message, "STORE_INTERVAL", c.StoreInterval)
	}
	if !c.EnvChanged["STORE_FILE"] {
		c.StoreFile = *f
		log.Printf(message, "STORE_FILE", c.StoreFile)
	}
	if !c.EnvChanged["RESTORE"] {
		c.Restore = *r
		log.Printf(message, "RESTORE", c.Restore)
	}
	if !c.EnvChanged["KEY"] {
		c.Key = *k
		log.Printf(message, "KEY", c.Key)
	}
}

func getEnv(variableName string, variableValue interface{}, changed map[string]bool) (changedValue interface{}) {
	var stringVal string
	var err error

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}
	var exists bool
	stringVal, exists = os.LookupEnv(variableName)
	if !exists {
		log.Printf("variable "+variableName+" not presented in environment, remains default:%v", variableValue)
		changed[variableName] = false
		return variableValue
	}
	changed[variableName] = true
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
			interval, err := time.ParseDuration(stringVal)
			if err != nil {
				log.Fatal("Duration Parse error")
			}
			changedValue = schema.Duration(interval)
		}
	default:
		log.Fatal("unknown type getEnv")
	}
	if stringVal != "" {
		log.Println("variable " + variableName + " presented in environment, value: " + stringVal)
	}

	return changedValue
}
