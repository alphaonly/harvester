package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/alphaonly/harvester/internal/agent"
	"github.com/alphaonly/harvester/internal/agent/crypto"
	"github.com/alphaonly/harvester/internal/common"
	"github.com/alphaonly/harvester/internal/common/logging"

	conf "github.com/alphaonly/harvester/internal/configuration"
	"github.com/go-resty/resty/v2"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	//Build tags
	common.PrintBuildTags(buildVersion, buildDate, buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Configuration parameters
	ac := conf.NewAgentConf(conf.UpdateACFromEnvironment, conf.UpdateACFromFlags)
	//resty client
	client := resty.New().SetRetryCount(10)

	// if ac.CryptoKey != "" {

		// file, err := os.OpenFile(ac.CryptoKey, os.O_RDONLY, 0777)
		// logging.LogFatal(err)
		// reader := bufio.NewReader(file)
		// b := make([]byte, 4096)
		// _, err = reader.Read(b)
		// logging.LogFatal(err)
		// publicKey, _ := pem.Decode(b)

		// file, err := os.OpenFile(ac.CryptoCert, os.O_RDONLY, 0777)
		// logging.LogFatal(err)
		// reader := bufio.NewReader(file)
		// b := make([]byte, 4096)
		// _, err = reader.Read(b)
		// logging.LogFatal(err)
		// certBytes, _ := pem.Decode(b)

		// cert2 := &x509.Certificate{

		// 	SerialNumber: big.NewInt(1658),
		// 	//
		// 	Subject: pkix.Name{
		// 		Organization: []string{"Yandex.Praktikum"},
		// 		Country:      []string{"RU"},
		// 	},
		// 	//
		// 	IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// 	//
		// 	NotBefore: time.Now(),
		// 	//
		// 	NotAfter:     time.Now().AddDate(10, 0, 0),
		// 	SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// 	//
		// 	//
		// 	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		// 	KeyUsage:    x509.KeyUsageDigitalSignature,
		// 	PublicKey:   publicKey,
		// }

		// certFromFile, err := x509.ParseCertificate(certBytes.Bytes)
		// logging.LogFatal(err)

		// tlsCert := tls.Certificate{Certificate: [][]byte{certFromFile.Raw}, Leaf: certFromFile}

		// client.SetCertificates(tlsCert)
		// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	// }

	//load public key pem file
	buf, err := crypto.ReadPublicKeyFile(ac)
	logging.LogFatal(err)

	//get public key for cm
	cm := crypto.NewRSA(ac).ReceivePublic(buf)
	logging.LogFatal(cm.Error())
	//Run agent
	agent.NewAgent(ac, client, cm).Run(ctx)
	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	select {
	case <-channel:
		log.Print("Agent shutdown by os signal")
	case <-ctx.Done():
		log.Print("Agent shutdown by cancelled context")
	}
}
