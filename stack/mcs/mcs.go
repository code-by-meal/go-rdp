package mcs

// MCS Message

type Message uint8

const (
	ConnectInitialM  Message = 0x65
	ConnectResponseM Message = 0x66
)

// MCS PDU type

type PDUType uint16

const (
	ErectDomainRequestT PDUType = 1
	AttachUserRequestT  PDUType = 10
	AttachUserConfirmT  PDUType = 11
	ChannelJoinRequestT PDUType = 14
	ChannelJoinConfirmT PDUType = 15
	SendDataRequestT    PDUType = 25
	SendDataIndicationT PDUType = 26
)

// MCS Channels

type Channel uint16

const (
	UserIDBase Channel = 1001
	Global     Channel = 1003
)
