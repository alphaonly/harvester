package configuration

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

const ServerDefaultJSON = `{"ADDRESS":"localhost:8080","STORE_INTERVAL": 300,"STORE_FILE":"/tmp/devops-metrics-db.json","RESTORE":true}`
const AgentDefaultJSON = `{"POLL_INTERVAL":2,"REPORT_INTERVAL":10,"ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":true}`

type AgentConfiguration struct {
	Address        string `json:"ADDRESS,omitempty"`
	Scheme         string `json:"SCHEME,omitempty"`
	CompressType   string `json:"COMPRESS_TYPE,omitempty"`
	PollInterval   int64  `json:"POLL_INTERVAL,omitempty"`
	ReportInterval int64  `json:"REPORT_INTERVAL,omitempty"`
	UseJSON        bool   `json:"USE_JSON,omitempty"`
}

type ServerConfiguration struct {
	Address       string `json:"ADDRESS,omitempty"`
	StoreInterval int64  `json:"STORE_INTERVAL,omitempty"`
	StoreFile     string `json:"STORE_FILE,omitempty"`
	Restore       bool   `json:"RESTORE,omitempty"`
	Port          string `json:"PORT,omitempty"` //additionally for listen and serve func
}

type FileArchiveConfiguration struct{
	StoreFile     string 
}

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

	c.Address = getEnv("ADDRESS", c.Address).(string)
	c.CompressType = getEnv("COMPRESS_TYPE", c.CompressType).(string)
	c.PollInterval = getEnv("POLL_INTERVAL", c.PollInterval).(int64)
	c.ReportInterval = getEnv("REPORT_INTERVAL", c.ReportInterval).(int64)
	c.Scheme = getEnv("SCHEME", c.Scheme).(string)
	c.UseJSON = getEnv("USE_JSON", c.UseJSON).(bool)
}
func (c *ServerConfiguration) UpdateFromEnvironment() {

	c.Address = getEnv("ADDRESS", c.Address).(string)
	c.Restore = getEnv("RESTORE", c.Restore).(bool)
	c.StoreFile = getEnv("STORE_FILE", c.StoreFile).(string)
	c.StoreInterval = getEnv("STORE_INTERVAL", c.StoreInterval).(int64)

	//PORT is derived from ADDRESS
	c.Port = ":" + strings.Split(c.Address, ":")[1]
}

func (c *AgentConfiguration) UpdateFromFlags() {
	var (
		a = flag.String("a", c.Address, "Domain name and :port")
		p = flag.Int64("p", c.PollInterval, "Poll interval")
		r = flag.Int64("r", c.ReportInterval, "Report interval")
		j = flag.Bool("j", c.UseJSON, "Use JSON true/false")
		t = flag.String("t", c.CompressType, "Compress type: \"deflate\" supported")
	)

	flag.Parse()
	dc := NewAgentConfiguration()

	//Если значение параметра из переменных окружения равно по умолчанию, то обновляем из флагов
	switch true {
	case c.Address == dc.Address:
		c.Address = *a
	case c.PollInterval == dc.PollInterval:
		c.PollInterval = *p
	case c.ReportInterval == dc.ReportInterval:
		c.ReportInterval = *r
	case c.UseJSON == dc.UseJSON:
		c.UseJSON = *j
	case c.CompressType == dc.CompressType:
		c.CompressType = *t
	}

}

func (c *ServerConfiguration) UpdateFromFlags() {
	var (
		a *string = flag.String("a", c.Address, "Domain name and :port")
		i *int64  = flag.Int64("i", c.StoreInterval, "Store interval")
		f *string = flag.String("f", c.StoreFile, "Store file full path")
		r *bool   = flag.Bool("r", c.Restore, "Restore from external storage:true/false")
	)
	flag.Parse()
	dc := NewServerConfiguration()

	switch true {
	case c.Address == dc.Address:
		c.Address = *a
		c.Port = ":" + strings.Split(c.Address, ":")[1]
	case c.StoreInterval == dc.StoreInterval:
		c.StoreInterval = *i
	case c.StoreFile == dc.StoreFile:
		c.StoreFile = *f
	case c.Restore == dc.Restore:
		c.Restore = *r
	}

	// flag.Parse()
	// if a != nil {
	// 	c.Address = *a
	// 	c.Port = ":" + strings.Split(c.Address, ":")[1]
	// }
	// if i != nil {
	// 	c.StoreInterval = *i
	// }
	// if f != nil {
	// 	c.StoreFile = *f
	// }
	// if r != nil {
	// 	c.Restore = *r
	// }
}

func getEnv(variableName string, variableValue interface{}) (changedValue interface{}) {
	var stringVal string
	var err error

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}

	stringVal = os.Getenv(variableName)
	if stringVal == "" {
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
	default:
		log.Fatal("unknown type getEnv")
	}
	if stringVal != "" {
		log.Println("variable " + variableName + "presented in environment, value: " + stringVal)
	}

	return changedValue
}
