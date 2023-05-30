package crypto

import (
	"bufio"
	"bytes"
)

const (
	CERTIFICATE_TYPE = "certificate"
	PUBLIC_TYPE      = "public"
	PRIVATE_TYPE     = "private"
)

type DataType struct {
	typ string
}

func (d DataType) GetType() string { return d.typ }

func Certificate() DataType { return DataType{typ: CERTIFICATE_TYPE} }
func Private() DataType     { return DataType{typ: PRIVATE_TYPE} }
func Public() DataType      { return DataType{typ: PUBLIC_TYPE} }

type CertificateManager interface {
	GetPublic() (*bytes.Buffer, error)                                        //returns public key
	GetPrivate() (*bytes.Buffer, error)                                       //returns private key
	Send(dataType DataType, buf *bufio.Writer) error                          //sends crypto data cert, private, public keys to writer
	Receive(dataType DataType, buf *bufio.Reader) (CertificateManager, error) // receives crypto data cert, private, public keys from reader
}

type DataTypeMap map[string]func() DataType
