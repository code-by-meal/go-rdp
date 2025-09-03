package x224

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"

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
	PDUType              TypePDU
	DestinationReference uint16
	SourceReference      uint16
	Flags                uint8
}

var (
	HeaderLength = 7
)

func Write(stream io.Writer, data *bytes.Buffer, pdu TypePDU) error {
	x224Header := Header{
		PDUType:              pdu,
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
	case DataPDU:
		log.Err(fmt.Errorf("x224 data pdu not implemented!"))
	default:
		log.Err(fmt.Sprintf("<e>x224</>: unproccessed pdu type <i>%d</>", pdu))
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

	log.Dbg("<i>[X224-HEADER]</> ", x224Header)
	log.Dbg("<i>[X224-WRITE]</> ", buff.Bytes())

	if err := tpkt.Write(stream, buff); err != nil {
		return fmt.Errorf("x224: %v", err)
	}

	return nil
}

func Read(stream io.Reader) (*bytes.Buffer, error) {
	buff, err := tpkt.Read(stream)

	if err != nil {
		return buff, fmt.Errorf("x224: %v", err)
	}

	var x224Header Header

	if buff.Len() <= HeaderLength {
		return buff, fmt.Errorf("x224: invalid packet length: %d", buff.Len())
	}

	if err := core.Unserialize(buff, &x224Header); err != nil {
		return buff, fmt.Errorf("x224 unserialize: %v", err)
	}

	switch x224Header.PDUType {
	case ConnectionConfirmPDU:
		if x224Header.Length != uint8(buff.Len()+HeaderLength-1) {
			log.Dbg("Header size: ", int(unsafe.Sizeof(x224Header)))
			return buff, fmt.Errorf("x224: invalid header length: %d need: %d", x224Header.Length, HeaderLength)
		}
	case DataPDU:
		if buff.Len() <= 3 || (x224Header.Length != 2 && x224Header.DestinationReference != 0x80) {
			return buff, fmt.Errorf("x224: invalid header")
		}

		log.Err(fmt.Errorf("x224 data pdu not implemented!"))
	default:
	}

	log.Dbg("<i>[X224-HEADER]</> ", x224Header)
	log.Dbg("<i>[X224-READ]</> ", buff.Bytes())

	return buff, nil
}
