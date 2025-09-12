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
	ErectDomainRequest PDUType = 1
	AttachUserRequest  PDUType = 10
	AttachUserConfirm  PDUType = 11
	ChannelJoinRequest PDUType = 14
	ChannelJoinConfirm PDUType = 15
	SendDataRequest    PDUType = 25
	SendDataIndication PDUType = 26
)

// MCS Channels

type Channel uint16

const (
	UserIDBase Channel = 1001
	Global     Channel = 1003
)
