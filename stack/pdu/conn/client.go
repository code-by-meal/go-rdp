package conn

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

func NewNegoRequest(username string, protocol NegoProtocol) *NegoRequest {
	return &NegoRequest{
		Cookie: fmt.Sprintf("Cookie: msthash=%s\r\n", username),
		Nego: Nego{
			Type:               Request,
			Flags:              0,
			Length:             8,
			RequestedProtocols: protocol,
		},
	}
}

func (c *NegoRequest) Write(stream io.Writer) error {
	packet, err := core.Serialize(c)

	if err != nil {
		return fmt.Errorf("nego req: %w", err)
	}

	log.Dbg("<i>[PDU-WRITE]</> ", packet)

	buff := bytes.NewBuffer(packet)

	if err := x224.Write(stream, buff, x224.ConnectionRequestPDU); err != nil {
		return fmt.Errorf("nego req: %w", err)
	}

	return nil
}
