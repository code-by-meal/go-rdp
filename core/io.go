package core

import "encoding/binary"

type Tag string

const (
	OrderTag Tag = "o"
)

// Serializing of any structures to slice of bytes
func Serialize() []byte {
	data := []byte{}

	return data
}

// Try unserialize slice of bytes to struct
func Unserialize(data []byte, obj any) error {
	return nil
}

func _GetOrder() binary.ByteOrder {
	return binary.LittleEndian
}
