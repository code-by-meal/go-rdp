package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/mcs/per"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type ErrectDomainRequest struct {
	Type        PDUType
	SubHeight   uint16
	SubInterval uint16
}

func NewErrectDomainRequest() *ErrectDomainRequest {
	return &ErrectDomainRequest{
		Type:        ErectDomainRequestT,
		SubHeight:   0,
		SubInterval: 0,
	}
}

func (e *ErrectDomainRequest) Write(stream io.Writer) error {
	var buff bytes.Buffer

	prefix := "mcs: errect-domain: write: %w"

	if err := per.WriteChoice(&buff, (byte(ErectDomainRequestT)<<2 | 0)); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := per.WriteInteger(&buff, 0); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := per.WriteInteger(&buff, 0); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, &buff, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
