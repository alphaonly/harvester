package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/alphaonly/harvester/internal/schema"
	"github.com/go-resty/resty/v2"
)

// create HTTP client without redirects support

func main() {

	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})
	httpc := resty.New().
		SetHostURL("http://127.0.0.1:8080").
		SetRedirectPolicy(redirPolicy)

	req := httpc.R().
		SetHeader("Content-Type", "application/json")

	// Вдруг на сервере уже есть значение, на всякий случай запросим.
	var result schema.MetricsJSON
	id := "GetSet" + strconv.Itoa(rand.Intn(256) )
	resp, err := req.
		SetBody(&schema.MetricsJSON{
			ID:    id,
			MType: "counter"}).
		SetResult(&result).
		Post("value/")
	if err != nil {
		log.Println("error")
		log.Println(resp.Body())
	}
	// dumpErr := false //suite.Assert().NoError(err, "Ошибка при попытке сделать запрос с получением значения counter")
	// var value0 int64
	switch resp.StatusCode() {
	case http.StatusOK:
		{
			if http.StatusOK != resp.StatusCode() {
				log.Printf("Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
			}
			if !strings.Contains(resp.Header().Get("Content-Type"), "application/json") {

				log.Println("Заголовок ответа Content-Type содержит несоответствующее значение")
			}

			if result.Delta == nil {
				log.Printf("Получено не инициализированное значение Delta '%q %s'", req.Method, req.URL)
			}
		}

	case http.StatusNotFound:
		{
		}

	default:
		// dumpErr = false
		log.Fatalf("Несоответствие статус кода %d ответа ожидаемому http.StatusNotFound или http.StatusOK в хендлере %q: %q", resp.StatusCode(), req.Method, req.URL)
		return
	}
	var value0 int64 = 800
	resp, err = req.
		SetBody(&schema.MetricsJSON{
			ID:    "PollCounter",
			MType: "counter",
			Delta: &value0,
		}).
		Post("update/")

	if http.StatusOK != resp.StatusCode() {
		log.Printf("Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
	}
	var value1 int64 = 800
	resp, err = req.
		SetBody(&schema.MetricsJSON{
			ID:    "PollCounter",
			MType: "counter",
			Delta: &value1,
		}).
		Post("update/")

	if http.StatusOK != resp.StatusCode() {
		log.Printf("Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
	}
	resp, err = req.
		SetBody(&schema.MetricsJSON{
			ID:    "PollCounter",
			MType: "counter"}).
		SetResult(&result).
		Post("/value/")

	if http.StatusOK != resp.StatusCode() {
		log.Printf("Несоответствие статус кода ответа ожидаемому в хендлере %q: %q ", req.Method, req.URL)
	}
	if !strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		log.Println("Заголовок ответа Content-Type содержит несоответствующее значение")
	}
	if result.Delta == nil {
		log.Printf("Получено не инициализированное значение Delta '%q %s'", req.Method, req.URL)
	}

}
