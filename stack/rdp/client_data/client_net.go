package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

// Channel Option

type ChannelOption uint32

const (
	Initialized  ChannelOption = 0x80000000
	EncryptRDP   ChannelOption = 0x40000000
	EncryptSC    ChannelOption = 0x20000000
	EncryptCS    ChannelOption = 0x10000000
	PriHight     ChannelOption = 0x08000000
	PriMed       ChannelOption = 0x04000000
	ProLow       ChannelOption = 0x02000000
	CompressRDP  ChannelOption = 0x00800000
	Compress     ChannelOption = 0x00400000
	ShowProtocol ChannelOption = 0x00200000
	Persistence  ChannelOption = 0x00100000
)

type ChannelDef struct {
	Name    [8]byte
	Options ChannelOption `order:"l"`
}

type ClientNetworkData struct {
	HeaderType      rdp.ClientHeaderType `order:"l"`
	HeaderLength    uint16               `order:"l"`
	ChannelCount    uint32               `order:"l"`
	ChannelDefArray []ChannelDef
}

func _NewClientNetworkData() *ClientNetworkData {
	oRdpdr := [8]byte{}
	copy(oRdpdr[:], "rdprd")

	oRdpsnd := [8]byte{}
	copy(oRdpsnd[:], "rdpsnd")

	oCliprdr := [8]byte{}
	copy(oCliprdr[:], "cliprdr")

	oDrvyvnc := [8]byte{}
	copy(oDrvyvnc[:], "drdyvnc")

	return &ClientNetworkData{
		HeaderType:   rdp.NetC,
		HeaderLength: 8 + 12*4,
		ChannelCount: 0x04,
		ChannelDefArray: []ChannelDef{
			ChannelDef{Name: oRdpdr, Options: Initialized | EncryptRDP | CompressRDP},
			ChannelDef{Name: oRdpsnd, Options: Initialized | EncryptRDP},
			ChannelDef{Name: oCliprdr, Options: Initialized | EncryptRDP | CompressRDP | ShowProtocol},
			ChannelDef{Name: oDrvyvnc, Options: Initialized | EncryptRDP | CompressRDP},
		},
	}
}

func (c *ClientNetworkData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	ser, err := core.Serialize(c)
	prefix := "rdp: client-net: serialize: %w"

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	for _, cd := range c.ChannelDefArray {
		ser, err = core.Serialize(cd)

		if err != nil {
			return &buff, fmt.Errorf(prefix, err)
		}

		if _, err := buff.Write(ser); err != nil {
			return &buff, fmt.Errorf(prefix, err)
		}
	}

	return &buff, nil
}
