package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AgentEnvConfiguration struct {
	variablesMap map[string]string
	Cfg          AgentCfg
}

func NewAgentEnvConfiguration() *AgentEnvConfiguration {

	v := copyMap(agentDefaults)
	cfg := UnMarshalAgentDefaults(AgentDefaultJSON)

	ac := &AgentEnvConfiguration{
		variablesMap: v,
		Cfg:          cfg,
	}
	return ac.update()

}

func (ac *AgentEnvConfiguration) Get(name string) (value string) {
	v := ac.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentEnvConfiguration) GetBool(name string) (value bool) {
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
func (ac *AgentEnvConfiguration) GetInt(name string) (value int64) {
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

func (ac *AgentEnvConfiguration) read() error {
	for k := range ac.variablesMap {
		v := os.Getenv(k)
		if v != "" {
			ac.variablesMap[k] = v
			log.Println("variable " + k + " presented in environment, value: " + ac.variablesMap[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + ac.variablesMap[k])
		}
	}
	return nil
}

func (ac *AgentEnvConfiguration) write() {
	for k := range ac.variablesMap {
		if v := ac.variablesMap[k]; v != "" {
			err := os.Setenv(k, v)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("variable" + k + " has no value to write")
		}
	}
}
func (ac *AgentEnvConfiguration) update() *AgentEnvConfiguration {
	acv := *ac
	err := acv.read()
	if err != nil {
		fmt.Println(acv.variablesMap)
		acv.write()

	}
	return ac
}
func (ac *AgentEnvConfiguration) UpdateNotGiven(fromConf *AgentFlagConfiguration) {
	acv := *ac
	fc := *fromConf
	for k := range ac.variablesMap {
		if ac.variablesMap[k] == agentDefaults[k] &&
			(fc.variablesMap[k] != agentDefaults[k]) {
			ac.variablesMap[k] = fc.variablesMap[k]
			log.Println("variable " + k + " updated from flags, value:" + acv.variablesMap[k])
		}
	}

}

type ServerEnvConfiguration struct {
	variablesMap map[string]string
	Cfg          ServerCfg
}

func NewServerEnvConfiguration() *ServerEnvConfiguration {
	v := copyMap(serverDefaults)

	cfg := UnMarshalServerDefaults(ServerDefaultJSON)

	sc := ServerEnvConfiguration{
		variablesMap: v,
		Cfg:          cfg,
	}
	sc.update()
	return &sc

}

func (sc *ServerEnvConfiguration) Get(name string) (value string) {
	v := sc.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerEnvConfiguration) read() error {
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

func (sc *ServerEnvConfiguration) write() {
	for k := range sc.variablesMap {
		if v := sc.variablesMap[k]; v != "" {
			err := os.Setenv(k, v)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("variable" + k + " has no value to write")
		}
	}
}
func (sc *ServerEnvConfiguration) update() *ServerEnvConfiguration {
	err := sc.read()
	if err != nil {
		sc.write()
	}
	return sc
}

func (sc *ServerEnvConfiguration) UpdateNotGiven(fromConf *ServerFlagConfiguration) {
	fc := *fromConf
	for k := range sc.variablesMap {
		if sc.variablesMap[k] == serverDefaults[k] &&
			fc.variablesMap[k] != serverDefaults[k] {
			sc.variablesMap[k] = fc.variablesMap[k]
			log.Println("variable " + k + " updated from flags, value:" + sc.variablesMap[k])
		}
	}
}

func (sc *ServerEnvConfiguration) GetBool(name string) (value bool) {

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
func (sc *ServerEnvConfiguration) GetInt(name string) (value int64) {
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

// // implementation check
// var ac Configuration = &AgentEnvConfiguration{}
// var sc Configuration = &ServerEnvConfiguration{}
