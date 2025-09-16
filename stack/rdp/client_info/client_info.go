package clientinfo

import (
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
	securitydata "github.com/code-by-meal/go-rdp/stack/rdp/security_data"
	"github.com/code-by-meal/go-rdp/stack/sec"
)

type Request struct {
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Write(stream io.Writer, proto nego.NegoProtocol, intiator uint16, sessionKeys *sec.SessionKey) error {
	prefix := "rdp: client-info: write: %w"

	switch proto { // nolint
	case nego.RDP:
		r := securitydata.NewRequest(0x0848)

		if err := r.Write(stream, intiator, sessionKeys, []byte("ENCRYPTED DATA")); err != nil {
			return fmt.Errorf(prefix, err)
		}
	default:
		return fmt.Errorf(prefix, fmt.Errorf("protocol %v are implemetned to proccess", proto))
	}

	return nil
}
