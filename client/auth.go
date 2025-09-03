package client

import (
	"fmt"

	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/pdu/conn"
)

func (c *Client) _Negotiation() error {
	log.Zebra("\n\n[NEGOTIATIOIN-CONNECTION-REQUEST]", log.SuccessColor)

	negoReq := conn.NewNegoRequest(c.Username)

	if err := negoReq.Write(c.Stream); err != nil {
		return fmt.Errorf("nego con-req: %v", err)
	}

	log.Zebra("\n\n[NEGOTIATION-CONFIRM-RESPONSE]", log.SuccessColor)

	negoRes := conn.NewNegoResponse()

	if err := negoRes.Read(); err != nil {
		return fmt.Errorf("nego conf-resp: %v", err)
	}

	return nil
}
