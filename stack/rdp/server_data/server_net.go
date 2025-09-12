package serverdata

type NetworkData struct {
	MCSChannelID   uint16 `order:"l"`
	ChannelCount   uint16 `order:"l"`
	ChannelIDArray []uint16
}
