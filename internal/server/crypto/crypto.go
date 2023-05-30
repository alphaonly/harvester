package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"net"
	"time"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
	"github.com/alphaonly/harvester/internal/configuration"
)

type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	certBytes  []byte
}

func NewRSA(serialNumber int64, configuration *configuration.ServerConfiguration) cryptoCommon.CertificateManager {
	if !configuration.EnableHTTPS {
		return &RSA{}
	}
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(serialNumber),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	
	return &RSA{privateKey: privateKey,
		publicKey: &privateKey.PublicKey,
		certBytes: certBytes}
}

func (r *RSA) GetPublic() (*bytes.Buffer, error) {

	if r.certBytes == nil {
		return nil, errors.New("")
	}

	return bytes.NewBuffer(r.certBytes), nil
}
func (r *RSA) GetPrivate() (*bytes.Buffer, error) {

	return nil, nil
}
func (r RSA) Send(dataType cryptoCommon.DataType, b *bufio.Writer) (err error) {
	switch dataType.GetType() {
	case cryptoCommon.CERTIFICATE_TYPE:
		{
			// encode certificate and public key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: r.certBytes,
			})
			if err != nil {
				return err
			}
		}
	case cryptoCommon.PUBLIC_TYPE:
		{
			// encode  public key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(r.publicKey),
			})
			if err != nil {
				return err
			}
		}
	case cryptoCommon.PRIVATE_TYPE:
		{
			// encode  private key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(r.privateKey),
			})
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unknown data type to publish")
	}

	return nil
}

func (r *RSA) Receive(dataType cryptoCommon.DataType, buf *bufio.Reader) (cryptoCommon.CertificateManager, error) {

	 bytesPEM := make([]byte,4096)
	_, err := buf.Read(bytesPEM)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	switch dataType.GetType() {
	case cryptoCommon.CERTIFICATE_TYPE:
		{
			// decode certificate and public key in PEM format
			block, _ := pem.Decode(bytesPEM)
			if block == nil {
				log.Fatal("certificate is not found")
			}
			r.certBytes = block.Bytes
		}
	case cryptoCommon.PUBLIC_TYPE:
		{
			// decode   public key in PEM format
			block, _ := pem.Decode(bytesPEM)
			if block == nil {
				log.Fatal("public key is not found")
			}
			r.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}

		}
	case cryptoCommon.PRIVATE_TYPE:
		{
			// decode  private key in PEM format
			block, _ := pem.Decode(bytesPEM)
			if block == nil {
				log.Fatal("private key is not found")
			}

			r.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		return nil, errors.New("unknown data type to publish")
	}
	return r, nil
}
