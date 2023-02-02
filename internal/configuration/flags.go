package configuration

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type AgentFlagConfiguration struct {
	variables *map[string]string
}

// Аргументы агента:
// ADDRESS через флаг a=<ЗНАЧЕНИЕ>,
// REPORT_INTERVAL через флаг r=<ЗНАЧЕНИЕ>,
// POLL_INTERVAL через флаг p=<ЗНАЧЕНИЕ>.

func NewAgentFlagConfiguration() *Configuration {
	m := copyMap(agentDefaults)

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

	var c Configuration = &AgentFlagConfiguration{
		variables: &m,
	}
	return &c
}

func (ac *AgentFlagConfiguration) DefaultConf() *Configuration {

	var c Configuration = &ServerEnvConfiguration{
		variables: &agentDefaults,
	}
	return &c
}

func (ac *AgentFlagConfiguration) Get(name string) (value string) {
	var v = (*(*ac).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentFlagConfiguration) GetBool(name string) (value bool) {
	var v = (*(*ac).variables)[name]
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
	var v = (*(*ac).variables)[name]
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
func (ac *AgentFlagConfiguration) Update() *Configuration {
	log.Fatal("method AgentFlagConfiguration.Update is not supported by implementation")
	return nil
}

func (ac *AgentFlagConfiguration) UpdateNotGiven(fromConf *Configuration) {
	//Type assertion
	_fromConf := *fromConf

	switch fc := _fromConf.(type) {
	case *AgentFlagConfiguration:
		{
			for k := range *ac.variables {
				if (*ac.variables)[k] == agentDefaults[k] &&
					((*(*fc).variables)[k] != agentDefaults[k]) {
					(*ac.variables)[k] = (*(*fc).variables)[k]
				}
			}
		}
	case *AgentEnvConfiguration:
		{
			for k := range *ac.variables {
				if (*ac.variables)[k] == agentDefaults[k] &&
					((*(*fc).variables)[k] != agentDefaults[k]) {
					(*ac.variables)[k] = (*(*fc).variables)[k]
				}
			}
		}

	default:
		log.Fatal("UpdateNotGiven illegal type assertion")
	}

}

type ServerFlagConfiguration struct {
	variables *map[string]string
}

func NewServerFlagConfiguration() *Configuration {

	m := copyMap(serverDefaults)

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

	var c Configuration = &ServerFlagConfiguration{
		variables: &m,
	}

	return &c
}

func (sc *ServerFlagConfiguration) Get(name string) (value string) {
	var v = (*(*sc).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerFlagConfiguration) read() error {
	for k := range *(*sc).variables {
		v := os.Getenv(k)
		if v != "" {
			(*(*sc).variables)[k] = v
			log.Println("variable " + k + " presented in environment, value: " + (*(*sc).variables)[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + (*(*sc).variables)[k])

		}
	}
	return nil
}

func (sc *ServerFlagConfiguration) write() {
	// is not supported by implementation
	log.Fatal("method ServerFlagConfiguration.write is not supported by implementation")
}
func (sc *ServerFlagConfiguration) Update() *Configuration {
	err := (*sc).read()
	if err != nil {
		(*sc).write()
	}
	var c Configuration = sc
	return &c
}
func (sc *ServerFlagConfiguration) UpdateNotGiven(fromConf *Configuration) {

	//Type assertion
	_fromConf := *fromConf

	switch fc := _fromConf.(type) {
	case *AgentFlagConfiguration:
		{
			for k := range *sc.variables {
				if (*sc.variables)[k] == serverDefaults[k] &&
					((*(*fc).variables)[k] != serverDefaults[k]) {
					(*sc.variables)[k] = (*(*fc).variables)[k]
				}
			}
		}
	case *AgentEnvConfiguration:
		{
			for k := range *sc.variables {
				if (*sc.variables)[k] == serverDefaults[k] &&
					((*(*fc).variables)[k] != serverDefaults[k]) {
					(*sc.variables)[k] = (*(*fc).variables)[k]
				}
			}
		}

	default:
		log.Fatal("UpdateNotGiven illegal type assertion")
	}

}

func (sc *ServerFlagConfiguration) GetBool(name string) (value bool) {
	var v = (*(*sc).variables)[name]
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
	var v = (*(*sc).variables)[name]
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

// // implementation check
// var ac Configuration = &AgentFlagConfiguration{}
// var sc Configuration = &ServerFlagConfiguration{}
