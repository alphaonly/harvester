package configuration

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type AgentFlagConfiguration struct {
	variablesMap map[string]string
	Cfg          AgentCfg
}

// Аргументы агента:
// ADDRESS через флаг a=<ЗНАЧЕНИЕ>,
// REPORT_INTERVAL через флаг r=<ЗНАЧЕНИЕ>,
// POLL_INTERVAL через флаг p=<ЗНАЧЕНИЕ>.

func NewAgentFlagConfiguration() *AgentFlagConfiguration {
	m := copyMap(agentDefaults)

	cfg := UnMarshalAgentDefaults(AgentDefaultJSON)
	//default bucket of parameters and their values

	a := flag.String("a", agentDefaults["ADDRESS"], "Domain name and :port")
	p := flag.String("p", agentDefaults["POLL_INTERVAL"], "Poll interval")
	r := flag.String("r", agentDefaults["REPORT_INTERVAL"], "Report interval")
	j := flag.String("j", agentDefaults["USE_JSON"], "Use JSON true/false")
	t := flag.String("t", agentDefaults["COMPRESS_TYPE"], "Compress type: \"deflate\" supported")
	flag.Parse()

	m["ADDRESS"] = *a
	m["POLL_INTERVAL"] = *p
	m["REPORT_INTERVAL"] = *r
	m["USE_JSON"] = *j
	m["COMPRESS_TYPE"] = *t

	cfg.ADDRESS = *a
	cfg.POLL_INTERVAL, _ = strconv.ParseInt(*p, 10, 64)
	cfg.REPORT_INTERVAL, _ = strconv.ParseInt(*r, 10, 64)
	cfg.USE_JSON, _ = strconv.ParseBool(*j)
	cfg.COMPRESS_TYPE = *t

	return &AgentFlagConfiguration{
		variablesMap: m,
		Cfg:          cfg,
	}
}

func (ac *AgentFlagConfiguration) Get(name string) (value string) {
	v := ac.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentFlagConfiguration) GetBool(name string) (value bool) {
	v := ac.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal("bool parsing error value" + v)
	}
	return b
}
func (ac *AgentFlagConfiguration) GetInt(name string) (value int64) {
	v := ac.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
		return
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatal("int parsing error value" + v)
	}
	return i
}

func (ac *AgentFlagConfiguration) read() error {
	log.Fatal("method AgentFlagConfiguration.read is not supported by implementation")
	return nil
}

func (ac *AgentFlagConfiguration) write() {
	// is not supported by implementation
	log.Fatal("method AgentFlagConfiguration.write is not supported by implementation")
}
func (ac *AgentFlagConfiguration) Update() *AgentFlagConfiguration {
	log.Fatal("method AgentFlagConfiguration.Update is not supported by implementation")
	return nil
}

type ServerFlagConfiguration struct {
	variablesMap map[string]string
	Cfg          ServerCfg
}

func NewServerFlagConfiguration() *ServerFlagConfiguration {

	m := copyMap(serverDefaults)

	cfg := UnMarshalServerDefaults(ServerDefaultJSON)

	//default bucket of parameters and their values
	a := flag.String("a", serverDefaults["ADDRESS"], "Address")
	i := flag.String("i", serverDefaults["STORE_INTERVAL"], "Store interval")
	f := flag.String("f", serverDefaults["STORE_FILE"], "Store interval")
	r := flag.String("r", serverDefaults["RESTORE"], "Restore from external storage:true/false")

	flag.Parse()

	m["ADDRESS"] = *a
	m["STORE_INTERVAL"] = *i
	m["STORE_FILE"] = *f
	m["RESTORE"] = *r

	cfg.ADDRESS = *a
	cfg.STORE_INTERVAL, _ = strconv.ParseInt(*i, 10, 64)
	cfg.STORE_FILE = *f
	cfg.RESTORE, _ = strconv.ParseBool(*r)

	return &ServerFlagConfiguration{
		variablesMap: m,
		Cfg:          cfg,
	}
}

func (sc *ServerFlagConfiguration) Get(name string) (value string) {
	v := sc.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerFlagConfiguration) read() error {
	for k := range sc.variablesMap {
		v := os.Getenv(k)
		if v != "" {
			sc.variablesMap[k] = v
			log.Println("variable " + k + " presented in environment, value: " + sc.variablesMap[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + sc.variablesMap[k])

		}
	}
	return nil
}

func (sc *ServerFlagConfiguration) write() {
	// is not supported by implementation
	log.Fatal("method ServerFlagConfiguration.write is not supported by implementation")
}
func (sc *ServerFlagConfiguration) Update() *ServerFlagConfiguration {
	err := sc.read()
	if err != nil {
		sc.write()
	}

	return sc
}

func (sc *ServerFlagConfiguration) GetBool(name string) (value bool) {
	v := sc.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal("bool parsing error value" + v)
	}
	return b
}
func (sc *ServerFlagConfiguration) GetInt(name string) (value int64) {
	v := sc.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
		return
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatal("int parsing error value" + v)
	}
	return i
}

func copyMap[K, V comparable](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}
