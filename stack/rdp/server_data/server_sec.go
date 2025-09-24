package serverdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	"github.com/code-by-meal/go-rdp/stack/certs"
)

// Encryption methods

type EncryptMethod uint32

const (
	None     EncryptMethod = 0x00000000
	Four0BIT EncryptMethod = 0x00000001
	One28BIT EncryptMethod = 0x00000002
	Five6BIT EncryptMethod = 0x00000008
	Fips     EncryptMethod = 0x00000010
)

func (e *EncryptMethod) Name() string {
	switch *e {
	case None:
		return "NONE"
	case Four0BIT:
		return "40BIT"
	case One28BIT:
		return "128BIT"
	case Five6BIT:
		return "56BIT"
	case Fips:
		return "FIPS"
	default:
		return "UNKNOWN"
	}
}

// Encryption level

type EncryptLevel uint32

const (
	LevelNone             EncryptLevel = 0x00000000
	LevelLow              EncryptLevel = 0x00000001
	LevelClientCompatible EncryptLevel = 0x00000002
	LevelHigh             EncryptLevel = 0x00000003
	LevelFips             EncryptLevel = 0x00000004
)

func (e *EncryptLevel) Name() string {
	switch *e {
	case LevelNone:
		return "NONE"
	case LevelLow:
		return "LOW"
	case LevelClientCompatible:
		return "CLIENT-COMPATIBLE"
	case LevelHigh:
		return "HIGHT"
	case LevelFips:
		return "FIPS"
	default:
		return "UNKNOWN"
	}
}

type SecurityData struct {
	EncryptionMethod  EncryptMethod `order:"l"`
	EncryptionLevel   EncryptLevel  `order:"l"`
	ServerRandomLen   uint32        `order:"l"`
	ServerCertLen     uint32        `order:"l"`
	ServerRandom      []byte
	ServerCertificate []byte
	Certificate       *certs.Certificate
}

func Log(p any, title string) {
	buff := bytes.NewBuffer([]byte{})

	if err := core.WriteSingleAny(buff, &p, binary.LittleEndian); err != nil {
		log.Err("rdp: client info: ", err)
	}

	log.Dbg(fmt.Sprintf("<s>%s</>", title), buff.Bytes())
}

func (s *SecurityData) Read(buff io.Reader) error {
	prefix := "server sec data: read: %w"

	if err := core.ReadSingleAny(buff, &s.EncryptionMethod, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &s.EncryptionLevel, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	// Verify encryption level/method combinations according to MS-RDPBCGR Section 5.3.2
	errInvalidCryptoConfig := fmt.Errorf("invalid cryptgraphic configuration: method: %s\tlevel: %s", s.EncryptionMethod.Name(), s.EncryptionLevel.Name())
	if s.EncryptionLevel == LevelNone {
		if s.EncryptionMethod != None {
			return fmt.Errorf(prefix, errInvalidCryptoConfig)
		}
	}

	if s.EncryptionLevel == LevelFips {
		if s.EncryptionMethod != Fips {
			return fmt.Errorf(prefix, errInvalidCryptoConfig)
		}
	}

	if s.EncryptionLevel == LevelLow || s.EncryptionLevel == LevelHigh || s.EncryptionLevel == LevelClientCompatible {
		if s.EncryptionMethod != Four0BIT && s.EncryptionMethod != Five6BIT && s.EncryptionMethod != One28BIT && s.EncryptionMethod != Fips {
			return fmt.Errorf(prefix, errInvalidCryptoConfig)
		}
	}

	if err := core.ReadSingleAny(buff, &s.ServerRandomLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &s.ServerCertLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if s.ServerRandomLen == 0 || s.ServerCertLen == 0 {
		return fmt.Errorf(prefix, fmt.Errorf("invalid server random len=%d or server cert len=%d", s.ServerRandomLen, s.ServerCertLen))
	}

	// Server Random
	serverRandom, err := core.ReadFull(buff, int(s.ServerRandomLen))

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.ServerRandom = serverRandom

	// Server Certificate
	serverCert, err := core.ReadFull(buff, int(s.ServerCertLen))

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.ServerCertificate = serverCert

	// Log Section
	Log(s.EncryptionMethod, "EncryptionMethod:")
	Log(s.EncryptionLevel, "EncryptionLevel:")
	Log(s.ServerRandomLen, "ServerRandomLen:")
	Log(s.ServerCertLen, "ServerCertLen:")
	log.Dbg("<s>Server Random:</>", serverRandom)
	log.Dbg("<s>Server Cer:</>", serverCert)

	return nil
}
