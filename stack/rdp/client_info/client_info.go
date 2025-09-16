package clientinfo

import (
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/stack/rdp/nego"
	securitydata "github.com/code-by-meal/go-rdp/stack/rdp/security_data"
	"github.com/code-by-meal/go-rdp/stack/sec"
)

// Address Families
type AddressFamily uint16

const (
	IPV4 AddressFamily = 0x0002
	IPV6 AddressFamily = 0x0017
)

// Info flags

type InfoFlag uint32

const (
	Mouse               InfoFlag = 0x00000001
	DisableCtrlAltDel   InfoFlag = 0x00000002
	AutoLogon           InfoFlag = 0x00000008
	Unicode             InfoFlag = 0x00000010
	MaximizeShell       InfoFlag = 0x00000020
	LogonNotify         InfoFlag = 0x00000040
	Compression         InfoFlag = 0x00000080
	EnableWindowsKey    InfoFlag = 0x00000100
	RemoteConsoleAudio  InfoFlag = 0x00002000
	ForceEncryptedCsPdu InfoFlag = 0x00004000
	Rail                InfoFlag = 0x00008000
	LogonErrors         InfoFlag = 0x00010000
	MouseHasWheel       InfoFlag = 0x00020000
	PasswordIsSCPin     InfoFlag = 0x00040000
	NoAudioPlay         InfoFlag = 0x00080000
	UsingSavedCreds     InfoFlag = 0x00100000
	AudioCapture        InfoFlag = 0x00200000
	VideoDisabled       InfoFlag = 0x00400000
	HidefRailSupported  InfoFlag = 0x02000000
)

type ExtraInfo struct {
	AddressFamily
	CbClientAddress             uint16
	ClientAddress               []byte
	CbClientDir                 uint16
	ClientDir                   []byte
	ClientTimeZone              [172]byte
	ClientSessionID             uint32
	PerformanceFlag             uint32
	CbAutoReconnectCookie       uint16
	AutoReconnectCookie         []byte
	Reserved1                   uint16
	Reserved2                   uint16
	CbDynamicDSTTimeZoneName    uint16
	DynamicDSTTimeZoneName      []byte
	DynamicDaylightTimeDisabled uint16
}

type Request struct {
	CodePage         uint32
	Flag             InfoFlag
	CbDomain         uint16
	CbUserName       uint16
	CbPassword       uint16
	CbAlternateShell uint16
	CbWorkingDir     uint16
	Domain           []byte
	UserName         []byte
	Password         []byte
	AlternateShell   []byte
	WorkingDir       []byte
	ExtraInfo
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Write(stream io.Writer, proto nego.NegoProtocol, intiator uint16, sessionKeys *sec.SessionKey) error {
	prefix := "rdp: client-info: write: %w"

	switch proto { // nolint
	case nego.RDP:
		r := securitydata.NewRequest(0x0848)

		if err := r.Write(stream, intiator, sessionKeys, []byte("ENCRYPTED DATA")); err != nil {
			return fmt.Errorf(prefix, err)
		}
	default:
		return fmt.Errorf(prefix, fmt.Errorf("protocol %v are implemetned to proccess", proto))
	}

	return nil
}
