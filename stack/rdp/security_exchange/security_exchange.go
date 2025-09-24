package securityexchange

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
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
	prefix := "rdp: sec-exchange: write: %w"

	if err := r.PrepareEncryptedClientRandom(cert); err != nil {
		return fmt.Errorf(prefix, err)
	}

	r.Length = uint32(len(r.EncryptedClientRandom) + 8)

	var buff bytes.Buffer

	ser, err := core.Serialize(r)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := binary.Write(&buff, binary.LittleEndian, r.EncryptedClientRandom); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(make([]byte, 8)); err != nil {
		return fmt.Errorf(prefix, err)
	}

	sdr := mcs.NewSendDataRequest(userID)

	if err := sdr.Write(buff.Bytes(), stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (r *Request) PrepareEncryptedClientRandom(cert certs.Certificate) error {
	prefix := "sec exhacnge: encr client random: %w"

	// Getting income data from server sertificate
	modulus, E := cert.TargetCertifacate.PublicKey()

	if len(modulus) == 0 || E == 0 {
		return fmt.Errorf(prefix, fmt.Errorf("invalid rsa key: exp=%d len(mod)=%d", E, len(modulus)))
	}

	modulus = bytes.TrimRight(modulus, "\x00")
	modulusLen := len(modulus)

	// Forming client random
	clientRandom := make([]byte, 32)

	if _, err := rand.Read(clientRandom); err != nil {
		return fmt.Errorf(prefix, err)
	}

	// Preparing income args
	modulus = core.Reverse(modulus)

	N := new(big.Int).SetBytes(modulus)
	m := make([]byte, modulusLen)

	copy(m, clientRandom)

	m = core.Reverse(m)
	M := new(big.Int).SetBytes(m)

	// ClientRandom (LE->BE->big.Int) ^ E mod Modulus (LE->BE->big.Int) = Result (big.Int->BE->LE)
	C := new(big.Int).Exp(M, big.NewInt(int64(E)), N)

	r.EncryptedClientRandom = core.Reverse(C.FillBytes(make([]byte, modulusLen)))

	return nil
}
