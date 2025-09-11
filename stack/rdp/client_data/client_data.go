package clientdata

import (
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/gcc"
)

type ClientDataRequest struct {
	ClientCoreData
	ClientClusterData
	ClientSecurityData
	ClientNetworkData
	ClientMsgChannelData
	ClientMultiTransportData
}

type ClientMsgChannelData struct{}

type ClientMultiTransportData struct{}

func NewClientDataRequest(hostname string) *ClientDataRequest {
	ccd := _NewClientCoreData(hostname)

	return &ClientDataRequest{
		ClientCoreData:           *ccd,
		ClientClusterData:        ClientClusterData{},
		ClientSecurityData:       ClientSecurityData{},
		ClientNetworkData:        ClientNetworkData{},
		ClientMsgChannelData:     ClientMsgChannelData{},
		ClientMultiTransportData: ClientMultiTransportData{},
	}
}

func (c *ClientDataRequest) Write(stream io.Writer) error {
	prefix := "rdp: client-data: write: %w"
	buffCCD, err := c.ClientCoreData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	ccr := gcc.NewCCR(buffCCD.Bytes())

	if err := ccr.Write(stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
