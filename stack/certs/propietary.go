package certs

import "io"

type RSAPublicKey struct {
	Magic   uint32 // 0x31415352
	KeyLen  uint32 // MUST be ((BitLen / 8) + 8) bytes.
	BitLen  uint32 // The number of bits in the public key modulus.
	DataLen uint32 // This value is directly related to the BitLen field and MUST be ((BitLen / 8) - 1) bytes.
	PubExp  uint32 // The public exponent of the public key.
	Modulus []byte // The modulus field contains all (BitLen / 8) bytes of the public key
	Padding []byte
}

type ProprietaryServerCertificate struct {
	DwSigAlgID        uint32 // This field MUST be set to SIGNATURE_ALG_RSA (0x00000001).
	DwKeyAlgID        uint32 // This field MUST be set to KEY_EXCHANGE_ALG_RSA (0x00000001).
	PublicKeyBlobType uint16 // This field MUST be set to BB_RSA_KEY_BLOB (0x0006).
	PublicKeyBlobLen  uint16 // The size in bytes of the PublicKeyBlob field.
	PublicKeyBlob     RSAPublicKey
	SignatureBlobType uint16 // This field is set to BB_RSA_SIGNATURE_BLOB (0x0008).
	SignatureBlobLen  uint16 // The size in bytes of the SignatureBlob field.
	SignatureBlob     []byte
	Padding           []byte
}

type Propietary struct{}

func NewPropietary() *Propietary {
	return &Propietary{}
}

func (x *Propietary) PublicKey() ([]byte, uint32) {
	return []byte{}, 0
}

func (x *Propietary) Read(r io.Reader) error {
	return nil
}

func (x *Propietary) Verify() bool {
	return false
}
