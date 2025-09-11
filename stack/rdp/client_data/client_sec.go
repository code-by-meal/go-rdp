package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

// Encryption Methods

type EncMethod uint32

const (
	Four0BIT EncMethod = 0x00000001
	One28BIT EncMethod = 0x00000002
	Five6BIT EncMethod = 0x00000008
	Fips     EncMethod = 0x00000010
)

type ClientSecurityData struct {
	HeaderType           rdp.ClientHeaderType `order:"l"`
	HeaderLength         uint16               `order:"l"`
	EncryptionMethods    EncMethod            `order:"l"`
	ExtEncryptionMethods uint32               `order:"l"`
}

func _NewClientSecurityData() *ClientSecurityData {
	return &ClientSecurityData{
		HeaderType:           rdp.SecurityC,
		HeaderLength:         12,
		EncryptionMethods:    Four0BIT | One28BIT | Five6BIT,
		ExtEncryptionMethods: 0,
	}
}

func (c *ClientSecurityData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	ser, err := core.Serialize(c)
	prefix := "rdp: client-sec: serialize: %w"

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
