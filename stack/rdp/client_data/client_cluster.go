package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

type ClientClusterData struct {
	HeaderType          rdp.ClientHeaderType `order:"l"`
	HeaderLength        uint16               `order:"l"`
	ClusterFlags        uint32               `order:"l"`
	RedirectedSessionID uint32               `order:"l"`
}

func _NewClientClusterData() *ClientClusterData {
	return &ClientClusterData{
		HeaderType:          rdp.ClusterC,
		HeaderLength:        12,
		ClusterFlags:        0,
		RedirectedSessionID: 0,
	}
}

func (c *ClientClusterData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	ser, err := core.Serialize(c)
	prefix := "rdp: client_cluster: serialize: %w"

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
