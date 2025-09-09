package gcc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/mcs"
	"github.com/code-by-meal/go-rdp/stack/mcs/per"
)

type ConferenceCreateRequest struct {
	UserData []byte
}

func NewCCR(userData []byte) *ConferenceCreateRequest {
	return &ConferenceCreateRequest{
		UserData: userData,
	}
}

func (c *ConferenceCreateRequest) Write(stream io.Writer) error {
	prefix := "gcc-1.24: write: %w"
	buff, err := c.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	log.Dbg(buff.Bytes())

	log.Dbg([]byte(c.UserData))

	ci := mcs.NewConnectInitial(buff.Bytes())

	if err := ci.Write(stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (c *ConferenceCreateRequest) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	prefix := "gcc-1.24: serialize: %w"

	if err := per.WriteChoice(&buff, 0); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteOID(&buff, []byte{0, 0, 20, 124, 0, 1}); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteLength(&buff, len(c.UserData)+14); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteChoice(&buff, 0); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteSelection(&buff, 0x08); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteNumericString(&buff, "1", 1); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WritePadding(&buff, 1); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteNumberOfSet(&buff, 1); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteChoice(&buff, 0xc0); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteOctetString(&buff, "Duca", 4); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := per.WriteOctetString(&buff, string(c.UserData), 0); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}

type ConferenceCreateResponse struct{}
