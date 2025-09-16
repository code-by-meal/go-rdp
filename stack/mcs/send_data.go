package mcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/mcs/per"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type SendDataRequest struct {
	InitiatorID  uint16
	ChannelID    uint16
	DataPriority uint8
	Segmentation uint8
	UserData     []byte
}

func NewSendDataRequest(userID uint16) *SendDataRequest {
	return &SendDataRequest{
		InitiatorID:  userID - uint16(UserIDBase),
		ChannelID:    uint16(Global),
		DataPriority: 0x1,
		Segmentation: 0xc0, // Priority + segmentation
	}
}

func (r *SendDataRequest) Write(userData []byte, stream io.Writer) error {
	var buff bytes.Buffer

	prefix := "mcs(1.25): domain request: write: %w"

	if err := per.WriteChoice(&buff, byte(SendDataRequestT)<<2); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.WriteSingleAny(&buff, &r.InitiatorID, binary.BigEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.WriteSingleAny(&buff, &r.ChannelID, binary.BigEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := buff.WriteByte((r.DataPriority & 0x03) << 6); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := per.WriteOctetString(&buff, string(userData), 0); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, &buff, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
