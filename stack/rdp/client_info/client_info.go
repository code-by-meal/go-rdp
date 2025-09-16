package clientinfo

import (
	"io"

	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
)

type Request struct{}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Write(stream io.Writer, proto nego.NegoProtocol) error {
	prefix := "rdp: client-info: write: %w"
	_ = prefix

	switch proto {
	case nego.RDP:

	case nego.Hybrid:
	}

	return nil
}
