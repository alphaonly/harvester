package agent

import (
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type AgentClient struct {
	Client     *http.Client
	Retries    int
	RetryPause time.Duration
}

func (c AgentClient) DoWithRetry(r *http.Request) (body []byte, err error) {
	var mutex sync.Mutex

	if c.Retries > 0 {
		for tries := 1; tries <= c.Retries; tries++ {
			mutex.Lock()  
			response, err := c.Client.Do(r)
			if err != nil {
				log.Printf("sending request:%v", r)
				log.Printf("getting response:%v", response)
				log.Printf("clientDo error:%v", err)
			}

			if err == nil {
				defer response.Body.Close()
				log.Println("agent:response from server:" + response.Status)
				bytes, err := io.ReadAll(response.Body)
				if err != nil {
					log.Fatalf("read all response body error:%v", err)
				}
				return bytes, nil
			}
			mutex.Unlock()  
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
