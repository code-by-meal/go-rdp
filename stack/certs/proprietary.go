package certs

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
)

type RSAPublicKey struct {
	Magic   uint32 // 0x31415352
	KeyLen  uint32 // MUST be ((BitLen / 8) + 8) bytes.
	BitLen  uint32 // The number of bits in the public key modulus.
	DataLen uint32 // This value is directly related to the BitLen field and MUST be ((BitLen / 8) - 1) bytes.
	PubExp  uint32 // The public exponent of the public key.
	Modulus []byte // The modulus field contains all (BitLen / 8) bytes of the public key
}

type ProprietaryCertificate struct {
	DwSigAlgID        uint32 // This field MUST be set to SIGNATURE_ALG_RSA (0x00000001.
	DwKeyAlgID        uint32 // This field MUST be set to KEY_EXCHANGE_ALG_RSA (0x00000001).
	PublicKeyBlobType uint16 // This field MUST be set to BB_RSA_KEY_BLOB (0x0006).
	PublicKeyBlobLen  uint16 // The size in bytes of the PublicKeyBlob field.
	PublicKeyBlob     RSAPublicKey
	SignatureBlobType uint16 // This field is set to BB_RSA_SIGNATURE_BLOB (0x0008).
	SignatureBlobLen  uint16 // The size in bytes of the SignatureBlob field.
	SignatureBlob     []byte
}

func NewPropietary() *ProprietaryCertificate {
	return &ProprietaryCertificate{
		PublicKeyBlob: RSAPublicKey{},
	}
}

func (x *ProprietaryCertificate) Serialize(buff io.Reader) error {
	prefix := "serialize cert: %w"

	if err := core.ReadSingleAny(buff, &x.DwSigAlgID, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.DwKeyAlgID, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlobType, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlobLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlob.Magic, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlob.KeyLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlob.BitLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlob.DataLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.PublicKeyBlob.PubExp, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	modulus := make([]byte, x.PublicKeyBlob.KeyLen)

	if err := binary.Read(buff, binary.LittleEndian, &modulus); err != nil {
		return fmt.Errorf(prefix, err)
	}

	x.PublicKeyBlob.Modulus = modulus

	if err := core.ReadSingleAny(buff, &x.SignatureBlobType, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if err := core.ReadSingleAny(buff, &x.SignatureBlobLen, binary.LittleEndian); err != nil {
		return fmt.Errorf(prefix, err)
	}

	blob, err := core.ReadFull(buff, int(x.SignatureBlobLen))

	if err != nil {
		return fmt.Errorf(prefix, err)
	}

	x.SignatureBlob = blob

	return nil
}

func (x *ProprietaryCertificate) PublicKey() ([]byte, uint32) {
	return []byte{}, 0
}

func (x *ProprietaryCertificate) Read(buff io.Reader) error {
	prefix := "proprietary: read: %w"

	if err := x.Serialize(buff); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func (x *ProprietaryCertificate) Verify() bool {
	return false
}
