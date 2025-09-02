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

type NegotiationResult uint8

// Negotiation Result
const (
	ProtocolRDP        NegotiationResult = 0x00000000 //Standard RDP Security
	ProtocolTLS        NegotiationResult = 0x00000001 //TLS1.0/1.1/1.2
	PROTOCOL_HYBRID    NegotiationResult = 0x00000002 //CredSSP
	PROTOCOL_RDSTLS    NegotiationResult = 0x00000004
	PROTOCOL_HYBRID_EX NegotiationResult = 0x00000008
	PROTOCOL_RDSAAD    NegotiationResult = 0x00000010
)

type ConnectionRequest struct {
	Cookie             string `order:"l"`
	Type               uint8  `order:"l"`
	Flags              uint8  `order:"l"`
	Length             uint8  `order:"l"`
	RequestedProtocols uint32 `order:"l"`
}

func NewConnectionRequest(username string) *ConnectionRequest {
	return &ConnectionRequest{
		Cookie:             fmt.Sprintf("Cookie: mstshas=%s", username),
		Type:               uint8(NegotiationRequest),
		Flags:              0,
		Length:             8,
		RequestedProtocols: uint32(ProtocolTLS),
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
