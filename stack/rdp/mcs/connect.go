package mcs

type ClientDataRequest struct {
	ClientCoreData
	ClientClusterData
	ClientSecurityData
	ClientNetworkData
	ClientMsgChannelData
	ClientMultiTransportData
}

type ClientCoreData struct {
	HeaderType             uint16
	HeaderLength           uint16
	Version                uint32
	DesktopWidth           uint16
	DesktopHeight          uint16
	ColorDepth             uint16
	SasSequence            uint16
	KbdLayout              uint32
	ClientBuild            uint32
	ClientName             [32]byte
	KeyboardType           uint32
	KeyboardSubType        uint32
	KeyboardFnKeys         uint32
	ImeFileName            [64]byte
	PostBeta2ColorDepth    uint16
	ClientProductId        uint16
	SerialNumber           uint32
	HighColorDepth         uint16
	SupportedColorDepths   uint16
	EarlyCapabilityFlags   uint16
	ClientDigProductId     [64]byte
	ConnectionType         uint8
	Pad1octet              uint8
	ServerSelectedProtocol uint32
}

type ClientClusterData struct{}

type ClientSecurityData struct{}

type ClientNetworkData struct{}

type ClientMsgChannelData struct{}

type ClientMultiTransportData struct{}

func NewClientDataRequest() *ClientDataRequest {
	return &ClientDataRequest{
		ClientCoreData:           ClientCoreData{},
		ClientClusterData:        ClientClusterData{},
		ClientSecurityData:       ClientSecurityData{},
		ClientNetworkData:        ClientNetworkData{},
		ClientMsgChannelData:     ClientMsgChannelData{},
		ClientMultiTransportData: ClientMultiTransportData{},
	}
}
