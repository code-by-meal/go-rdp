package mcs

type ConnectInitial struct {
	CallingDomainSelector uint8 // FIXME: change -> []byte default: 0x01
	CalledDomainSelector  uint8 // FIXME: change -> []byte default: 0x01
	UpwardFlag            uint8
}

type Parameters struct {
	MaxChannelIDs uint8
	MaxUserIDs    uint8
	MaxTokenIDs   uint8
}
