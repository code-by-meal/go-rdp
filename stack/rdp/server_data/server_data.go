package serverdata

import (
	"errors"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/gcc"
	"github.com/code-by-meal/go-rdp/stack/rdp"
)

type Response struct {
	ServerCoreData                  CoreData
	ServerNetworkData               NetworkData
	ServerSecurityData              SecurityData
	ServerMessageChannelData        MsgChannelData
	ServerMultitransportChannelData MultyTransportData
}

type Header struct {
	Type   rdp.ServerHeaderType `order:"l"`
	Length uint16               `order:"l"`
}

func NewResponse() *Response {
	return &Response{
		ServerCoreData:                  CoreData{},
		ServerNetworkData:               NetworkData{},
		ServerSecurityData:              SecurityData{},
		ServerMessageChannelData:        MsgChannelData{},
		ServerMultitransportChannelData: MultyTransportData{},
	}
}

func (r *Response) Read(stream io.Reader) error {
	prefix := "rdp: server-data: read: %w"
	ccrsp := gcc.NewConfernceCreateResponse()

	buff, err := ccrsp.Read(stream)

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	for {
		header := Header{}

		if err := core.Unserialize(buff, &header); err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf(prefix, err)
		}

		if header.Length <= 0 {
			break
		}

		switch header.Type {
		case rdp.CoreS:
			if err := core.Unserialize(buff, &r.ServerCoreData); err != nil {
				return fmt.Errorf(prefix, err)
			}
		case rdp.NetS:
			if err := core.Unserialize(buff, &r.ServerNetworkData); err != nil {
				return fmt.Errorf(prefix, err)
			}

			for i := 0; i < int(r.ServerNetworkData.ChannelCount); i++ {
				var tmp struct {
					ChannelID uint16 `order:"l"`
				}

				if err := core.Unserialize(buff, &tmp); err != nil {
					return fmt.Errorf(prefix, err)
				}

				r.ServerNetworkData.ChannelIDArray = append(r.ServerNetworkData.ChannelIDArray, tmp.ChannelID)
			}
		case rdp.SecurityS:
			if err := core.Unserialize(buff, &r.ServerSecurityData); err != nil {
				return fmt.Errorf(prefix, err)
			}

			if r.ServerSecurityData.EncryptionLevel == 0 || r.ServerSecurityData.EncryptionMethod == 0 {
				continue
			}

			random, err := core.ReadFull(buff, int(r.ServerSecurityData.ServerRandomLen))

			if err != nil {
				return fmt.Errorf(prefix, err)
			}

			r.ServerSecurityData.ServerRandom = random

			cert, err := core.ReadFull(buff, int(r.ServerSecurityData.ServerCertLen))

			if err != nil {
				return fmt.Errorf(prefix, err)
			}

			r.ServerSecurityData.ServerCertificate = cert

		case rdp.MsgChannelS:
			if err := core.Unserialize(buff, &r.ServerMessageChannelData); err != nil {
				return fmt.Errorf(prefix, err)
			}

		case rdp.MultyTransportS:
			if err := core.Unserialize(buff, &r.ServerMultitransportChannelData); err != nil {
				return fmt.Errorf(prefix, err)
			}

		default:
			return fmt.Errorf(prefix, fmt.Errorf("unexpected server-data header-type: 0x%X", header.Type))
		}
	}

	return nil
}
