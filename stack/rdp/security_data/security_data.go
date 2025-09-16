package securitydata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/mcs"
	"github.com/code-by-meal/go-rdp/stack/sec"
)

// Request

type Request struct {
	Flags         uint16 `order:"l"`
	FlagsHi       uint16 `order:"l"`
	DataSignature [8]byte
	EncryptedData []byte
}

func NewRequest(flags uint16) *Request {
	return &Request{
		Flags: flags,
	}
}

func (r *Request) Write(stream io.Writer, initiator uint16, sessionKey *sec.SessionKey, dataPlain []byte) error {
	prefix := "rdp: security data: write: %w"

	seq := 0
	signature := sec.MAC8(sessionKey.MAC8[:], uint32(seq), []byte("ENCRYPTED DATA"))
	encryptedData := make([]byte, len(dataPlain))
	buff := new(bytes.Buffer)

	copy(encryptedData, dataPlain)

	sessionKey.RC4Out.XORKeyStream(encryptedData, encryptedData)

	if err := core.WriteSingleAny(buff, &r.Flags, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.WriteSingleAny(buff, &r.FlagsHi, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(signature[:]); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(encryptedData); err != nil {
		return fmt.Errorf(prefix, err)
	}

	sdr := mcs.NewSendDataRequest(initiator)

	if err := sdr.Write(buff.Bytes(), stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

// Response

type Response struct{}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Read(stream io.Reader) error {
	return nil
}
