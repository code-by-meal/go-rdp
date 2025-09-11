package core

import (
	"encoding/binary"
	"unicode/utf16"
)

func _UTF16ToLE(u []uint16) []byte {
	b := make([]byte, 2*len(u))

	for index, value := range u {
		binary.LittleEndian.PutUint16(b[index*2:], value)
	}

	return b
}

// utf-16le
func UTF16toLE(p string) []byte {
	return _UTF16ToLE(utf16.Encode([]rune(p)))
}
