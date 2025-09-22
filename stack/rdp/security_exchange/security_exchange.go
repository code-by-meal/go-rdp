package securityexchange

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/certs"
	"github.com/code-by-meal/go-rdp/stack/mcs"
)

type Request struct {
	Flags                 uint16 `order:"l"`
	FlagsHi               uint16 `order:"l"`
	Length                uint32 `order:"l"`
	EncryptedClientRandom []byte
}

func NewRequest() *Request {
	return &Request{
		Flags:   0x0201,
		FlagsHi: 0x0,
	}
}

func (r *Request) Write(stream io.Writer, userID uint16, cert certs.Certificate) error {
	modulus, exp := cert.TargetCertifacate.PublicKey()
	prefix := "rdp: sec-exchange: write: %w"

	if len(modulus) == 0 || exp == 0 {
		return fmt.Errorf(prefix, fmt.Errorf("empty rsa key"))
	}

	n := new(big.Int).SetBytes(modulus)
	publicKey := &rsa.PublicKey{N: n, E: int(exp)}
	clientRandom := make([]byte, 32)

	if _, err := rand.Read(clientRandom); err != nil {
		return fmt.Errorf(prefix, err)
	}

	enc, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, clientRandom)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	r.Length = uint32(len(enc))
	r.EncryptedClientRandom = enc

	var buff bytes.Buffer

	ser, err := core.Serialize(r)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(r.EncryptedClientRandom); err != nil {
		return fmt.Errorf(prefix, err)
	}

	sdr := mcs.NewSendDataRequest(userID)

	if err := sdr.Write(buff.Bytes(), stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
