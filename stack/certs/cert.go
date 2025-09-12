package certs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
)

// Certificate version
type ChainVersion uint32

const (
	Version1 ChainVersion = 1
	Version2 ChainVersion = 2
)

type Cert interface {
	Verify() bool
	PublicKey() ([]byte, uint32)
	Read(io.Reader) error
}

type Certificate struct {
	DwVersion         ChainVersion `order:"l"`
	Raw               []byte
	TargetCertifacate Cert
}

func NewCertificate(buff *bytes.Buffer) (*Certificate, error) {
	c := &Certificate{}
	prefix := "certs: new cert: %w"

	if err := core.Unserialize(buff, c); err != nil {
		return c, fmt.Errorf(prefix, err)
	}

	c.Raw = buff.Bytes()

	switch c.DwVersion {
	case Version2:
		c.TargetCertifacate = NewX509()

		log.Dbg("<d>Certificate type</> : <i>X509</>")
	case Version1:
		c.TargetCertifacate = NewPropietary()

		log.Dbg("<d>Certificate type</> : <i>PROPIETARY</>")
	default:
		return nil, fmt.Errorf(prefix, fmt.Errorf("unknown version of certificate: %d", c.DwVersion))
	}

	return c, nil
}
