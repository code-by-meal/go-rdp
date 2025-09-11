package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

type ClientMultiTransportData struct {
	HeaderType          rdp.ClientHeaderType `order:"l"`
	HeaderLength        uint16               `order:"l"`
	MultiTransportFlags uint32               `order:"l"`
}

func _NewClientMultiTransport() *ClientMultiTransportData {
	return &ClientMultiTransportData{
		HeaderType:          rdp.MultiTransportData,
		HeaderLength:        8,
		MultiTransportFlags: 1,
	}
}

func (c *ClientMultiTransportData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	prefix := "rdp: client-multi-transport: serialize: %w"
	ser, err := core.Serialize(c)

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
