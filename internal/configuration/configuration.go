package configuration

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alphaonly/harvester/internal/common/logging"

	"github.com/alphaonly/harvester/internal/schema"
)

const ServerDefaultJSON = `{"ADDRESS":"localhost:8080","STORE_INTERVAL": "300s","STORE_FILE":"/tmp/devops-metrics-db.json","RESTORE":true,"KEY":"","CRYPTO_KEY":"","ENABLE_HTTPS":false}`
const AgentDefaultJSON = `{"POLL_INTERVAL":"2s","REPORT_INTERVAL":"10s","ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":1,"KEY":"","RATE_LIMIT":1,"CRYPTO_KEY":""}`

type AgentConfiguration struct {
	Address        string          `json:"ADDRESS,omitempty"`
	Scheme         string          `json:"SCHEME,omitempty"`
	CompressType   string          `json:"COMPRESS_TYPE,omitempty"`
	PollInterval   schema.Duration `json:"POLL_INTERVAL,omitempty"`
	ReportInterval schema.Duration `json:"REPORT_INTERVAL,omitempty"`
	UseJSON        int             `json:"USE_JSON,omitempty"`
	Key            string          `json:"KEY,omitempty"`
	RateLimit      int             `json:"RATE_LIMIT,omitempty"`
	CryptoKey      string          `json:"CRYPTO_KEY,omitempty"` //path to public key file
	EnvChanged     map[string]bool
}

type ServerConfiguration struct {
	Address       string          `json:"ADDRESS,omitempty"`
	StoreInterval schema.Duration `json:"STORE_INTERVAL,omitempty"`
	StoreFile     string          `json:"STORE_FILE,omitempty"`
	Restore       bool            `json:"RESTORE,omitempty"`
	Port          string          `json:"PORT,omitempty"` //additionally for listen and serve func
	Key           string          `json:"KEY,omitempty"`
	DatabaseDsn   string          `json:"DATABASE_DSN,omitempty"`
	CryptoKey     string          `json:"CRYPTO_KEY,omitempty"`     //path to private key file
	Config        string          `json:"CONFIG,omitempty"`         //path to config file
	TrustedSubnet string          `json:"TRUSTED_SUBNET,omitempty"` //trusted subnet declaration
	EnvChanged    map[string]bool
}

type AgentConfigurationOption func(*AgentConfiguration)
type ServerConfigurationOption func(*ServerConfiguration)

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

func NewAgentConf(options ...AgentConfigurationOption) *AgentConfiguration {
	c := UnMarshalAgentDefaults(AgentDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}
func NewServerConf(options ...ServerConfigurationOption) *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}

func UpdateACFromEnvironment(c *AgentConfiguration) {

	c.Address = getEnv("ADDRESS", &StrValue{c.Address}, c.EnvChanged).(string)
	c.CompressType = getEnv("COMPRESS_TYPE", &StrValue{c.CompressType}, c.EnvChanged).(string)
	c.PollInterval = getEnv("POLL_INTERVAL", &DurValue{c.PollInterval}, c.EnvChanged).(schema.Duration)
	c.ReportInterval = getEnv("REPORT_INTERVAL", &DurValue{c.ReportInterval}, c.EnvChanged).(schema.Duration)
	c.Scheme = getEnv("SCHEME", &StrValue{c.Scheme}, c.EnvChanged).(string)
	c.UseJSON = getEnv("USE_JSON", &IntValue{c.UseJSON}, c.EnvChanged).(int)
	c.Key = getEnv("KEY", &StrValue{c.Key}, c.EnvChanged).(string)
	c.RateLimit = getEnv("RATE_LIMIT", &IntValue{c.RateLimit}, c.EnvChanged).(int)
	c.CryptoKey = getEnv("CRYPTO_KEY", &StrValue{c.CryptoKey}, c.EnvChanged).(string)
}

func UpdateSCFromEnvironment(c *ServerConfiguration) {
	c.Address = getEnv("ADDRESS", &StrValue{c.Address}, c.EnvChanged).(string)
	c.Restore = getEnv("RESTORE", &BoolValue{c.Restore}, c.EnvChanged).(bool)
	c.StoreFile = getEnv("STORE_FILE", &StrValue{c.StoreFile}, c.EnvChanged).(string)
	c.StoreInterval = getEnv("STORE_INTERVAL", &DurValue{c.StoreInterval}, c.EnvChanged).(schema.Duration)
	c.Key = getEnv("KEY", &StrValue{c.Key}, c.EnvChanged).(string)
	//PORT is derived from ADDRESS
	c.Port = ":" + strings.Split(c.Address, ":")[1]
	c.DatabaseDsn = getEnv("DATABASE_DSN", &StrValue{c.DatabaseDsn}, c.EnvChanged).(string)
	c.CryptoKey = getEnv("CRYPTO_KEY", &StrValue{c.CryptoKey}, c.EnvChanged).(string)
	c.Config = getEnv("CONFIG", &StrValue{c.Config}, c.EnvChanged).(string)
	c.TrustedSubnet = getEnv("TRUSTED_SUBNET", &StrValue{c.TrustedSubnet}, c.EnvChanged).(string)
}

func UpdateACFromFlags(c *AgentConfiguration) {
	dc := NewAgentConfiguration()
	var (
		a  = flag.String("a", dc.Address, "Domain name and :port")
		p  = flag.Duration("p", time.Duration(dc.PollInterval), "Poll interval")
		r  = flag.Duration("r", time.Duration(dc.ReportInterval), "Report interval")
		j  = flag.Int("j", dc.UseJSON, "Use JSON 0-No JSON,1- JSON, 2-JSON Batch")
		t  = flag.String("t", dc.CompressType, "Compress type: \"deflate\" supported")
		k  = flag.String("k", dc.Key, "string key for hash signing")
		l  = flag.Int("l", dc.RateLimit, "Number of parallel inbound requests ")
		cr = flag.String("crypto-key", dc.CryptoKey, "string contains a full path to a public key file ")
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
	if !c.EnvChanged["RATE_LIMIT"] {
		c.RateLimit = *l
		log.Printf(message, "RATE_LIMIT", c.RateLimit)
	}
	if !c.EnvChanged["CRYPTO_KEY"] {
		c.CryptoKey = *cr
		log.Printf(message, "CRYPTO_KEY", c.CryptoKey)
	}
}

func UpdateSCFromFlags(c *ServerConfiguration) {

	dc := NewServerConfiguration()

	var (
		a   = flag.String("a", dc.Address, "Domain name and :port")
		i   = flag.Duration("i", time.Duration(dc.StoreInterval), "Store interval")
		f   = flag.String("f", dc.StoreFile, "Store file full path")
		r   = flag.Bool("r", dc.Restore, "Restore from external storage:true/false")
		k   = flag.String("k", dc.Key, "string key for hash signing")
		d   = flag.String("d", dc.DatabaseDsn, "database destination string")
		cr  = flag.String("crypto-key", dc.CryptoKey, "string contains a full path to a private key file ")
		cf1 = flag.String("c", dc.Config, "string contains a full path to configuration JSON File")
		cf2 = flag.String("config", dc.Config, "string contains a full path to configuration JSON File")
		t   = flag.String("t", dc.Config, "string contains a full path to configuration JSON File")
	)
	flag.Parse()

	message := "variable %v  updated from flags, value %v"
	//Если значение из переменных равно значению по умолчанию, тогда берем из flagS
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
	if !c.EnvChanged["DATABASE_DSN"] {
		c.DatabaseDsn = *d
		log.Printf(message, "DATABASE_DSN", c.DatabaseDsn)
	}
	if !c.EnvChanged["CRYPTO_KEY"] {
		c.CryptoKey = *cr
		log.Printf(message, "CRYPTO_KEY", c.CryptoKey)
	}
	if !c.EnvChanged["CONFIG"] {
		switch {
		case *cf1 != "":
			c.Config = *cf1
		case *cf2 != "":
			c.Config = *cf2
		}
		log.Printf(message, "CONFIG", c.Config)
	}
	if !c.EnvChanged["TRUSTED_SUBNET"] {
		c.TrustedSubnet = *t
		log.Printf(message, "TRUSTED_SUBNET", c.TrustedSubnet)
	}
}

// UpdateSCFromConfigFile - read server configuration from JSON file
func UpdateSCFromConfigFile(c *ServerConfiguration) {
	fc := NewServerConfiguration()
	//check if there is a previously given path to config JSON JSONConfigFile
	if c.Config == "" {
		return
	}
	//read JSONConfigFile
	JSONConfigFile, err := os.Open(c.Config)
	logging.LogFatal(err)

	buf := bufio.NewReader(JSONConfigFile)
	logging.LogFatal(err)
	b := make([]byte, 4096)
	readBytes, err := buf.Read(b)

	//Unmarshal JSON bytes from file
	err = json.Unmarshal(b[:readBytes], fc)
	logging.LogFatal(err)
	//Analyze which parameters are present and changed
	message := "variable %v  updated from file configuration, value %v"

	if !c.EnvChanged["ADDRESS"] {
		c.Address = fc.Address
		c.Port = ":" + strings.Split(c.Address, ":")[1]
		log.Printf(message, "ADDRESS", c.Address)
		log.Printf(message, "PORT", c.Port)
	}
	if !c.EnvChanged["STORE_INTERVAL"] {
		c.StoreInterval = fc.StoreInterval
		log.Printf(message, "STORE_INTERVAL", c.StoreInterval)
	}
	if !c.EnvChanged["STORE_FILE"] {
		c.StoreFile = fc.StoreFile
		log.Printf(message, "STORE_FILE", c.StoreFile)
	}
	if !c.EnvChanged["RESTORE"] {
		c.Restore = fc.Restore
		log.Printf(message, "RESTORE", c.Restore)
	}
	if !c.EnvChanged["KEY"] {
		c.Key = fc.Key
		log.Printf(message, "KEY", c.Key)
	}
	if !c.EnvChanged["DATABASE_DSN"] {
		c.DatabaseDsn = fc.DatabaseDsn
		log.Printf(message, "DATABASE_DSN", c.DatabaseDsn)
	}
	if !c.EnvChanged["CRYPTO_KEY"] {
		c.CryptoKey = fc.CryptoKey
		log.Printf(message, "CRYPTO_KEY", c.CryptoKey)
	}
	if !c.EnvChanged["TRUSTED_SUBNET"] {
		c.TrustedSubnet = fc.TrustedSubnet
		log.Printf(message, "TRUSTED_SUBNET", c.TrustedSubnet)
	}
}

type VariableValue interface {
	Get() interface{}
	Set(string)
}
type StrValue struct {
	value string
}

func (v *StrValue) Get() interface{} {
	return v.value
}

func (v *StrValue) Set(s string) {
	v.value = s
}

type IntValue struct {
	value int
}

func (v IntValue) Get() interface{} {
	return v.value
}
func (v *IntValue) Set(s string) {
	var err error
	v.value, err = strconv.Atoi(s)
	if err != nil {
		log.Fatal("Int Parse error")
	}
}

type BoolValue struct {
	value bool
}

func (v BoolValue) Get() interface{} {
	return v.value
}
func (v *BoolValue) Set(s string) {
	var err error
	v.value, err = strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Bool Parse error")
	}
}

type DurValue struct {
	value schema.Duration
}

func (v DurValue) Get() interface{} {
	return v.value
}
func (v *DurValue) Set(s string) {
	var err error
	interval, err := time.ParseDuration(s)
	if err != nil {
		log.Fatal("Duration Parse error")
	}
	v.value = schema.Duration(interval)
}

func getEnv(variableName string, variableValue VariableValue, changed map[string]bool) (changedValue interface{}) {
	var stringVal string

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}
	var exists bool
	stringVal, exists = os.LookupEnv(variableName)
	if !exists {
		log.Printf("variable "+variableName+" not presented in environment, remains default:%v", variableValue.Get())
		changed[variableName] = false
		return variableValue.Get()
	}
	variableValue.Set(stringVal)
	changed[variableName] = true
	log.Println("variable " + variableName + " presented in environment, value: " + stringVal)

	return variableValue.Get()
}
