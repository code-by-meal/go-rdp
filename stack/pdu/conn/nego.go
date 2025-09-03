package conn

type NegoType uint8

const (
	Request  NegoType = 0x01
	Response NegoType = 0x02
	Failure  NegoType = 0x03
)

// Negotiation Support Protocols
type NegoProtocol uint32

const (
	RDP            NegoProtocol = 0x00000000 // Standard RDP Security
	TLS            NegoProtocol = 0x00000001 // TLS1.0/1.1/1.2
	Hybrid         NegoProtocol = 0x00000002 // CredSSP
	RDSTLS         NegoProtocol = 0x00000004 // RDSTLS protocol
	HybridExtended NegoProtocol = 0x00000008 // Credential Security Support Provider protocol coupled with the Early User Authorization Result PDU
	RDSAAD         NegoProtocol = 0x00000010 // RDS-AAD-Auth Security
)

// Negotiation Request Flags
type RequestFlag uint8

const (
	RestrictedAdminModeRequired            RequestFlag = 0x01
	RedirectedAuthentificationModeRequired RequestFlag = 0x02
	CorelationInfoPresent                  RequestFlag = 0x08
)

// Negotiation Response Flags
type ResponseFlag uint8

const (
	ExtendedClientDataSupported             ResponseFlag = 0x01
	DynvcGFXProtocolSupported               ResponseFlag = 0x02
	NegrspFlagReserved                      ResponseFlag = 0x04
	RestrictedAdminModeSupported            ResponseFlag = 0x08
	RedirectedAuthentificationModeSupported ResponseFlag = 0x10
)

// Negotiation Responses
type FailureCode uint32

const (
	SSLRequiredByServer             FailureCode = 0x00000001
	SSLNotAllowedByServer           FailureCode = 0x00000002
	SSLCertNotOnServer              FailureCode = 0x00000003
	InconsistentFlags               FailureCode = 0x00000004
	HybridRequiredByServer          FailureCode = 0x00000005
	SSLWithUserAuthRequiredByServer FailureCode = 0x00000006
)

var (
	Protocols = map[NegoProtocol]string{
		RDP:            "RDP",
		TLS:            "TLS/SSL",
		Hybrid:         "Hybrid",
		RDSTLS:         "RDSTLS",
		HybridExtended: "Hybrid Extended",
		RDSAAD:         "RDSAAD",
	}
)

type Nego struct {
	Type               NegoType     `order:"l"`
	Flags              uint8        `order:"l"`
	Length             uint16       `order:"l"`
	RequestedProtocols NegoProtocol `order:"l"`
}

type NegoRequest struct {
	Cookie string
	Nego
}
