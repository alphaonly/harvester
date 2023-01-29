package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AgentConfiguration struct {
	variables *map[string]string
}

func NewAgentConfiguration() *Configuration {
	m := make(map[string]string)

	var c Configuration = &AgentConfiguration{
		variables: &m,
	}
	//default bucket of parameters and their values
	m["POLL_INTERVAL"] = "2"
	m["REPORT_INTERVAL"] = "3" //10
	m["HOST"] = "127.0.0.1"
	m["PORT"] = "8080"
	m["SCHEME"] = "http"
	m["USE_JSON"] = "false"
	return &c
}

func (ac *AgentConfiguration) Get(name string) (value string) {
	var v string = (*(*ac).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentConfiguration) GetBool(name string) (value bool) {
	var v string = (*(*ac).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal("bool parsing error value" + v)
	}
	return b
}
func (ac *AgentConfiguration) GetInt(name string) (value int64) {
	var v string = (*(*ac).variables)[name]
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

func (ac *AgentConfiguration) read() error {
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

func (ac *AgentConfiguration) write() {
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
func (ac *AgentConfiguration) Update() *Configuration {

	err := (*ac).read()
	if err != nil {
		fmt.Println((*ac).variables)
		(*ac).write()
	}
	var c Configuration = ac
	return &c
}

type ServerConfiguration struct {
	variables *map[string]string
}

func NewServerConfiguration() *Configuration {
	m := make(map[string]string)

	var c Configuration = &AgentConfiguration{
		variables: &m,
	}
	//default bucket of parameters and their values
	m["SERVER_PORT"] = "8080"
	m["STORE_INTERVAL"] = "10" //300
	//m["STORE_FILE"] = "/tmp/devops-metrics-db.json"
	m["STORE_FILE"] = "devops-metrics-db.json"
	m["RESTORE"] = "true"
	return &c
}

func (sc *ServerConfiguration) Get(name string) (value string) {
	var v string = (*(*sc).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (sc *ServerConfiguration) read() error {
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

func (sc *ServerConfiguration) write() {
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
func (sc *ServerConfiguration) Update() *Configuration {
	err := (*sc).read()
	if err != nil {
		(*sc).write()
	}
	var c Configuration = sc
	return &c
}
func (sc *ServerConfiguration) GetBool(name string) (value bool) {
	var v string = (*(*sc).variables)[name]
	if v == "" {
		log.Println("no variable:" + name)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatal("bool parsing error value" + v)
	}
	return b
}
func (sc *ServerConfiguration) GetInt(name string) (value int64) {
	var v string = (*(*sc).variables)[name]
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
var ac Configuration = &AgentConfiguration{}
var sc Configuration = &ServerConfiguration{}
