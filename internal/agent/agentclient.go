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
		for tries := 1; tries <= c.Retries; tries++ {
			response, err := c.Client.Do(r)
			if err == nil {
				return response, nil
			}
			if c.Retries > 1 {
				log.Printf("Request error: %v", err)
				log.Printf("retry %v time...", tries)
				time.Sleep(c.RetryPause)
			}
		}
		log.Fatalf("agent gave up after %v attempts,exit", c.Retries)
	}

	return nil, errors.New("retries int was not noticed")
}
