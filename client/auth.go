package client

import (
	"fmt"

	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/pdu/conn"
)

func (c *Client) _Negotiation() error {
	log.Zebra("\n[NEGOTIATIOIN-CONNECTION-REQUEST]", log.SuccessColor)

	negoReq := conn.NewNegoRequest(c.Username)

	if err := negoReq.Write(c.Stream); err != nil {
		return fmt.Errorf("nego con-req: %v", err)
	}

	log.Zebra("\n[NEGOTIATION-CONFIRM-RESPONSE]", log.SuccessColor)

	negoRes := conn.NewNegoResponse()

	if err := negoRes.Read(c.Stream); err != nil {
		return fmt.Errorf("nego conf-resp: %v", err)
	}

	switch negoRes.Nego.RequestedProtocols {
	case conn.ProtocolRDP:

	case conn.ProtocolTLS:
	case conn.ProtocolHybrid:
	case conn.ProtocolRDSTLS:
	case conn.ProtocolHybridExtended:
	case conn.ProtocolRDSAAD:
	default:
		log.Info(fmt.Sprintf("<e>nego</>: not implemented protocol (code: 0x%X)", negoRes.Nego.RequestedProtocols))
	}

	log.Dbg(fmt.Sprintf("[<i>SELECTED NEGO PROTOCOL</>: <d>%s</>]", conn.Protocols[negoRes.Nego.RequestedProtocols]))

	return nil
}
