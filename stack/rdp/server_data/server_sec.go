package serverdata

// Encryption methods

type EncryptMethod uint32

const (
	None     EncryptMethod = 0x00000000
	Four0BIT EncryptMethod = 0x00000001
	One28BIT EncryptMethod = 0x00000002
	Five6BIT EncryptMethod = 0x00000008
	Fips     EncryptMethod = 0x00000010
)

// Encryption level

type EncryptLevel uint32

const (
	LevelNone             EncryptLevel = 0x00000000
	LevelLow              EncryptLevel = 0x00000001
	LevelClientCompatible EncryptLevel = 0x00000002
	LevelHigh             EncryptLevel = 0x00000003
	LevelFips             EncryptLevel = 0x00000004
)

type SecurityData struct {
	EncryptionMethod  EncryptMethod `order:"l"`
	EncryptionLevel   EncryptLevel  `order:"l"`
	ServerRandomLen   uint32        `order:"l"`
	ServerCertLen     uint32        `order:"l"`
	ServerRandom      []byte
	ServerCertificate []byte
}
