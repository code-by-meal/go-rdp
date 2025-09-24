package certs

import (
	"bytes"
	"fmt"
)

type X509 struct {
}

func NewX509() *X509 {
	return &X509{}
}

func (x *X509) PublicKey() ([]byte, uint32) {
	return []byte{}, 0
}

func (x *X509) Read(r *bytes.Buffer) error {
	return fmt.Errorf("X509 read not implemented!")
}

func (x *X509) Verify() bool {
	return false
}
