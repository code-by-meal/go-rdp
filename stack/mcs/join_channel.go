package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/mcs/per"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

// Request
type JoinChannelRequest struct {
	InitiatorID uint16
	ChannelID   uint16
}

func NewJoinChannelRequest(initiator, channel uint16) *JoinChannelRequest {
	return &JoinChannelRequest{
		InitiatorID: initiator - uint16(UserIDBase),
		ChannelID:   channel,
	}
}

func (j *JoinChannelRequest) Write(stream io.Writer) error {
	var buff bytes.Buffer

	prefix := "join-channel: write: %w"

	if err := per.WriteChoice(&buff, byte(ChannelJoinRequestT)<<2); err != nil {
		return fmt.Errorf(prefix, err)
	}

	ser, err := core.Serialize(j)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, &buff, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

// Response
type JoinChannelConfirm struct {
	Confirm   uint8
	UserID    uint16
	ChannelID uint16
}

func NewJoinChannelConfirm() *JoinChannelConfirm {
	return &JoinChannelConfirm{}
}

func (j *JoinChannelConfirm) Read(stream io.Reader) error {
	prefix := "join-channel: read: %w"

	buff, err := x224.Read(stream, x224.DataPDU)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	c, err := per.ReadChoice(buff)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if (c >> 2) != byte(ChannelJoinConfirmT) {
		return fmt.Errorf(prefix, fmt.Errorf("read invalid header: %d", c>>2))
	}

	enum, err := per.ReadEnumerated(buff)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	j.Confirm = enum

	ui, err := per.ReadInteger16(buff, 0)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	j.UserID = ui + uint16(UserIDBase)

	cid, err := per.ReadInteger16(buff, 0)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	j.ChannelID = cid

	if j.Confirm != 0 {
		return fmt.Errorf(prefix, fmt.Errorf("not confirmed confirm from server"))
	}

	return nil
}
