package client

import (
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/pdu/conn"
)

func (c *Client) _Negotiation() error {
	var selectedProtocol conn.NegoProtocol
	requestedProtocols := conn.RDP | conn.Hybrid | conn.TLS // Some commonly used by RDP clients RequestedProtocol flag
	tryCount := 1

	// Try to reconnect if server and client protocols wont mathcing
	for {
		log.Zebra(fmt.Sprintf("\n[NEGOTIATIOIN-CONNECTION-REQUEST] Try: %d", tryCount), log.SuccessColor)

		negoReq := conn.NewNegoRequest(c.Username, requestedProtocols)

		if err := negoReq.Write(c.Stream); err != nil {
			return fmt.Errorf("nego con-req: %v", err)
		}

		log.Zebra(fmt.Sprintf("\n[NEGOTIATION-CONFIRM-RESPONSE] Try: %d", tryCount), log.SuccessColor)

		negoRes := conn.NewNegoResponse()

		if err := negoRes.Read(c.Stream); err != nil {
			return fmt.Errorf("nego conf-resp: %v", err)
		}

		switch negoRes.Type {
		case conn.Failure:
			switch conn.FailureCode(negoRes.RequestedProtocols) {
			// The server requires that the client support Enhanced RDP Security with either TLS 1.0, 1.1 or 1.2 or CredSSP. If only CredSSP was requested then the server only supports TLS.
			case conn.SSLRequiredByServer:
				requestedProtocols |= conn.TLS
				requestedProtocols |= conn.Hybrid
			// The server is configured to only use Standard RDP Security mechanisms and does not support any External Security Protocols.
			case conn.SSLNotAllowedByServer:
				requestedProtocols = conn.RDP
			// The server does not possess a valid authentication certificate and cannot initialize the External Security Protocol Provider.
			case conn.SSLCertNotOnServer:
				return fmt.Errorf("nego: ssl certs not on server")
			// The list of requested security protocols is not consistent with the current security protocol in effect. This error is only possible when the Direct Approach is used and an External Security Protocol is already being used.
			case conn.InconsistentFlags:
				return fmt.Errorf("nego: inconsistent flags")
			// The server requires that the client support Enhanced RDP Security with CredSSP.
			case conn.HybridRequiredByServer:
				requestedProtocols |= conn.Hybrid
			// The server requires that the client support Enhanced RDP Security with TLS 1.0, 1.1 or 1.2  and certificate-based client authentication.
			case conn.SSLWithUserAuthRequiredByServer:
				// Server need authentication from user by trust certificates.
				// Need TLS setting.
				return fmt.Errorf("nego: server need tls + trust certs")
			}

			// You absolutely need to close connection TCP
			if err := c.Stream.Conn.Close(); err != nil {
				return fmt.Errorf("nego: close tcp: %v", err)
			}

			stream, err := core.NewStream(c.Host, c.Port, c.Timeout, c.Context)

			if err != nil {
				return fmt.Errorf("nego: new stream: %v", err)
			}

			c.Stream = stream
		case conn.Response:
			selectedProtocol = negoRes.RequestedProtocols

			goto out_loop
		default:
			log.Info("<e>[UNKNOWN NEGOTIATION RESPONSE TYPE]</> ", fmt.Sprintf("Code: %d", negoReq.Type))
		}

		tryCount += 1
	}

out_loop:
	switch selectedProtocol {
	case conn.RDP:
	case conn.TLS:
	case conn.Hybrid:
	case conn.RDSTLS:
	case conn.HybridExtended:
	case conn.RDSAAD:
	default:
		log.Info(fmt.Sprintf("<e>nego</>: not implemented protocol (code: 0x%X)", selectedProtocol))
	}

	log.Dbg(fmt.Sprintf("[<i>SELECTED NEGO PROTOCOL</>: <d>%s</>]", conn.Protocols[selectedProtocol]))

	return nil
}
