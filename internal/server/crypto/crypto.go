package crypto

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
	"github.com/alphaonly/harvester/internal/common/logging"
	"github.com/alphaonly/harvester/internal/configuration"
)

type RSA struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	certBytes  []byte
	err        error
}

func NewRSA(serialNumber int64, cfg *configuration.ServerConfiguration) cryptoCommon.ServerCertificateManager {
	if cfg.CryptoKey == "" {
		return nil
	}
	ca := &x509.Certificate{

		SerialNumber: big.NewInt(serialNumber),

		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},

		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},

		NotBefore: time.Now(),

		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,

		DNSNames: []string{"localhost"},
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	logging.LogFatal(err)

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	logging.LogFatal(err)

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	cert := &x509.Certificate{

		SerialNumber: big.NewInt(serialNumber),

		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},

		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},

		NotBefore: time.Now(),

		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,

		DNSNames: []string{"localhost"},
	}

	var certPrivateKey *rsa.PrivateKey
	// var err error
	//Analyze privateKeyData if it given create certificate based on it
	if cfg.CryptoKey != "" {
		//read file with private key
		file, err := os.OpenFile(cfg.CryptoKey, os.O_RDONLY, 0777)
		logging.LogFatal(err)

		//make new instance of RSA to get private key data
		r := &RSA{}
		privateKeyBuf := r.Receive(cryptoCommon.PRIVATE, file).GetPrivate()
		logging.LogFatal(r.Error())

		//parsing PEM decoded
		certPrivateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBuf.Bytes())
		logging.LogFatal(err)
	}
	//if privateKeyData was not given, then generate  private key
	if certPrivateKey == nil {
		certPrivateKey, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			log.Fatal(err)
		}
	}
	// create cert x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivateKey.PublicKey, certPrivateKey)
	logging.LogFatal(err)

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivateKey),
	})

	// serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	// logging.LogFatal(err)

	cm := &RSA{privateKey: certPrivateKey,
		publicKey: &certPrivateKey.PublicKey,
		certBytes: certBytes}

	MakeCryptoFiles("/rsa/", cfg, cm)

	return cm
}

func (r *RSA) GetPublic() *bytes.Buffer {
	if r.Error() != nil {
		logging.LogPrintln(r.Error())
		return nil
	}
	b := x509.MarshalPKCS1PublicKey(r.publicKey)

	return bytes.NewBuffer(b)
}
func (r *RSA) GetPrivate() *bytes.Buffer {
	if r.Error() != nil {
		logging.LogPrintln(r.Error())
		return nil
	}

	b := x509.MarshalPKCS1PrivateKey(r.privateKey)
	return bytes.NewBuffer(b)
}

// Sends encoded in PEM cert or keys in writer
func (r *RSA) Send(dataType cryptoCommon.KeyType, b io.Writer) cryptoCommon.ServerCertificateManager {
	if r.Error() != nil {
		return r
	}
	switch dataType {
	case cryptoCommon.CERTIFICATE:
		{
			// encode certificate and public key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: r.certBytes,
			})
			logging.LogPrintln(err)
			r.err = err
		}
	case cryptoCommon.PUBLIC:
		{
			// encode  public key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(r.publicKey),
			})
			logging.LogPrintln(err)
			r.err = err

		}
	case cryptoCommon.PRIVATE:
		{
			// encode  private key in PEM format
			err := pem.Encode(b, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(r.privateKey),
			})
			logging.LogPrintln(err)
			r.err = err
		}
	default:
		r.err = errors.New("unknown data type to publish")
	}
	return r
}

func (r *RSA) Receive(dataType cryptoCommon.KeyType, buf io.Reader) cryptoCommon.ServerCertificateManager {
	if r.err != nil {
		return r
	}

	bytesPEM := make([]byte, 4096)
	_, err := buf.Read(bytesPEM)
	if err != nil {
		log.Println(err)
		r.err = err
		return &RSA{err: err}
	}

	block, _ := pem.Decode(bytesPEM)

	switch dataType {
	case cryptoCommon.CERTIFICATE:
		{
			// decode certificate and public key in PEM format
			if block == nil {
				r.err = errors.New("certificate is not found")
				logging.LogPrintln(err)
				return r
			}
			r.certBytes = block.Bytes
		}
	case cryptoCommon.PUBLIC:
		{
			// decode   public key in PEM format
			if block == nil {
				r.err = errors.New("public key is not found")
				logging.LogPrintln(err)
				return r
			}
			r.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			r.err = err
			logging.LogPrintln(err)
		}
	case cryptoCommon.PRIVATE:
		{
			// decode  private key in PEM format
			if block == nil {
				r.err = errors.New("private key is not found")
				logging.LogPrintln(err)
				return r
			}

			r.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			r.err = err
			logging.LogPrintln(err)

		}
	default:
		{
			r.err = errors.New("unknown data type to publish")
			logging.LogPrintln(err)
			return r

		}
	}
	return r
}

// Decrypt -  Decrypts in data and return it to out
func (r *RSA) DecryptData(in []byte) []byte {
	var decryptedBytes []byte

	//message length
	msgLen := len(in)

	//picked hash function
	hash := sha256.New()
	//message length for one interation
	step := r.privateKey.PublicKey.Size()

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedPart, err := rsa.DecryptOAEP(
			hash,
			rand.Reader,
			r.privateKey,
			in[start:finish],
			// in,
			nil)
		if err != nil {
			r.err = err
			logging.LogPrintln(err)
			return nil
		}
		decryptedBytes = append(decryptedBytes, decryptedPart...)

	}
	return decryptedBytes
}

func (r *RSA) Error() error {
	return r.err
}

func (r *RSA) IsError() bool {
	return r.err != nil
}

// MakeCryptoFiles - Makes files from data in certificate manager to /rsa/ folder
func MakeCryptoFiles(subFolder string, cfg *configuration.ServerConfiguration, cm cryptoCommon.ServerCertificateManager) {
	prefPath, _ := os.Getwd()

	//save cert and  keys in another folder
	pathCert := prefPath + "/" + subFolder + "/cert.rsa"
	pathPub := prefPath + "/" + subFolder + "/public.rsa"
	pathPriv := prefPath + "/" + subFolder + "/private.rsa"

	fileMap := cryptoCommon.DataTypeMap{
		pathCert: cryptoCommon.CERTIFICATE,
		pathPub:  cryptoCommon.PUBLIC,
		pathPriv: cryptoCommon.PRIVATE}
	if cm == nil {
		cm = NewRSA(6996, &configuration.ServerConfiguration{CryptoKey: pathPriv})
	}
	for fileName, dataType := range fileMap {
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
		logging.LogFatal(err)
		fileWriter := bufio.NewWriter(file)
		//Send certificate manager data to file writer
		cm.Send(dataType, fileWriter)
		logging.LogFatal(cm.Error())
		//write down data to file
		err = fileWriter.Flush()
		logging.LogFatal(err)

		//close file
		err = file.Close()
		logging.LogFatal(err)
	}

}
