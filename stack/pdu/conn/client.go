package conn

import (
	"bytes"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type NegotiationType uint8

const (
	NegotiationRequest  NegotiationType = 0x01
	NegotiationResponse NegotiationType = 0x02
	NegotiationFailure  NegotiationType = 0x03
)

type NegotiationResult uint32

// Negotiation Result
const (
	ProtocolRDP            NegotiationResult = 0x00000000 //Standard RDP Security
	ProtocolTLS            NegotiationResult = 0x00000001 //TLS1.0/1.1/1.2
	ProtocolHybrid         NegotiationResult = 0x00000002 //CredSSP
	ProtocolRDSTLS         NegotiationResult = 0x00000004
	ProtocolHybridExtended NegotiationResult = 0x00000008
	ProtocolRDSAAD         NegotiationResult = 0x00000010
)

type ConnectionRequest struct {
	Cookie             string `order:"l"`
	Type               uint8  `order:"l"`
	Flags              uint8  `order:"l"`
	Length             uint16 `order:"l"`
	RequestedProtocols uint32 `order:"l"`
}

func NewConnectionRequest(username string) *ConnectionRequest {
	return &ConnectionRequest{
		Cookie:             fmt.Sprintf("Cookie: msthash=%s", username),
		Type:               uint8(NegotiationRequest),
		Flags:              0,
		Length:             8,
		RequestedProtocols: uint32(ProtocolTLS) | uint32(ProtocolHybrid) | uint32(ProtocolRDP),
	}
}

func (c *ConnectionRequest) Write(stream io.Writer) error {
	packet, err := core.Serialize(c)

	if err != nil {
		return fmt.Errorf("pdu connection-request: %v", err)
	}

	log.Dbg("<i>[PDU-WRITE]</> ", packet)

	buff := bytes.NewBuffer(packet)

	if err := x224.Write(stream, buff, x224.ConnectionRequestPDU); err != nil {
		return fmt.Errorf("pdu connection-request: %v", err)
	}

	return nil
}
