package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/mcs/ber"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type ConnectInitial struct {
	CallingDomainSelector []byte // FIXME: change -> []byte default: 0x01
	CalledDomainSelector  []byte // FIXME: change -> []byte default: 0x01
	UpwardFlag            bool
	TargetParameters      Parameters
	MinimumParameters     Parameters
	MaximumParameters     Parameters
	UserData              []byte
}

type Parameters struct {
	MaxChannelIDs   int
	MaxUserIDs      int
	MaxTokenIDs     int
	NumPriorities   int
	MinThoughput    int
	MaxHeight       int
	MaxMCSPDUsize   int
	ProtocolVersion int
}

func NewConnectInitial(userData []byte) *ConnectInitial {
	return &ConnectInitial{
		CallingDomainSelector: []byte{0x1},
		CalledDomainSelector:  []byte{0x1},
		UpwardFlag:            true,
		TargetParameters:      Parameters{34, 2, 0, 1, 0, 1, 0xffff, 2},
		MinimumParameters:     Parameters{1, 1, 1, 1, 0, 1, 0x420, 2},
		MaximumParameters:     Parameters{0xffff, 0xfc17, 0xffff, 1, 0, 1, 0xffff, 2},
		UserData:              userData,
	}
}

// Connect intial
func (ci *ConnectInitial) Write(stream io.Writer) error {
	prefix := "t125(mcs): write: %w"
	buff, err := ci.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buff1 := new(bytes.Buffer)

	if err := ber.WriteApplicationTag(buff1, ber.Tag(ConnectInitialM), buff.Bytes()); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := x224.Write(stream, buff1, x224.DataPDU); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (ci *ConnectInitial) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	prefix := "t125: conn-init: %w"

	if err := ber.WriteOctetString(&buff, string(ci.CallingDomainSelector)); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteOctetString(&buff, string(ci.CalledDomainSelector)); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteBool(&buff, ci.UpwardFlag); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	params, err := ci.TargetParameters.Serialize()

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteParameters(&buff, params.Bytes()); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	params, err = ci.MinimumParameters.Serialize()

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteParameters(&buff, params.Bytes()); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	params, err = ci.MaximumParameters.Serialize()

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteParameters(&buff, params.Bytes()); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteOctetString(&buff, string(ci.UserData)); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}

// Parameters
func (p *Parameters) Serialize() (*bytes.Buffer, error) {
	prefix := "t125(mcs): serialize: %w"

	var buff bytes.Buffer

	if err := ber.WriteInteger(&buff, p.MaxChannelIDs); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.MaxUserIDs); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.MaxTokenIDs); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.NumPriorities); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.MinThoughput); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.MaxHeight); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.MaxMCSPDUsize); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if err := ber.WriteInteger(&buff, p.ProtocolVersion); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
