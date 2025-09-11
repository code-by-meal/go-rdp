package clientdata

import (
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/gcc"
	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
)

type Request struct {
	ClientCoreData
	ClientClusterData
	ClientSecurityData
	ClientNetworkData
	ClientMsgChannelData
	ClientMultiTransportData
}

func NewRequest(hostname string, protocol nego.NegoProtocol) *Request {
	ccd := _NewClientCoreData(hostname)
	ccld := _NewClientClusterData()
	csd := _NewClientSecurityData()
	cnd := _NewClientNetworkData()
	cmcd := _NewClientMsgChannelData()
	cmtd := _NewClientMultiTransport()

	ccd.ServerSelectedProtocol = protocol

	return &Request{
		ClientCoreData:           *ccd,
		ClientClusterData:        *ccld,
		ClientSecurityData:       *csd,
		ClientNetworkData:        *cnd,
		ClientMsgChannelData:     *cmcd,
		ClientMultiTransportData: *cmtd,
	}
}

func (c *Request) Write(stream io.Writer) error {
	prefix := "rdp: client-data: write: %w"
	buffCCD, err := c.ClientCoreData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buffCCLD, err := c.ClientClusterData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buffCSD, err := c.ClientSecurityData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buffCND, err := c.ClientNetworkData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buffCMCD, err := c.ClientMsgChannelData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buffCMTD, err := c.ClientMultiTransportData.Serialize()

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	buff := []byte{}
	buff = append(buff, buffCCD.Bytes()...)
	buff = append(buff, buffCCLD.Bytes()...)
	buff = append(buff, buffCSD.Bytes()...)
	buff = append(buff, buffCND.Bytes()...)
	buff = append(buff, buffCMCD.Bytes()...)
	buff = append(buff, buffCMTD.Bytes()...)

	ccr := gcc.NewCCR(buff)

	if err := ccr.Write(stream); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}
