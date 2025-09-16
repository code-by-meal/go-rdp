package sec

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"encoding/binary"
	"fmt"

	serverdata "github.com/code-by-meal/go-rdp/stack/rdp/server_data"
)

type SessionKey struct {
	MAC8             [16]byte
	ClientEncryptKey [8]byte
	ClientDecryptKey [8]byte
	RC4In            *rc4.Cipher
	RC4Out           *rc4.Cipher
	EncryptMethod    serverdata.EncryptMethod
}

func NewSessionKey() *SessionKey {
	sk := SessionKey{}

	return &sk
}

func (s *SessionKey) Calc(clientRandom, serverRandom []byte) error {
	prefix := "sec: calc session key: %w"

	if len(clientRandom) == 0 || len(serverRandom) == 0 {
		return fmt.Errorf(prefix, fmt.Errorf("server or client rabdom is empty"))
	}

	master := _SSL3Gen(clientRandom, _Join(clientRandom, serverRandom))

	keyBlob := _SSL3Gen(master, _Join(serverRandom, clientRandom))

	copy(s.MAC8[:], keyBlob[:16])
	copy(s.ClientEncryptKey[:], keyBlob[16:32])
	copy(s.ClientDecryptKey[:], keyBlob)

	switch s.EncryptMethod { // nolint
	case serverdata.Four0BIT:
		_Weaken56to40(s.ClientEncryptKey[:])
		_Weaken56to40(s.ClientDecryptKey[:])
	case serverdata.Five6BIT:
		_WeakenTo56(s.ClientEncryptKey[:])
		_WeakenTo56(s.ClientDecryptKey[:])
	default:
		return fmt.Errorf(prefix, fmt.Errorf("unknown encryption method: %v", s.EncryptMethod))
	}

	rcOut, err := rc4.NewCipher(s.ClientEncryptKey[:])

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.RC4Out = rcOut

	rcIn, err := rc4.NewCipher(s.ClientDecryptKey[:])

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	s.RC4In = rcIn

	return nil
}

func _Join(a, b []byte) []byte {
	t := []byte{}

	t = append(t, a...)
	t = append(t, b...)

	return t
}

func _SSL3Gen(secret, seed []byte) []byte {
	out := make([]byte, 0, 48)

	for i := 1; i <= 3; i++ {
		label := bytes.Repeat([]byte{byte('A' - 1 + i)}, i)

		h := sha1.New()
		h.Write(label)
		h.Write(secret)
		h.Write(seed)
		sha := h.Sum(nil) // 20B

		m := md5.New()
		m.Write(secret)
		m.Write(sha)
		out = append(out, m.Sum(nil)...) // 16B
	}

	return out // 48B
}

func _WeakenTo56(k []byte) {
	if len(k) < 16 {
		return
	}

	for i := 8; i < 16; i++ {
		k[i] = 0
	}
}

func _Weaken56to40(k []byte) {
	if len(k) < 16 {
		return
	}

	for i := 5; i < 16; i++ {
		k[i] = 0
	}
}

func MAC8(macKey []byte, seq uint32, data []byte) [8]byte {
	k := make([]byte, 64)

	copy(k, macKey)

	ipad := bytes.Repeat([]byte{0x36}, 64)
	opad := bytes.Repeat([]byte{0x5c}, 64)

	for i := 0; i < 64; i++ {
		ipad[i] ^= k[i]
	}

	for i := 0; i < 64; i++ {
		opad[i] ^= k[i]
	}

	var seqLE [4]byte

	binary.LittleEndian.PutUint32(seqLE[:], seq)

	md := md5.New()
	md.Write(ipad)
	md.Write(seqLE[:])
	md.Write(data)
	inner := md.Sum(nil)

	md = md5.New()
	md.Write(opad)
	md.Write(inner)
	sum := md.Sum(nil)

	var out [8]byte

	copy(out[:], sum[:8])

	return out
}
