package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/mcs/ber"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type ConnectResponse struct {
	Result           uint8
	CalledConnectID  int
	DomainParameters Parameters
	UserData         []byte
}

func NewConnectResponse() *ConnectResponse {
	return &ConnectResponse{}
}

func (c *ConnectResponse) Read(stream io.Reader) (*bytes.Buffer, error) {
	prefix := "mcs-1.25: connect-response: %w"

	buff, err := x224.Read(stream, x224.DataPDU)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	rdrS := bytes.NewReader(buff.Bytes())

	userData, err := ber.ReadApplicationTag(rdrS, ber.Tag(ConnectResponseM))

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	rdrU := bytes.NewReader(userData)

	elem, err := ber.ReadEnumerated(rdrU)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	c.Result = elem

	id, err := ber.ReadInteger(rdrU)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	c.CalledConnectID = id

	p, err := ber.ReadDomainParameters(rdrU)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if err := core.Unserialize(bytes.NewBuffer(p), &(c.DomainParameters)); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	if err := ber.ReadUniversalTag(rdrU, ber.TagOctetString, false); err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	length, err := ber.ReadLength(rdrU)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	data, err := core.ReadFull(rdrU, length)

	if err != nil {
		return buff, fmt.Errorf(prefix, err)
	}

	c.UserData = data

	return bytes.NewBuffer(c.UserData), nil
}
