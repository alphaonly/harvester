package crypto

import (
	"bytes"
	"io"
)

type KeyType string

const (
	CERTIFICATE KeyType = "certificate"
	PUBLIC      KeyType = "public"
	PRIVATE     KeyType = "private"
)

type DataTypeMap map[string]KeyType

// An interface that handles certificates and keys for server
type ServerCertificateManager interface {
	GetPublic() *bytes.Buffer                                         // returns public key decode from PEM
	GetPrivate() *bytes.Buffer                                        // returns private key
	Send(dataType KeyType, buf io.Writer) ServerCertificateManager    // sends crypto data cert, private, public keys to writer
	Receive(dataType KeyType, buf io.Reader) ServerCertificateManager // receives crypto data cert, private, public keys from reader
	Error() error                                                     // returns error if appeared
}

// An interface that handles certificates and keys for agent
type AgentCertificateManager interface {
	GetPublic() *bytes.Buffer                        // returns public key bytes decoded from PEM
	ReceivePublic(io.Reader) AgentCertificateManager // sets public key decode
	EncryptData(in []byte) (out *[]byte)            // encrypts in data and returns to out
	Error() error                                    // returns error if appeared
}
