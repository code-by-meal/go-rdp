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
	PDUType              TypePDU
	DestinationReference uint16
	SourceReference      uint16
	Flags                uint8
}
type DataHeader struct {
	Length               uint8
	PDUType              TypePDU
	DestinationReference uint8
}

var (
	HeaderLength = 7
)

func _WriteHeader(pdu TypePDU, length int) ([]byte, error) {
	prefix := "x224: get header: %w"

	switch pdu {
	case DataPDU:
		header := DataHeader{
			Length:               2,
			PDUType:              DataPDU,
			DestinationReference: 0x80,
		}

		headerBody, err := core.Serialize(&header)

		if err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}

		return headerBody, nil
	case ConnectionRequestPDU:
		if length > 0xf9 {
			return []byte{}, fmt.Errorf(prefix, fmt.Errorf("invalid data length: %d, cant more then 0xf9", length))
		}

		header := Header{
			PDUType:              ConnectionRequestPDU,
			DestinationReference: 0,
			SourceReference:      0,
			Length:               uint8(length + 6),
		}

		headerBody, err := core.Serialize(&header)

		if err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}

		return headerBody, nil

	default:
		log.Info(fmt.Sprintf("<e>x224 unknown header: </> <d>%s</>", string(pdu)))
	}

	return []byte{}, fmt.Errorf("x224: get header: invalid pdu type")
}

func Write(stream io.Writer, data *bytes.Buffer, pdu TypePDU) error {
	buff := new(bytes.Buffer)
	x224Packet, err := _WriteHeader(pdu, data.Len())

	if err != nil {
		return fmt.Errorf("x224: %w", err)
	}

	if _, err := buff.Write(x224Packet); err != nil {
		return fmt.Errorf("x224: %w", err)
	}

	if _, err := buff.Write(data.Bytes()); err != nil {
		return fmt.Errorf("x224: %w", err)
	}

	log.Dbg("<i>[X224-WRITE]</> ", buff.Bytes())

	if err := tpkt.Write(stream, buff); err != nil {
		return fmt.Errorf("x224: %w", err)
	}

	return nil
}

func Read(stream io.Reader, pdu TypePDU) (*bytes.Buffer, error) {
	buff, err := tpkt.Read(stream)

	if err != nil {
		return buff, fmt.Errorf("x224: %w", err)
	}

	if buff.Len() <= HeaderLength {
		return buff, fmt.Errorf("x224: invalid packet length: %d", buff.Len())
	}

	switch pdu { // nolint
	case ConnectionConfirmPDU:
		var x224Header Header

		if err := core.Unserialize(buff, &x224Header); err != nil {
			return buff, fmt.Errorf("x224 unserialize: %w", err)
		}

		if x224Header.Length != uint8(buff.Len()+HeaderLength-1) {
			return buff, fmt.Errorf("x224: invalid header length: %d need: %d", x224Header.Length, HeaderLength)
		}
	case DataPDU:
		var x224Header DataHeader

		if err := core.Unserialize(buff, &x224Header); err != nil {
			return buff, fmt.Errorf("x224 unserialize: %w", err)
		}

		if buff.Len() <= 3 || (x224Header.Length != 2 && x224Header.DestinationReference != 0x80) {
			return buff, fmt.Errorf("x224: invalid header")
		}
	default:
	}

	log.Dbg("<i>[X224-READ]</> ", buff.Bytes())

	return buff, nil
}
