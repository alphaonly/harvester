package configuration

import (
	"encoding/json"
	"log"
)

var serverDefaults = map[string]string{
	"ADDRESS":        "localhost:8080",
	"STORE_INTERVAL": "300",
	"STORE_FILE":     "/tmp/devops-metrics-db.json",
	"RESTORE":        "true",
}

var agentDefaults = map[string]string{
	"POLL_INTERVAL":   "2",
	"REPORT_INTERVAL": "10",
	"ADDRESS":         "localhost:8080",
	"SCHEME":          "http",
	"USE_JSON":        "true",
	"COMPRESS_TYPE":   "",
}

var ServerDefaultJSON = `{"ADDRESS":"localhost:8080","STORE_INTERVAL": 300,"STORE_FILE":"/tmp/devops-metrics-db.json","RESTORE":true}`
var AgentDefaultJSON = `{"POLL_INTERVAL":2,"REPORT_INTERVAL":10,"ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":true}`

type ServerCfg struct {
	ADDRESS        string `json:"ADDRESS,omitempty"`
	STORE_INTERVAL int64  `json:"STORE_INTERVAL,omitempty"`
	STORE_FILE     string `json:"STORE_FILE,omitempty"`
	RESTORE        bool   `json:"RESTORE,omitempty"`
	PORT           string `json:"PORT,omitempty"` //additionally for listen and serve func
}
type AgentCfg struct {
	POLL_INTERVAL   int64  `json:"POLL_INTERVAL,omitempty"`
	REPORT_INTERVAL int64  `json:"REPORT_INTERVAL,omitempty"`
	ADDRESS         string `json:"ADDRESS,omitempty"`
	SCHEME          string `json:"SCHEME,omitempty"`
	USE_JSON        bool   `json:"USE_JSON ,omitempty"`
	COMPRESS_TYPE   string `json:"COMPRESS_TYPE,omitempty"`
}

func UnMarshalAgentDefaults(s string) AgentCfg {
	ac := AgentCfg{}
	err := json.Unmarshal([]byte(s), &ac)
	if err != nil {
		log.Fatal("cannot unmarshal agent configuration")
	}
	return ac
}
func UnMarshalServerDefaults(s string) ServerCfg {
	sc := ServerCfg{}
	err := json.Unmarshal([]byte(s), &sc)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return sc

}
