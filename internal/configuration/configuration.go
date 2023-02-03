package configuration

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

	"SCHEME":        "http",
	"USE_JSON":      "true",
	"COMPRESS_TYPE": "",
}
