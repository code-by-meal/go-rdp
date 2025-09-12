package client

import (
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/mcs"
	clientdata "github.com/code-by-meal/go-rdp/stack/rdp/client_data"
	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
	serverdata "github.com/code-by-meal/go-rdp/stack/rdp/server_data"
)

func (c *Client) _ChannelConnection() error {
	prefix := "auth: channel-connection: %w"

	// Erect domain request
	log.Zebra("[ERECT-DOMAIN-REQUEST]", log.SuccessColor)

	edr := mcs.NewErrectDomainRequest()

	if err := edr.Write(c.Stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (c *Client) _BasicSettingExchange() error {
	prefix := "basic set exchange: %w"

	// Doing client data request
	log.Zebra("[CLIENT-DATA-REQUEST]", log.SuccessColor)

	cdr := clientdata.NewRequest(c.Hostname, c.SelectedProtocol)

	if err := cdr.Write(c.Stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	// Getting server data response
	log.Zebra("[SERVER-DATA-RESPONSE]", log.SuccessColor)

	sdr := serverdata.NewResponse()

	if err := sdr.Read(c.Stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (c *Client) _Negotiation() error {
	var selectedProtocol nego.NegoProtocol

	requestedProtocols := nego.RDP | nego.Hybrid | nego.TLS // Some commonly used by RDP clients RequestedProtocol flag
	tryCount := 1

	// Try to reconnect if server and client protocols wont mathcing
	for {
		log.Zebra(fmt.Sprintf("[NEGOTIATIOIN-CONNECTION-REQUEST] Try: %d", tryCount), log.SuccessColor)

		negoReq := nego.NewNegoRequest(c.Username, requestedProtocols)

		if err := negoReq.Write(c.Stream); err != nil {
			return fmt.Errorf("nego con-req: %w", err)
		}

		log.Zebra(fmt.Sprintf("[NEGOTIATION-CONFIRM-RESPONSE] Try: %d", tryCount), log.SuccessColor)

		negoRes := nego.NewNegoResponse()

		if err := negoRes.Read(c.Stream); err != nil {
			return fmt.Errorf("nego conf-resp: %w", err)
		}

		switch negoRes.Type { // nolint
		case nego.Failure:
			switch nego.FailureCode(negoRes.RequestedProtocols) {
			// The server requires that the client support Enhanced RDP Security with either TLS 1.0, 1.1 or 1.2 or CredSSP. If only CredSSP was requested then the server only supports TLS.
			case nego.SSLRequiredByServer:
				requestedProtocols |= nego.TLS
				requestedProtocols |= nego.Hybrid
			// The server is configured to only use Standard RDP Security mechanisms and does not support any External Security Protocols.
			case nego.SSLNotAllowedByServer:
				requestedProtocols = nego.RDP
			// The server does not possess a valid authentication certificate and cannot initialize the External Security Protocol Provider.
			case nego.SSLCertNotOnServer:
				return fmt.Errorf("nego: ssl certs not on server")
			// The list of requested security protocols is not consistent with the current security protocol in effect. This error is only possible when the Direct Approach is used and an External Security Protocol is already being used.
			case nego.InconsistentFlags:
				return fmt.Errorf("nego: inconsistent flags")
			// The server requires that the client support Enhanced RDP Security with CredSSP.
			case nego.HybridRequiredByServer:
				requestedProtocols |= nego.Hybrid
			// The server requires that the client support Enhanced RDP Security with TLS 1.0, 1.1 or 1.2  and certificate-based client authentication.
			case nego.SSLWithUserAuthRequiredByServer:
				// Server need authentication from user by trust certificates.
				// Need TLS setting.
				return fmt.Errorf("nego: server need tls + trust certs")
			}

			// You absolutely need to close nego.ction TCP
			if err := c.Stream.Conn.Close(); err != nil {
				return fmt.Errorf("nego: close tcp: %w", err)
			}

			stream, err := core.NewStream(c.Context, c.Host, c.Port, c.Timeout)

			if err != nil {
				return fmt.Errorf("nego: new stream: %w", err)
			}

			c.Stream = stream
		case nego.Response:
			selectedProtocol = negoRes.RequestedProtocols

			goto out_loop
		default:
			log.Info("<e>[UNKNOWN NEGOTIATION RESPONSE TYPE]</> ", fmt.Sprintf("Code: %d", negoReq.Type))
		}

		tryCount++
	}

out_loop:
	log.Dbg(fmt.Sprintf("[<i>SELECTED NEGO PROTOCOL</>: <d>%s</>]", nego.Protocols[selectedProtocol]))

	switch selectedProtocol {
	case nego.RDP:
		// No need to init TLS layer
		// Working without any encryption layer
	case nego.TLS:
		// Need TLS layer
		return fmt.Errorf("not implemented")
	case nego.Hybrid:
		// CredSSP / NLA
		// TLS -> SPNEGO (Kerberos or NTLM)
		return fmt.Errorf("not implemented")
	case nego.RDSTLS:
		// Need TLS layer
		return fmt.Errorf("not implemented")
	case nego.HybridExtended:
		// Need TLS layer
		return fmt.Errorf("not implemented")
	case nego.RDSAAD:
		return fmt.Errorf("not implemented")
	default:
		log.Info(fmt.Sprintf("<e>nego</>: not implemented protocol (code: 0x%X)", selectedProtocol))
	}

	c.SelectedProtocol = selectedProtocol

	return nil
}
