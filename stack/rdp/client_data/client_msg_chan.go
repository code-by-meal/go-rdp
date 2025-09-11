package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

type ClientMsgChannelData struct {
	HeaderType      rdp.ClientHeaderType `order:"l"`
	HeaderLength    uint16               `order:"l"`
	MsgChannelFlags uint32               `order:"l"`
}

func _NewClientMsgChannelData() *ClientMsgChannelData {
	return &ClientMsgChannelData{
		HeaderType:      rdp.MsgChannelData,
		HeaderLength:    8,
		MsgChannelFlags: 0x00000000,
	}
}

func (c *ClientMsgChannelData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	prefix := "rdp: client-msg-channel: serialize: %w"
	ser, err := core.Serialize(c)

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
