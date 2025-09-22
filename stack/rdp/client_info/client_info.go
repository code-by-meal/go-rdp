package clientinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
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
	AddressFamily               AddressFamily
	CbClientAddress             uint16
	ClientAddress               []byte
	CbClientDir                 uint16
	ClientDir                   []byte
	ClientTimeZone              []byte
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

func (e *ExtraInfo) Serialize(buff *bytes.Buffer) error {
	prefix := "extra info: serialize: %w"

	values := []any{e.AddressFamily, e.CbClientAddress, e.ClientAddress, e.CbClientDir, e.ClientDir, e.ClientTimeZone, e.ClientSessionID, e.PerformanceFlag, e.CbAutoReconnectCookie, e.AutoReconnectCookie, e.Reserved1, e.Reserved2, e.CbDynamicDSTTimeZoneName, e.DynamicDSTTimeZoneName, e.DynamicDaylightTimeDisabled}

	// Log(e.AddressFamily, "Address Family:")
	// Log(e.CbClientAddress, "CbClientAddress:")
	// log.UTF16LE("<s>ClientAddress:</>", e.ClientAddress)
	// Log(e.CbClientDir, "CbClientDir:")
	// log.UTF16LE("<s>ClientDir:</>", e.ClientDir)
	// log.Dbg("<s>ClientTimeZone:</>", e.ClientTimeZone)
	// Log(e.ClientSessionID, "ClientSessionID:")
	// Log(e.PerformanceFlag, "PerformanceFlag:")
	// Log(e.CbAutoReconnectCookie, "CbAutoReconnectCookie:")
	// log.Dbg("<s>AutoReconnectCookie:</>", e.AutoReconnectCookie)
	// Log(e.Reserved1, "Reserved1:")
	// Log(e.Reserved2, "Reserved2:")
	// Log(e.CbDynamicDSTTimeZoneName, "CbDynamicDSTTimeZoneName:")
	// log.UTF16LE("<s>DynamicDSTTimeZoneName:</>", e.DynamicDSTTimeZoneName)
	// Log(e.DynamicDaylightTimeDisabled, "DynamicDaylightTimeDisabled:")

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
		default:
			fmt.Println("Unkonwn type: extra: ", vv)
		}

	}

	return nil
}

func Log(p any, title string) {
	buff := bytes.NewBuffer([]byte{})

	if err := core.WriteSingleAny(buff, &p, binary.LittleEndian); err != nil {
		log.Err("rdp: client info: ", err)
	}

	log.Dbg(fmt.Sprintf("<s>%s</>", title), buff.Bytes())
}

func (r *Request) Serialize(buff *bytes.Buffer) error {
	prefix := "req serialize: %w"

	values := []any{r.CodePage, r.Flag, r.CbDomain, r.CbUserName, r.CbPassword, r.CbAlternateShell, r.CbWorkingDir}

	// fmt.Print("\n\n")
	// Log(r.CodePage, "CodePage:")
	// Log(r.Flag, "Flag:")
	// Log(r.CbDomain, "CbDomain:")
	// Log(r.CbUserName, "CbUserName:")
	// Log(r.CbPassword, "CbPassword:")
	// Log(r.CbAlternateShell, "CbAlternateShell:")
	// Log(r.CbWorkingDir, "CbWorkingDir:")
	//
	// log.UTF16LE("<s>Domain:</>", r.Domain)
	// log.UTF16LE("<s>Username:</>", r.UserName)
	// log.UTF16LE("<s>Password:</>", r.Password)
	// log.UTF16LE("<s>AlternateShell:</>", r.AlternateShell)
	// log.UTF16LE("<s>WorkingDir:</>", r.WorkingDir)
	//
	// fmt.Print("\n")

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

	// if RDP_VERSION >= 5 then add extra packet
	if err := r.ExtraInfo.Serialize(buff); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func NewRequest(
	ip string,
	domain string,
	username string,
	password string,
) *Request {
	r := &Request{
		Flag:           Mouse | Unicode | LogonErrors | LogonNotify | DisableCtrlAltDel | EnableWindowsKey | AutoLogon,
		Domain:         append(core.UTF16toLE(domain), 0, 0),
		UserName:       append(core.UTF16toLE(username), 0, 0),
		Password:       append(core.UTF16toLE(password), 0, 0),
		AlternateShell: []byte{0, 0},
		WorkingDir:     []byte{0, 0},
		ExtraInfo: ExtraInfo{
			AddressFamily: IPV4,
			ClientAddress: append(core.UTF16toLE(ip), 0, 0),
			ClientDir:     append(core.UTF16toLE("C:\\Windows\\System32\\mstscax.dll"), 0, 0),
			ClientTimeZone: []byte{
				0x88, 0xFF, 0xFF, 0xFF, // Bias
				0x46, 0x00, 0x4C, 0x00, 0x45, 0x00, 0x20, 0x00, 0x53, 0x00, 0x74, 0x00, 0x61, 0x00, 0x6E, 0x00, 0x64, 0x00, 0x61, 0x00, 0x72, 0x00, 0x64, 0x00, 0x20, 0x00, 0x54, 0x00, 0x69, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Standart name
				0x00, 0x00, // Year
				0xA0, 0x00, // Month
				0x00, 0x00, // Day of week
				0x05, 0x00, // Day
				0x04, 0x00, // Hour
				0x00, 0x00, // Minute
				0x00, 0x00, // Second
				0x00, 0x00, // Milisecond
				0x00, 0x00, 0x00, 0x00, // Standart Bias
				0x46, 0x00, 0x4C, 0x00, 0x45, 0x00, 0x20, 0x00, 0x44, 0x00, 0x61, 0x00, 0x79, 0x00, 0x6C, 0x00, 0x69, 0x00, 0x67, 0x00, 0x68, 0x00, 0x74, 0x00, 0x20, 0x00, 0x54, 0x00, 0x69, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Daylight Name
				0x00, 0x00, // Year
				0xA0, 0x00, // Month
				0x00, 0x00, // Day of week
				0x05, 0x00, // Day
				0x04, 0x00, // Hour
				0x00, 0x00, // Minute
				0x00, 0x00, // Second
				0x00, 0x00, // Milisecond
				0xC4, 0xFF, 0xFF, 0xFF, // Daylight BIAS
			},
			DynamicDSTTimeZoneName: core.UTF16toLE("FLE Standard Time"),
		},
	}

	r.CbUserName = uint16(len(r.UserName) - 2)
	r.CbPassword = uint16(len(r.Password) - 2)
	r.CbDomain = uint16(len(r.Domain) - 2)

	r.ExtraInfo.CbClientAddress = uint16(len(r.ExtraInfo.ClientAddress) - 2)
	r.ExtraInfo.CbClientDir = uint16(len(r.ExtraInfo.ClientDir) - 2)
	r.ExtraInfo.CbDynamicDSTTimeZoneName = uint16(len(r.ExtraInfo.DynamicDSTTimeZoneName))
	r.ExtraInfo.CbClientDir = uint16(len(r.ExtraInfo.ClientDir) - 2)

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
