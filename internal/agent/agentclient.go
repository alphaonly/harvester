package agent

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type AgentClient struct {
	Client     *http.Client
	Retries    int
	RetryPause time.Duration
}

func (c AgentClient) DoWithRetry(r *http.Request) (body []byte, err error) {
	if c.Retries > 0 {
		for tries := 1; tries <= c.Retries; tries++ {
			response, err := c.Client.Do(r)
			if err != nil {
				log.Printf("sending request:%v", r)
				log.Printf("sending response:%v", response)
				log.Printf("clientDo error:%v", err)
			}

			if err == nil {
				defer response.Body.Close()
				log.Println("agent:response from server:" + response.Status)
				bytes, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("read all response body error:%v", err)
				}
				log.Println("agent:body from server:" + string(bytes))
				return bytes, nil
			}
			if c.Retries > 1 {
				log.Printf("Request error: %v", err)
				log.Printf("retry %v time...", tries)
				time.Sleep(c.RetryPause)
			}

		}
		log.Fatalf("agent gave up after %v attempts,exit", c.Retries)
		return nil, err
	}

	return nil, errors.New("retries int was not noticed")
}
