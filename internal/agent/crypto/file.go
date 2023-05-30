package crypto

import (
	"bufio"
	"errors"
	"log"
	"os"

	"github.com/alphaonly/harvester/internal/configuration"
)

func CheckCertificateFile(configuration *configuration.AgentConfiguration) (*bufio.Reader, error) {
	if configuration.CryptoKey == "" {
		log.Println("path to given public key file is not defined")
		return nil, errors.New("no given rsa file path")
	}
	//Reading file with rsa key from os
	file, err := os.OpenFile(configuration.CryptoKey, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Printf("error:file %v  is not read", file)
		return nil, err
	}
	//put data to read buffer
	return bufio.NewReader(file), nil

}
