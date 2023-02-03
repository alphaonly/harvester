package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AgentEnvConfiguration struct {
	variables *map[string]string
}

func NewAgentEnvConfiguration() *Configuration {

	v := copyMap(agentDefaults)

	c := Configuration(&AgentEnvConfiguration{
		variables: &v,
	})
	return &c
}

func (ac *AgentEnvConfiguration) Get(name string) (value string) {
	v := (*(*ac).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentEnvConfiguration) GetBool(name string) (value bool) {
	v := (*(*ac).variables)[name]
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
	v := (*(*ac).variables)[name]
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
	for k := range *(*ac).variables {
		v := os.Getenv(k)
		if v != "" {
			(*(*ac).variables)[k] = v
			log.Println("variable " + k + " presented in environment, value: " + (*(*ac).variables)[k])
		} else {
			log.Println("variable " + k + " not presented in environment, remains default: " + (*(*ac).variables)[k])
		}
	}
	return nil
}

func (ac *AgentEnvConfiguration) write() {
	for k := range *(*ac).variables {
		if v := (*(*ac).variables)[k]; v != "" {
			err := os.Setenv(k, v)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("variable" + k + " has no value to write")
		}
	}
}
func (ac *AgentEnvConfiguration) Update() *Configuration {

	err := (*ac).read()
	if err != nil {
		fmt.Println((*ac).variables)
		(*ac).write()

	}
	c := Configuration(ac)

	return &c
}
func (ac *AgentEnvConfiguration) UpdateNotGiven(fromConf *Configuration) {
	//Type assertion
	_fromConf := *fromConf

	switch fc := _fromConf.(type) {
	case *AgentFlagConfiguration:
		{
			for k := range *ac.variables {
				if (*ac.variables)[k] == agentDefaults[k] &&
					((*(*fc).variables)[k] != agentDefaults[k]) {
					(*ac.variables)[k] = (*(*fc).variables)[k]
					log.Println("variable " + k + " updated from flags, value:" + (*ac.variables)[k])
				}
			}
		}
	case *AgentEnvConfiguration:
		{
			for k := range *ac.variables {
				if (*ac.variables)[k] == agentDefaults[k] &&
					((*(*fc).variables)[k] != agentDefaults[k]) {
					(*ac.variables)[k] = (*(*fc).variables)[k]
					log.Println("variable " + k + " updated from flags, value:" + (*ac.variables)[k])
				}
			}
		}

	default:
		log.Fatal("UpdateNotGiven illegal type assertion")
	}

}

type ServerEnvConfiguration struct {
	variables *map[string]string
}

func NewServerEnvConfiguration() *Configuration {
	v := copyMap(serverDefaults)
	c := Configuration(&ServerEnvConfiguration{
		variables: &v,
	})
	return &c
}

func (sc *ServerEnvConfiguration) Get(name string) (value string) {
	v := (*(*sc).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerEnvConfiguration) read() error {
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

func (sc *ServerEnvConfiguration) write() {
	for k := range *(*sc).variables {
		if v := (*(*sc).variables)[k]; v != "" {
			err := os.Setenv(k, v)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("variable" + k + " has no value to write")
		}
	}
}
func (sc *ServerEnvConfiguration) Update() *Configuration {
	err := (*sc).read()
	if err != nil {
		(*sc).write()
	}
	c := Configuration(sc)
	return &c
}

func (sc *ServerEnvConfiguration) UpdateNotGiven(fromConf *Configuration) {

	//Type assertion
	_fromConf := *fromConf

	switch fc := _fromConf.(type) {
	case *ServerFlagConfiguration:
		{
			for k := range *sc.variables {
				if (*sc.variables)[k] == serverDefaults[k] &&
					((*(*fc).variables)[k] != serverDefaults[k]) {
					(*sc.variables)[k] = (*(*fc).variables)[k]
					log.Println("variable " + k + " updated from flags, value:" + (*sc.variables)[k])
				}
			}
		}
	case *ServerEnvConfiguration:
		{
			for k := range *sc.variables {
				if (*sc.variables)[k] == serverDefaults[k] &&
					((*(*fc).variables)[k] != serverDefaults[k]) {
					(*sc.variables)[k] = (*(*fc).variables)[k]
					log.Println("variable " + k + " updated from flags, value:" + (*sc.variables)[k])
				}
			}
		}

	default:
		log.Fatal("UpdateNotGiven illegal type assertion")
	}

}

func (sc *ServerEnvConfiguration) GetBool(name string) (value bool) {
	v:= (*(*sc).variables)[name]
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
	v:= (*(*sc).variables)[name]
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
