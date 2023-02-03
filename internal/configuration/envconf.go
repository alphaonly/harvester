package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AgentEnvConfiguration struct {
	variables map[string]string
}

func NewAgentEnvConfiguration() *AgentEnvConfiguration {

	v := copyMap(agentDefaults)
	ac := &AgentEnvConfiguration{
		variables: v,
	}
	return ac.update()

}

func (ac *AgentEnvConfiguration) Get(name string) (value string) {
	v := ac.variables[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentEnvConfiguration) GetBool(name string) (value bool) {
	v := ac.variables[name]
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
	v := ac.variables[name]
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
	for k := range ac.variables {
		v := os.Getenv(k)
		if v != "" {
			ac.variables[k] = v
			log.Println("variable " + k + " presented in environment, value: " + ac.variables[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + ac.variables[k])
		}
	}
	return nil
}

func (ac *AgentEnvConfiguration) write() {
	for k := range ac.variables {
		if v := ac.variables[k]; v != "" {
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
		fmt.Println(acv.variables)
		acv.write()

	}
	return ac
}
func (ac *AgentEnvConfiguration) UpdateNotGiven(fromConf *AgentFlagConfiguration) {
	acv := *ac
	fc := *fromConf
	for k := range ac.variables {
		if ac.variables[k] == agentDefaults[k] &&
			(fc.variables[k] != agentDefaults[k]) {
			ac.variables[k] = fc.variables[k]
			log.Println("variable " + k + " updated from flags, value:" + acv.variables[k])
		}
	}

}

type ServerEnvConfiguration struct {
	variables map[string]string
}

func NewServerEnvConfiguration() *ServerEnvConfiguration {
	v := copyMap(serverDefaults)
	sc := ServerEnvConfiguration{
		variables: v,
	}
	sc.update()
	return &sc

}

func (sc *ServerEnvConfiguration) Get(name string) (value string) {
	v := sc.variables[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerEnvConfiguration) read() error {
	for k := range sc.variables {
		v := os.Getenv(k)
		if v != "" {
			sc.variables[k] = v
			log.Println("variable " + k + " presented in environment, value: " + sc.variables[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + sc.variables[k])

		}
	}
	return nil
}

func (sc *ServerEnvConfiguration) write() {
	for k := range sc.variables {
		if v := sc.variables[k]; v != "" {
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
	for k := range sc.variables {
		if sc.variables[k] == serverDefaults[k] &&
			fc.variables[k] != serverDefaults[k] {
			sc.variables[k] = fc.variables[k]
			log.Println("variable " + k + " updated from flags, value:" + sc.variables[k])
		}
	}
}

func (sc *ServerEnvConfiguration) GetBool(name string) (value bool) {

	v := sc.variables[name]
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
	v := sc.variables[name]
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
