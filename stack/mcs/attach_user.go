package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/mcs/per"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

// Request

type AttachUserRequest struct{}

func NewAttachUserRequest() *AttachUserRequest {
	return &AttachUserRequest{}
}

func (a *AttachUserRequest) Write(stream io.Writer) error {
	var buff bytes.Buffer

	prefix := "mcs: attach-user: write: %w"

	if err := per.WriteChoice(&buff, byte(AttachUserRequestT)<<2); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, &buff, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

// Response

type AttachUserConfirm struct {
	UserID uint16
}

func NewAttachUserConfirm() *AttachUserConfirm {
	return &AttachUserConfirm{
		UserID: 0,
	}
}

func (a *AttachUserConfirm) Read(stream io.Reader) error {
	prefix := "mcs: attach-user-confirm: read: %w"

	buff, err := x224.Read(stream, x224.DataPDU)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	typ, err := per.ReadChoice(buff)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if AttachUserConfirmT != PDUType(typ>>2) {
		return fmt.Errorf(prefix, fmt.Errorf("invalid type of UAC: %d", typ>>2))
	}

	enum, err := per.ReadEnumerated(buff)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	if enum != 0 {
		return fmt.Errorf(prefix, fmt.Errorf("invalid enumarated"))
	}

	userID, err := per.ReadInteger16(buff, 0)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	a.UserID = userID + uint16(UserIDBase)

	log.Dbg(fmt.Sprintf("[<i>USER-ID</>: <d>%d</>]", a.UserID))

	return nil
}
