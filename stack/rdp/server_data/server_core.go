package serverdata

type CoreData struct {
	Version                  uint32 `order:"l"`
	ClientRequestedProtocols uint32 `order:"l"`
	EarlyCapabilityFlags     uint32 `order:"l"`
}
