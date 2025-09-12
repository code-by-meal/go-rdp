package certs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
)

// Certificate version
type ChainVersion uint32

const (
	Vresion1 ChainVersion = 1
	Version2 ChainVersion = 2
)

type Cert interface {
	Verify() bool
	PubclicKey() ([]byte, uint32)
	Read(io.Reader) error
}

type Certificate struct {
	DwVersion         ChainVersion `order:"l"`
	Raw               []byte
	TargetCertifacate Cert
}

func NewCertificate(buff *bytes.Buffer) (*Certificate, error) {
	c := &Certificate{}

	if err := core.Unserialize(buff, c); err != nil {
		return c, fmt.Errorf("certs: new cert: %w", err)
	}

	return c, nil
}

func (c *Certificate) Proccess() error {
	return nil
}
