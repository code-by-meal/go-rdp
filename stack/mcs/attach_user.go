package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/mcs/per"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type AttachUserRequest struct{}

func NewAttachUserRequest() *AttachUserRequest {
	return &AttachUserRequest{}
}

func (a *AttachUserRequest) Write(stream io.Writer) error {
	var buff bytes.Buffer

	prefix := "mcs: attach-user: write: %w"

	if err := per.WriteChoice(&buff, byte(AttachUserRequestT)<<2); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, &buff, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
