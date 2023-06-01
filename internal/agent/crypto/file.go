package crypto

import (
	"bufio"
	"log"
	"os"

	"github.com/alphaonly/harvester/internal/configuration"
)

func ReadPublicKeyFile(configuration *configuration.AgentConfiguration) (*bufio.Reader, error) {
	if configuration.CryptoKey == "" {
		log.Println("path to given public key file is not defined")
		return nil, nil
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
