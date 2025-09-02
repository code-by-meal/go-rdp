package x224

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/tpkt"
)

type TypePDU uint8

const (
	ConnectionRequestPDU    TypePDU = 0xe0
	ConnectionConfirmPDU    TypePDU = 0xd0
	DisconnectionRequestPDU TypePDU = 0x80
	DataPDU                 TypePDU = 0xf0
	ErrorPDU                TypePDU = 0x70
)

type Header struct {
	Length               uint8
	PDUType              uint8
	DestinationReference uint16
	SourceReference      uint16
	Flags                uint8
}

func Write(stream io.Writer, data *bytes.Buffer, pdu TypePDU) error {
	x224Header := Header{
		PDUType:              uint8(pdu),
		DestinationReference: 0,
		SourceReference:      0,
		Flags:                0,
	}
	buff := new(bytes.Buffer)

	switch pdu {
	case ConnectionRequestPDU:
		if data.Len() > 0xf9 {
			return fmt.Errorf("x224: invalid data length: %d, can't be more than 0xf9", data.Len())
		}
		x224Header.Length = uint8(data.Len() + 6)
	case ConnectionConfirmPDU:
	case DataPDU:
	default:
		log.Err(fmt.Sprintf("<e>x224</>: unknown pdu type <i>%d</>", pdu))
	}

	x224Packet, err := core.Serialize(x224Header)

	if err != nil {
		return fmt.Errorf("x224: %v", err)
	}

	if _, err := buff.Write(x224Packet); err != nil {
		return fmt.Errorf("x224: %v", err)
	}

	if _, err := buff.Write(data.Bytes()); err != nil {
		return fmt.Errorf("x224: %v", err)
	}

	log.Dbg("<i>[X224-WRITE]</> ", buff.Bytes())

	if err := tpkt.Write(stream, buff); err != nil {
		return fmt.Errorf("x224: %v", err)
	}

	return nil
}

func Read() {
	// implement
}
