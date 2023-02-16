package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

var counterHashMessage = func(id string, delta int64) []byte {
	return []byte(fmt.Sprintf("%s:counter:%d", id, delta))
}

var gaugeHashMessage = func(id string, value float64) []byte {
	return []byte(fmt.Sprintf("%s:gauge:%f", id, value))
}

func CounterHash(id string, delta int64, key []byte) ([]byte, error) {

	msg := counterHashMessage(id, delta)
	h := hmac.New(sha256.New, key)

	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil

}
func GaugeHash(id string, value float64, key []byte) ([]byte, error) {
	msg := gaugeHashMessage(id, value)
	h := hmac.New(sha256.New, key)

	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil

}
