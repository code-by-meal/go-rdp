package conn

import (
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/x224"
)

type NegoResponse struct {
	Nego
}

func NewNegoResponse() *NegoResponse {
	return &NegoResponse{}
}

func (n *NegoResponse) Read(stream io.Reader) error {
	buff, err := x224.Read(stream)

	if err != nil {
		return fmt.Errorf("nego resp: %v", err)
	}

	log.Dbg("<i>[NEGO-RESPONSE-READ]</> ", buff.Bytes())

	if err := core.Unserialize(buff, n); err != nil {
		return fmt.Errorf("nego resp: unserialize %v", err)
	}

	return nil
}
