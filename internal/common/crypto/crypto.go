package crypto

import (
	"bytes"
	"io"
)

type DataType string

const (
	CERTIFICATE DataType = "certificate"
	PUBLIC      DataType = "public"
	PRIVATE     DataType = "private"
)

type DataTypeMap map[string]DataType

// An interface that handles certificates and keys for agent and server
type ServerCertificateManager interface {
	GetPublic() *bytes.Buffer                                          // returns public key decode from PEM
	GetPrivate() *bytes.Buffer                                         // returns private key
	Send(dataType DataType, buf io.Writer) ServerCertificateManager    // sends crypto data cert, private, public keys to writer
	Receive(dataType DataType, buf io.Reader) ServerCertificateManager // receives crypto data cert, private, public keys from reader
	Error() error                                                      // returns error if appeared
}

// An interface that handles certificates and keys for agent and server
type AgentCertificateManager interface {
	GetPublic() *bytes.Buffer                        // returns public key bytes decoded from PEM
	ReceivePublic(io.Reader) AgentCertificateManager // sets public key decode
	Error() error                                    // returns error if appeared
}
