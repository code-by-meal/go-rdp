package sec

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"

	serverdata "github.com/code-by-meal/go-rdp/stack/rdp/server_data"
)

var (
	A  = []byte("A")
	BB = []byte("BB")
	CC = []byte("CCC")
)

type SessionKeys struct {
	EncryptMethod      serverdata.EncryptMethod
	PreMasterSecretKey []byte // 48 bytes
	MasterSecretKey    []byte // 48 bytes
	SessionKeyBlob     []byte // 48 bytes
}

func NewSessionKey(encryptMethod serverdata.EncryptMethod) *SessionKeys {
	return &SessionKeys{
		EncryptMethod:      encryptMethod,
		PreMasterSecretKey: []byte{},
		MasterSecretKey:    []byte{},
		SessionKeyBlob:     []byte{},
	}
}

func (s *SessionKeys) Calc(clientRandom, serverRandom []byte) error {
	prefix := "sec: calc session key: %w"

	if len(clientRandom) == 0 || len(serverRandom) == 0 {
		return fmt.Errorf(prefix, fmt.Errorf("server or client rabdom is empty"))
	}

	return nil
}

func (s *SessionKeys) EstablishKeys() error {
	prefix := "sess keys: establish keys: %w"

	// Security A

	// Security X

	return nil
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
