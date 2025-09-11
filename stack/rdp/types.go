package rdp

// Client header types

type ClientHeaderType uint16

const (
	CoreC     ClientHeaderType = 0xC001
	SecurityC ClientHeaderType = 0xC002
	NetC      ClientHeaderType = 0xC003
	ClusterC  ClientHeaderType = 0xC004
	Monitor   ClientHeaderType = 0xC005
)

// Server header types

type ServerHeaderType uint16

const (
	CoreS          ServerHeaderType = 0x0C01
	SecurityS      ServerHeaderType = 0x0C02
	NetS           ServerHeaderType = 0x0C03
	MsgChannel     ServerHeaderType = 0x0C04
	MultyTransport ServerHeaderType = 0x0C08
)
