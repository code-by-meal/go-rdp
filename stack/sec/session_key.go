package sec

import (
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"fmt"

	"github.com/code-by-meal/go-rdp/core"
	"github.com/code-by-meal/go-rdp/log"
	serverdata "github.com/code-by-meal/go-rdp/stack/rdp/server_data"
)

var (
	A    = []byte("A")
	BB   = []byte("BB")
	CCC  = []byte("CCC")
	X    = []byte("X")
	YY   = []byte("YY")
	ZZZ  = []byte("ZZZ")
	Salt = []byte{0xD1, 0x26, 0x9E}
)

type SessionKeys struct {
	EncryptMethod           serverdata.EncryptMethod
	PreMasterSecretKey      []byte // 48 bytes
	MasterSecretKey         []byte // 48 bytes
	SessionKeyBlob          []byte // 48 bytes
	ClientRandom            []byte // 32 bytes
	ServerRandom            []byte // 32 bytes
	SignKey                 []byte // 16 bytes
	EncryptKey              []byte // 16 bytes
	DecryptKey              []byte // 16 bytes
	UpdateEncryptKey        []byte // 16 bytes
	UpdateDecryptKey        []byte // 16 bytes
	Padding1                []byte // 40 bytes with \x36
	Padding2                []byte // 48 bytes with \x5c
	RC4KeyLen               int
	DecryptUseCount         int
	DecryptChecksumUseCount int
	EncryptUseCount         int
	EncryptChecksumUseCount int
	// Every 4096 packets keys need be updated
}

func NewSessionKey(encryptMethod serverdata.EncryptMethod, clientRandom, serverRandom []byte) (*SessionKeys, error) {
	if len(clientRandom) == 0 || len(serverRandom) == 0 {
		return nil, fmt.Errorf("sec: new sess keys: %w", fmt.Errorf("server or client rabdom is empty"))
	}

	return &SessionKeys{
		EncryptMethod:           encryptMethod,
		PreMasterSecretKey:      []byte{},
		MasterSecretKey:         []byte{},
		SessionKeyBlob:          []byte{},
		ClientRandom:            clientRandom,
		ServerRandom:            serverRandom,
		DecryptChecksumUseCount: 0,
		DecryptUseCount:         0,
		EncryptChecksumUseCount: 0,
		EncryptUseCount:         0,
	}, nil
}

func (s *SessionKeys) UpdateKey(key, updateKey []byte) ([]byte, error) {
	prefix := "sess key: update key: %w"

	// SHA-1 digest
	s1 := sha1.New()

	for _, arr := range [][]byte{updateKey, s.Padding1, key} {
		if _, err := s1.Write(arr); err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}
	}

	// MD-5 digest
	m5 := md5.New()

	for _, arr := range [][]byte{updateKey, s.Padding2, s1.Sum(nil)} {
		if _, err := m5.Write(arr); err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}
	}

	copy(key, m5.Sum(nil)[:16])

	// RC-4 encryption
	ciph, err := rc4.NewCipher(key)

	if err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	ciph.XORKeyStream(key, key)

	switch s.EncryptMethod { // nolint
	case serverdata.Four0BIT:
		Copy(&key, &Salt, 3)
	case serverdata.Five6BIT:
		Copy(&key, &Salt, 1)
	}

	return key, nil
}

func (s *SessionKeys) Decrypt(data []byte) ([]byte, error) {
	return []byte{}, nil
}

func (s *SessionKeys) Encrypt(data []byte) ([]byte, error) {
	prefix := "ses key: encrypt: %w"

	if s.EncryptUseCount >= 4096 {
		key, err := s.UpdateKey(s.EncryptKey, s.UpdateEncryptKey)

		if err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}

		s.EncryptKey = key
	}

	s.EncryptUseCount += 1
	s.EncryptChecksumUseCount += 1

	return []byte{}, nil
}

func (s *SessionKeys) SaltedMACSignature(data []byte) ([]byte, error) {
	prefix := "ses keys: salt mac sign: %w"
	s1 := sha1.New()

	for _, arr := range [][]byte{s.SignKey[:s.RC4KeyLen], s.Padding1, core.U32ToLE(uint32(len(data))), data, core.U32ToLE(uint32(s.EncryptChecksumUseCount))} {
		if _, err := s1.Write(arr); err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}
	}

	m5 := md5.New()

	for _, arr := range [][]byte{s.SignKey, s.Padding2, s1.Sum(nil)} {
		if _, err := m5.Write(arr); err != nil {
			return []byte{}, fmt.Errorf(prefix, err)
		}
	}

	return m5.Sum(nil)[:8], nil
}

func (s *SessionKeys) MACSignature() error {
	return nil
}

func (s *SessionKeys) Calc() error {
	prefix := "sec: calc session key: %w"

	if err := s.EstablishKeys(); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (s *SessionKeys) EstablishKeys() error {
	prefix := "sess keys: establish keys: %w"
	preMasterSecretKey := []byte{}
	preMasterSecretKey = append(preMasterSecretKey, s.ClientRandom[:24]...)
	preMasterSecretKey = append(preMasterSecretKey, s.ServerRandom[:24]...)
	s.PreMasterSecretKey = preMasterSecretKey

	log.Dbg("<s>Pre Master Secret Key:</>", preMasterSecretKey)

	// Security A
	outA := []byte{}

	for _, input := range [][]byte{A, BB, CCC} {
		out, err := s.SaltedHash(input, s.PreMasterSecretKey, s.ClientRandom, s.ServerRandom)

		if err != nil {
			return fmt.Errorf(prefix, err)
		}

		outA = append(outA, out...)
	}

	log.Dbg("<s>Out A:</>", outA)
	s.MasterSecretKey = outA

	// Security X
	outX := []byte{}

	for _, input := range [][]byte{X, YY, ZZZ} {
		out, err := s.SaltedHash(input, s.MasterSecretKey, s.ClientRandom, s.ServerRandom)

		if err != nil {
			return fmt.Errorf(prefix, err)
		}

		outX = append(outX, out...)
	}

	log.Dbg("<s>Out X:</>", outX)
	s.SessionKeyBlob = outX

	// Calculating decrypt session key
	decryptKey, err := s.AllowsFIPS(s.SessionKeyBlob[16:32])

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.DecryptKey = decryptKey
	log.Dbg("<s>Decrypt key:</>", decryptKey)

	// Calculating encrypt session key
	encryptKey, err := s.AllowsFIPS(s.SessionKeyBlob[32:48])

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.EncryptKey = encryptKey
	log.Dbg("<s>Encrypt Key:</>", encryptKey)

	// Weaking sessin keys
	switch s.EncryptMethod { // nolint
	case serverdata.Four0BIT:
		Copy(&s.SignKey, &Salt, 3)
		Copy(&s.DecryptKey, &Salt, 3)
		Copy(&s.EncryptKey, &Salt, 3)

		s.RC4KeyLen = 8

	case serverdata.Five6BIT:
		Copy(&s.SignKey, &Salt, 1)
		Copy(&s.DecryptKey, &Salt, 1)
		Copy(&s.EncryptKey, &Salt, 1)

		s.RC4KeyLen = 8

	case serverdata.One28BIT:
		s.RC4KeyLen = 16
	}

	// Setting paddings
	s.Padding1 = make([]byte, 40)

	for i := range 40 {
		s.Padding1[i] = 0x36
	}

	s.Padding2 = make([]byte, 48)

	for i := range 48 {
		s.Padding2[i] = 0x5c
	}

	return nil
}

func Copy(dst, src *[]byte, length int) {
	for i := range length {
		(*dst)[i] = (*src)[i]
	}
}

func (s *SessionKeys) AllowsFIPS(input []byte) ([]byte, error) {
	m5 := md5.New()

	for _, arr := range [][]byte{input, s.ClientRandom, s.ServerRandom} {
		if _, err := m5.Write(arr); err != nil {
			return []byte{}, fmt.Errorf("sess keys: allow fips: %w", err)
		}
	}

	return m5.Sum(nil), nil
}

func (s *SessionKeys) SaltedHash(input, salt1, salt2, salt3 []byte) ([]byte, error) {
	prefix := "sess keys: satl hash: %w"
	empty := []byte{}

	if len(salt1) != 48 {
		return empty, fmt.Errorf(prefix, fmt.Errorf("salt 1 must be 48 byte array (len=%d)", len(salt1)))
	}

	if len(salt2) != 32 {
		return empty, fmt.Errorf(prefix, fmt.Errorf("salt 2 must be 32 byte array (len=%d)", len(salt2)))
	}

	if len(salt3) != 32 {
		return empty, fmt.Errorf(prefix, fmt.Errorf("salt 3 must be 32 byte array (len=%d)", len(salt1)))
	}

	// SHA-1 digest = SHA-1(INPUT + SALT1 + SALT2 + SALT3)
	s1 := sha1.New()

	for _, arr := range [][]byte{input, salt1, salt2, salt3} {
		if _, err := s1.Write(arr); err != nil {
			return empty, fmt.Errorf(prefix, err)
		}
	}

	shaDigest := s1.Sum(nil)

	// MD-5 digest = MD5(SALT1 + SHA-1_DIGEST)
	m5 := md5.New()

	for _, arr := range [][]byte{salt1, shaDigest} {
		if _, err := m5.Write(arr); err != nil {
			return empty, fmt.Errorf(prefix, err)
		}
	}

	var hash [16]byte

	copy(hash[:], m5.Sum(nil))

	return hash[:], nil
}
