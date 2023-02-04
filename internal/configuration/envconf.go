package configuration

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func (ac *AgentEnvConfiguration) get(name string) (value string) {
	v := ac.variablesMap[name]
	if v == "" {
		log.Println("no variable:" + name)
	}

	return v
}

func (ac *AgentEnvConfiguration) getBool(name string) (value bool) {
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
func (ac *AgentEnvConfiguration) getInt(name string) (value int64) {
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
	ac.Cfg.ADDRESS = ac.get("ADDRESS")
	ac.Cfg.SCHEME = ac.get("SCHEME")
	ac.Cfg.COMPRESS_TYPE = ac.get("COMPRESS_TYPE")
	ac.Cfg.POLL_INTERVAL = ac.getInt("POLL_INTERVAL")
	ac.Cfg.USE_JSON = ac.getBool("USE_JSON")

}

type ServerEnvConfiguration struct {
	variablesMap map[string]string
	Cfg          ServerCfg
}

func NewServerEnvConfiguration() *ServerEnvConfiguration {
	v := copyMap(serverDefaults)
	v["PORT"] = ":" + (strings.Split(v["ADDRESS"], ":"))[1]

	cfg := UnMarshalServerDefaults(ServerDefaultJSON)
	cfg.PORT = ":" + (strings.Split(cfg.ADDRESS, ":"))[1]

	sc := ServerEnvConfiguration{
		variablesMap: v,
		Cfg:          cfg,
	}
	sc.update()
	return &sc

}

func (sc *ServerEnvConfiguration) get(name string) (value string) {
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
		sc.write() //write defaults to environment
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

			//Парсим порт из адреса

		}
		if k == "ADDRESS" {
			sc.variablesMap["PORT"] = ":" + strings.Split(sc.variablesMap[k], ":")[1]
			sc.Cfg.PORT = sc.variablesMap["PORT"]
			log.Println("variable PORT updated from ADDRESS flags, value:" + sc.Cfg.PORT)
		}
	}
	sc.Cfg.STORE_FILE = sc.get("STORE_FILE")
	sc.Cfg.RESTORE = sc.getBool("RESTORE")
	sc.Cfg.STORE_INTERVAL = sc.getInt("STORE_INTERVAL")
	sc.Cfg.ADDRESS = sc.get("ADDRESS")

	if sc.get("PORT") != serverDefaults["PORT"] {
		sc.Cfg.PORT = sc.get("PORT")
		log.Println("variable PORT updated from ADDRESS flags, value:" + sc.Cfg.PORT)
	}

}

func (sc *ServerEnvConfiguration) getBool(name string) (value bool) {

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
func (sc *ServerEnvConfiguration) getInt(name string) (value int64) {
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
