package tpkt

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
)

type Header struct {
	Version  uint8
	Reserved uint8
	Length   uint16
}

var (
	HeaderLength = 4
)

func Write(stream io.Writer, buff *bytes.Buffer) error {
	tpktHeader := Header{
		Version:  3,
		Reserved: 0,
		Length:   uint16(buff.Len() + 4),
	}
	tpktPacket, err := core.Serialize(tpktHeader)

	if err != nil {
		return fmt.Errorf("tpkt: serialize bytes %v", err)
	}

	tpktPacket = append(tpktPacket, buff.Bytes()...)

	log.Dbg("<i>[TPKT-HEADER]</> ", tpktHeader)
	log.Dbg("<i>[TPKT-WRITE]</> ", tpktPacket)

	if _, err := stream.Write(tpktPacket); err != nil {
		return fmt.Errorf("tpkt: write buff: %v", err)
	}

	return nil
}

func Read(stream io.Reader) (*bytes.Buffer, error) {
	var tpktHeader Header
	buff := new(bytes.Buffer)
	tpktPacket, err := core.ReadFull(stream, HeaderLength)

	if err != nil {
		return buff, fmt.Errorf("tpkt: read full: %v", err)
	}

	if err := core.Unserialize(bytes.NewBuffer(tpktPacket), &tpktHeader); err != nil {
		return buff, fmt.Errorf("tpkt: unserialize: %v", err)
	}

	if tpktHeader.Version != 3 || tpktHeader.Length <= 4 {
		return buff, fmt.Errorf("tpkt: invalid packet versin: %d length: %d", tpktHeader.Version, tpktHeader.Length)
	}

	tpktData, err := core.ReadFull(stream, int(tpktHeader.Length)-HeaderLength)

	if err != nil {
		return buff, nil
	}

	if _, err := buff.Write(tpktData); err != nil {
		return buff, nil
	}

	log.Dbg("<i>[TPKT-HEADER]</> ", tpktHeader)
	log.Dbg("<i>[TPKT-READ]</> ", buff.Bytes())

	return buff, nil
}
