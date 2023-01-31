package configuration

type Configuration interface {
	Update() *Configuration
	UpdateNotGiven(fromConf *Configuration)
	Get(name string) (value string)
	GetInt(name string) (value int64)
	GetBool(name string) (value bool)
	read() error
	write()
}

var serverDefaults = map[string]string{
	"SERVER_PORT":    "8080",
	"STORE_INTERVAL": "300",
	"STORE_FILE":     "/tmp/devops-metrics-db.json",
	"RESTORE":        "true",
}

var agentDefaults = map[string]string{
	"POLL_INTERVAL":   "2",
	"REPORT_INTERVAL": "10",
	"HOST":            "localhost",
	"PORT":            "8080",
	"SCHEME":          "http",
	"USE_JSON":        "false",
}
