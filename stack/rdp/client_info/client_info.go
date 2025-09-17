package clientinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
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
	AddressFamily   AddressFamily
	CbClientAddress uint16
	ClientAddress   []byte
	CbClientDir     uint16
	ClientDir       []byte
	// ClientTimeZone  []byte
	// ClientSessionID uint32
	// PerformanceFlag uint32
	// CbAutoReconnectCookie       uint16
	// AutoReconnectCookie         []byte
	// Reserved1                   uint16
	// Reserved2                   uint16
	// CbDynamicDSTTimeZoneName    uint16
	// DynamicDSTTimeZoneName      []byte
	// DynamicDaylightTimeDisabled uint16
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

func (e *ExtraInfo) Serialize(buff *bytes.Buffer) error {
	prefix := "extra info: serialize: %w"

	values := []any{e.AddressFamily, e.CbClientAddress, e.ClientAddress, e.CbClientDir, e.ClientDir}

	for _, v := range values {
		switch vv := v.(type) {
		case []byte:
			if _, err := buff.Write(vv); err != nil {
				return fmt.Errorf(prefix, err)
			}
		case uint16, uint32, uint8, AddressFamily:
			if err := core.WriteSingleAny(buff, &vv, binary.LittleEndian); err != nil {
				return fmt.Errorf(prefix, err)
			}
		}

	}

	return nil
}

func (r *Request) Serialize(buff *bytes.Buffer) error {
	prefix := "req serialize: %w"

	values := []any{r.CodePage, r.Flag, r.CbDomain, r.CbUserName, r.CbPassword, r.CbAlternateShell, r.CbWorkingDir}

	for _, v := range values {
		if err := core.WriteSingleAny(buff, &v, binary.LittleEndian); err != nil {
			return fmt.Errorf(prefix, err)
		}
	}

	arrays := [][]byte{r.Domain, r.UserName, r.Password, r.AlternateShell, r.WorkingDir}

	for _, a := range arrays {
		if _, err := buff.Write(a); err != nil {
			return fmt.Errorf(prefix, err)
		}
	}

	if err := r.ExtraInfo.Serialize(buff); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func NewRequest(
	domain string,
	username string,
	password string,
) *Request {
	r := &Request{
		Flag:           Mouse | Unicode | LogonErrors | LogonNotify | DisableCtrlAltDel | EnableWindowsKey | AutoLogon,
		Domain:         []byte{0, 0},
		UserName:       append(core.UTF16toLE(username), 0, 0),
		Password:       append(core.UTF16toLE(password), 0, 0),
		AlternateShell: []byte{0, 0},
		WorkingDir:     []byte{0, 0},
		ExtraInfo: ExtraInfo{
			AddressFamily:   IPV4,
			CbClientAddress: 0,
			ClientAddress:   []byte{0, 0},
			CbClientDir:     0,
			ClientDir:       []byte{0, 0},
			// ClientTimeZone:  make([]byte, 172),
			// ClientSessionID: 0,
			// PerformanceFlag: 0,
		},
	}

	if len(domain) != 0 {
		r.Domain = append(core.UTF16toLE(domain), 0, 0)
		r.CbDomain = uint16(len(r.Domain) - 2)
	}

	r.CbUserName = uint16(len(r.UserName) - 2)
	r.CbPassword = uint16(len(r.Password) - 2)

	return r
}

func (r *Request) Write(stream io.Writer, proto nego.NegoProtocol, intiator uint16, sessionKeys *sec.SessionKey) error {
	prefix := "rdp: client-info: write: %w"

	var buff bytes.Buffer

	if err := r.Serialize(&buff); err != nil {
		return fmt.Errorf(prefix, err)
	}

	fmt.Println("Len: ", len(buff.Bytes()), "Client Info: ", buff.Bytes())

	switch proto { // nolint
	case nego.RDP:
		r := securitydata.NewRequest(0x0848)

		if err := r.Write(stream, intiator, sessionKeys, buff.Bytes()); err != nil {
			return fmt.Errorf(prefix, err)
		}
	default:
		return fmt.Errorf(prefix, fmt.Errorf("protocol %v are implemetned to proccess", proto))
	}

	return nil
}
