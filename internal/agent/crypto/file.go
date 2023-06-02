package crypto

import (
	"bufio"
	"errors"
	"log"
	"os"

	"github.com/alphaonly/harvester/internal/configuration"
)

func ReadPublicKeyFile(configuration *configuration.AgentConfiguration) (*bufio.Reader, error) {
	if configuration.CryptoKey == "" {
		mess := "path to given public key file is not defined"
		log.Println(mess)
		return nil, errors.New(mess)
	}
	//Reading file with rsa key from os
	file, err := os.OpenFile(configuration.CryptoKey, os.O_RDONLY, 0777)
	if err != nil {
		log.Printf("error:file %v  is not read", file)
		return nil, err
	}

	//put data to read buffer
	return bufio.NewReader(file), nil

}
