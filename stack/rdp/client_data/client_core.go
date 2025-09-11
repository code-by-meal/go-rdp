package clientdata

import (
	"bytes"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/stack/rdp"
	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
)

// RDP Versions

type Version uint32

const (
	V4     Version = 0x00080001
	V5Plus Version = 0x00080004 // RDP 5.0, 5.1, 5.2, 6.0, 6.1, 7.0, 7.1, 8.0, and 8.1 clients
	V10    Version = 0x00080005 // RDP 10.0 clients
	V101   Version = 0x00080006
	V102   Version = 0x00080007
	V103   Version = 0x00080008
	V104   Version = 0x00080009
	V105   Version = 0x0008000a
	V106   Version = 0x0008000b
	V107   Version = 0x0008000C
	V108   Version = 0x0008000D
	V109   Version = 0x0008000E
	V1010  Version = 0x0008000F
	V1011  Version = 0x00080010
)

// RDP Color depths

type ColorDepth uint16

const (
	EightBPP   ColorDepth = 0xCA01
	One6BPP555 ColorDepth = 0xCA02
	One6BPP565 ColorDepth = 0xCA03
	Two4BPP    ColorDepth = 0xCA04
)

// RDP SASS Sequence

type Sequence uint16

const (
	SasDel Sequence = 0xAA03
)

// RDP Keyboard layout

type KeyboardLayout uint32

const (
	Arabic    KeyboardLayout = 0x00000401
	Bulgarian KeyboardLayout = 0x00000402
	ChineseUS KeyboardLayout = 0x00000404
	Czech     KeyboardLayout = 0x00000405
	Danish    KeyboardLayout = 0x00000406
	German    KeyboardLayout = 0x00000407
	Greek     KeyboardLayout = 0x00000408
	US        KeyboardLayout = 0x00000409
	Spanish   KeyboardLayout = 0x0000040a
	Finnish   KeyboardLayout = 0x0000040b
	French    KeyboardLayout = 0x0000040c
	Hebrew    KeyboardLayout = 0x0000040d
	Hungurian KeyboardLayout = 0x0000040e
	Icelandic KeyboardLayout = 0x0000040f
	Italian   KeyboardLayout = 0x00000410
	Japanese  KeyboardLayout = 0x00000411
	Korean    KeyboardLayout = 0x00000412
	Dutch     KeyboardLayout = 0x00000413
	Norwegian KeyboardLayout = 0x00000414
)

// RDP Keyboard type

type KeyboardType uint32

const (
	IbmPcXT83Key  KeyboardType = 0x00000001
	Olivetti      KeyboardType = 0x00000002
	IbmPcAT84Keys KeyboardType = 0x00000003
	Ibm101102Keys KeyboardType = 0x00000004
	Nokia1050     KeyboardType = 0x00000005
	Nokia9140     KeyboardType = 0x00000006
	JapaneseKey   KeyboardType = 0x00000007
)

// RDP Hight color depth

type HightColorDepth uint16

const (
	FourBPP    HightColorDepth = 0x0004
	EightBPPHC HightColorDepth = 0x0008
	One5BPP    HightColorDepth = 0x000f
	One6BPP    HightColorDepth = 0x0010
	Two4BPPHC  HightColorDepth = 0x0018
)

// RDP Bit field

type SupportColorDepth uint16

const (
	Two4BPPSCD SupportColorDepth = 0x0001
	One6BPPSCD SupportColorDepth = 0x0002
	One5BPPSCD SupportColorDepth = 0x0004
	Three2BPP  SupportColorDepth = 0x0008
)

// RDP Capability flags

type CapabiltyFlag uint16

const (
	SupportErrInfoPdu       CapabiltyFlag = 0x0001
	Want32BPPSession        CapabiltyFlag = 0x0002
	SupportStatusInfoPdu    CapabiltyFlag = 0x0004
	StrongAsymetricKeys     CapabiltyFlag = 0x0008
	Unused                  CapabiltyFlag = 0x0010
	ValidConnectionType     CapabiltyFlag = 0x0020
	SupportMonitorLayoutPdu CapabiltyFlag = 0x0040
	NetcharAutodetect       CapabiltyFlag = 0x0080
	DyVNCGFXProtocol        CapabiltyFlag = 0x0100
	DynamicTimeZone         CapabiltyFlag = 0x0200
	HeartBeatPdu            CapabiltyFlag = 0x0400
)

type ClientCoreData struct {
	HeaderType             rdp.ClientHeaderType `order:"l"`
	HeaderLength           uint16               `order:"l"`
	Version                Version              `order:"l"`
	DesktopWidth           uint16               `order:"l"`
	DesktopHeight          uint16               `order:"l"`
	ColorDepth             ColorDepth           `order:"l"`
	SasSequence            Sequence             `order:"l"`
	KbdLayout              KeyboardLayout       `order:"l"`
	ClientBuild            uint32               `order:"l"`
	ClientName             [32]byte
	KeyboardType           KeyboardType `order:"l"`
	KeyboardSubType        uint32       `order:"l"`
	KeyboardFnKeys         uint32       `order:"l"`
	ImeFileName            [64]byte
	PostBeta2ColorDepth    ColorDepth        `order:"l"`
	ClientProductID        uint16            `order:"l"`
	SerialNumber           uint32            `order:"l"`
	HighColorDepth         HightColorDepth   `order:"l"`
	SupportedColorDepths   SupportColorDepth `order:"l"`
	EarlyCapabilityFlags   CapabiltyFlag     `order:"l"`
	ClientDigProductID     [64]byte
	ConnectionType         uint8
	Pad1octet              uint8
	ServerSelectedProtocol nego.NegoProtocol `order:"l"`
}

func _NewClientCoreData(hostname string) *ClientCoreData {
	ccd := ClientCoreData{
		HeaderType:             rdp.CoreC,
		HeaderLength:           0xd8,
		Version:                V5Plus,
		DesktopWidth:           1280,
		DesktopHeight:          800,
		ColorDepth:             EightBPP,
		SasSequence:            SasDel,
		KbdLayout:              US,
		ClientBuild:            3790,
		ClientName:             [32]byte{},
		KeyboardType:           Ibm101102Keys,
		KeyboardSubType:        0,
		KeyboardFnKeys:         12,
		ImeFileName:            [64]byte{},
		PostBeta2ColorDepth:    EightBPP,
		SerialNumber:           0,
		HighColorDepth:         Two4BPPHC,
		SupportedColorDepths:   One5BPPSCD | One6BPPSCD | Two4BPPSCD | Three2BPP,
		EarlyCapabilityFlags:   SupportErrInfoPdu,
		ClientDigProductID:     [64]byte{},
		ConnectionType:         0,
		Pad1octet:              0,
		ServerSelectedProtocol: 0,
	}

	copy(ccd.ClientName[:], core.UTF16toLE(hostname))

	return &ccd
}

func (c *ClientCoreData) Serialize() (*bytes.Buffer, error) {
	var buff bytes.Buffer

	prefix := "rdp: client_core: serialize: %w"
	ser, err := core.Serialize(c)

	if err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	if _, err := buff.Write(ser); err != nil {
		return &buff, fmt.Errorf(prefix, err)
	}

	return &buff, nil
}
