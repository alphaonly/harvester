package crypto

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"

	cryptoCommon "github.com/alphaonly/harvester/internal/common/crypto"
)

type RSA struct {
	publicKey *rsa.PublicKey
}

func NewRSA() cryptoCommon.CertificateManager { return &RSA{} }

func (r *RSA) GetPublic() (*bytes.Buffer, error) {

	if r.publicKey == nil {
		return nil, errors.New("")
	}

	return bytes.NewBuffer(r.publicKey.N.Bytes()), nil
}
func (r *RSA) GetPrivate() (*bytes.Buffer, error) {
	return nil, errors.New("getting private key is forbidden for agent")
}
func (r RSA) Send(dataType cryptoCommon.DataType, b *bufio.Writer) (err error) {
	return errors.New("sending rsa data is forbidden for agent")
}

func (r *RSA) Receive(dataType cryptoCommon.DataType, buf *bufio.Reader) (cryptoCommon.CertificateManager, error) {

	var bytesPEM []byte
	_, err := buf.Read(bytesPEM)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	switch dataType.GetType() {
	case cryptoCommon.CERTIFICATE_TYPE:
		{
			return nil, errors.New("receiving certificate data is forbidden for agent")
		}
	case cryptoCommon.PUBLIC_TYPE:
		{
			// decode   public key in PEM format
			block, _ := pem.Decode(bytesPEM)
			if block == nil {
				log.Fatal(errors.New("public key is not found"))
			}
			r.publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}

		}
	case cryptoCommon.PRIVATE_TYPE:
		{
			return nil, errors.New("receiving private key data is forbidden for agent")
		}
	default:
		return nil, errors.New("unknown data type to publish")
	}
	return r, nil
}
