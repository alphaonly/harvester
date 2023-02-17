package signchecker

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/alphaonly/harvester/internal/schema"
)

type Signer interface {
	IsValidSign(mj schema.MetricsJSON) bool
	CounterHash(id string, delta *int64) ([]byte, error)
	GaugeHash(id string, value *float64) ([]byte, error)
	Sign(mj *schema.MetricsJSON) (err error)
}

type CheckerSHA256 struct {
	key []byte
}

func NewSHA256(key string) Signer {
	return CheckerSHA256{key: make([]byte, len(key))}
}
func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func (c CheckerSHA256) IsValidSign(mj schema.MetricsJSON) (result bool) {

	var leftHash []byte
	var err error

	//mj hash came to confirm the sign
	leftHash,err= hex.DecodeString(mj.Hash)
	logFatal(err)
	//signing it again to update hash to compare
	err = c.Sign(&mj)
	logFatal(err)

	//get updated hash to compare with left(original)
	rightHash, err := hex.DecodeString(mj.Hash)
	logFatal(err)
	//compare
	return hmac.Equal(leftHash, rightHash)
}

var counterHashMessage = func(id string, delta *int64) []byte {
	return []byte(fmt.Sprintf("%s:counter:%d", id, *delta))
}

var gaugeHashMessage = func(id string, value *float64) []byte {
	return []byte(fmt.Sprintf("%s:gauge:%f", id, *value))
}

func (c CheckerSHA256) CounterHash(id string, delta *int64) ([]byte, error) {
	msg := counterHashMessage(id, delta)
	h := hmac.New(sha256.New, c.key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
func (c CheckerSHA256) GaugeHash(id string, value *float64) ([]byte, error) {
	msg := gaugeHashMessage(id, value)
	h := hmac.New(sha256.New, c.key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
func (c CheckerSHA256) Sign(mj *schema.MetricsJSON) (err error) {
	var hashBytes []byte

	switch mj.MType {
	case "gauge":
		{
			hashBytes, err = c.GaugeHash(mj.ID, mj.Value)
			return err
		}
	case "counter":
		{
			hashBytes, err = c.CounterHash(mj.ID, mj.Delta)
			return err
		}
	default:
		log.Panic("CheckJSON unknown type")
		return
	}

	mj.Hash = hex.EncodeToString(hashBytes)
	return nil
}
