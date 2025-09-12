package certs

import "io"

type X509 struct {
}

func NewX509() *X509 {
	return &X509{}
}

func (x *X509) PublicKey() ([]byte, uint32) {
	return []byte{}, 0
}

func (x *X509) Read(r io.Reader) error {
	return nil
}

func (x *X509) Verify() bool {
	return false
}
