package agent

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type AgentClient struct {
	Client     *http.Client
	Retries    int
	RetryPause time.Duration
}

func (c AgentClient) DoWithRetry(r *http.Request) (*http.Response, error) {
	if c.Retries > 0 {
		for tries := 0; tries < c.Retries; tries++ {
			response, err := c.Client.Do(r)
			if err == nil {
				return response, nil
			}
			if c.Retries > 1 {
				log.Println("Request error, retry")
				time.Sleep(c.RetryPause)
			}
		}
	}
	return nil, errors.New("retries int was not noticed")
}
