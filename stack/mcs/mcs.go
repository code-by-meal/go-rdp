package mcs

// MCS Message
const (
	MCSTypeConnectInitial  = 0x65
	MCSTypeConnectResponse = 0x66
)

// MCS PDU type
const (
	MCSPDUTypeErectDomainRequest = 1
	MCSPDUTypeAttachUserRequest  = 10
	MCSPDUTypeAttachUserConfirm  = 11
	MCSPDUTypeChannelJoinRequest = 14
	MCSPDUTypeChannelJoinConfirm = 15
	MCSPDUTypeSendDataRequest    = 25
	MCSPDUTypeSendDataIndication = 26
)

// MCS Channels
const (
	MCSChannelUserIDBase = 1001
	MCSChannelGlobal     = 1003
)
