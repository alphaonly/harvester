package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
	"github.com/alphaonly/harvester/internal/common/logging"
	"github.com/alphaonly/harvester/internal/configuration"
)

type RSA struct {
	publicKey *rsa.PublicKey
	err       error
}

func NewRSA(cfg *configuration.AgentConfiguration) cryptoCommon.AgentCertificateManager {
	if cfg.CryptoKey == "" {
		return &RSA{err: errors.New("path to public key is not defined in configuration")}
	}
	return &RSA{}
}

func (r *RSA) GetPublic() *bytes.Buffer {
	if r.Error() != nil {
		logging.LogPrintln(r.Error())
		return nil
	}
	b := x509.MarshalPKCS1PublicKey(r.publicKey)

	return bytes.NewBuffer(b)
}

// SetPublic receive public key from PEM format
func (r *RSA) ReceivePublic(buf io.Reader) cryptoCommon.AgentCertificateManager {
	if r.Error() != nil {
		return r
	}
	bytesPEM := make([]byte, 4096)
	_, err := buf.Read(bytesPEM)
	if err != nil {
		log.Println(err)
		r.err = err
		return r
	}
	// decode   public key in PEM format
	block, _ := pem.Decode(bytesPEM)
	if block == nil {
		r.err = errors.New("public key is not found")
		return r

	}
	r.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	logging.LogFatal(err)

	return r
}
func (r *RSA) Error() error {
	return r.err
}

// Encrypt -  Encrypts in data and return it to out
func (r *RSA) EncryptData(in []byte) *[]byte {
	var encryptedBytes []byte

	//message length
	msgLen := len(in)
	//picked hash function
	hash := sha256.New()
	//message length for one interation
	step := r.publicKey.Size() - 2*hash.Size() - 2

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedPart, err := rsa.EncryptOAEP(
			hash,
			rand.Reader,
			r.publicKey,
			in[start:finish],
			nil)
		if err != nil {
			r.err = err
			logging.LogPrintln(err)
			return nil
		}
		encryptedBytes = append(encryptedBytes, encryptedPart...)

	}
	return &encryptedBytes
}
