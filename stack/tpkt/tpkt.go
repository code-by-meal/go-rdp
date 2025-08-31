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

func Write(stream io.Writer, buff *bytes.Buffer) error {
	tpktHeader := Header{
		Version:  3,
		Reserved: 0,
		Length:   uint16(buff.Len() + 4),
	}
	tpktPacket, err := core.Serialize(tpktHeader)

	if err != nil {
		return fmt.Errorf("tpkt: write bytes %w", err)
	}

	tpktPacket = append(tpktPacket, buff.Bytes()...)

	log.Dbg("<i>[TPKT-WRITE]</> ", tpktPacket)

	return nil
}

func Read() {

}
