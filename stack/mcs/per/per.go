package per

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
)

// PER Writing functions

func WriteLength(stream io.Writer, length int) error {
	prefix := "per: write length: %w"

	if length > 0x7f {
		ptr := uint16(length | 0x8000)

		if err := core.WriteSingleAny(stream, &ptr, binary.BigEndian); err != nil {
			return fmt.Errorf(prefix, err)
		}
	} else {

		if _, err := stream.Write([]byte{uint8(length & 0xff)}); err != nil {
			return fmt.Errorf(prefix, err)
		}
	}

	return nil
}

func WriteInteger(stream io.Writer, v uint32) error {
	prefix := "per: write int: %w"

	switch {
	case v <= 0xff:
		if err := WriteLength(stream, 1); err != nil {
			return fmt.Errorf(prefix, err)
		}

		p := uint8(v)

		if err := core.WriteSingleAny(stream, &p, binary.BigEndian); err != nil {
			return fmt.Errorf(prefix, err)
		}
	case v <= 0xffff:
		if err := WriteLength(stream, 2); err != nil {
			return fmt.Errorf(prefix, err)
		}

		p := uint16(v)

		if err := core.WriteSingleAny(stream, &p, binary.BigEndian); err != nil {
			return fmt.Errorf(prefix, err)
		}
	default:
		if err := WriteLength(stream, 4); err != nil {
			return fmt.Errorf(prefix, err)
		}

		p := uint32(v)

		if err := core.WriteSingleAny(stream, &p, binary.BigEndian); err != nil {
			return fmt.Errorf(prefix, err)
		}
	}

	return nil
}

func WriteOctetString(stream io.Writer, str string, minLengt int) error {
	l := len(str)

	if l >= minLengt {
		minLengt = l - minLengt
	}

	prefix := "per: write octet str: %w"

	if err := WriteLength(stream, minLengt); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := stream.Write([]byte(str)); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func WriteNumberOfSet(stream io.Writer, n int) error {
	if _, err := stream.Write([]byte{uint8(n)}); err != nil {
		return fmt.Errorf("per: write numb of set: %w", err)
	}

	return nil
}

func WritePadding(stream io.Writer, padd int) error {
	if _, err := stream.Write(make([]byte, padd)); err != nil {
		return fmt.Errorf("per: write padding: %w", err)
	}

	return nil
}

func WriteNumericString(stream io.Writer, str string, minValue int) error {
	l := len(str)

	if l >= minValue {
		minValue = l - minValue
	}

	prefix := "per: write numeric str: %w"

	if err := WriteLength(stream, minValue); err != nil {
		return fmt.Errorf(prefix, err)
	}

	for i := 0; i < l; i += 2 {
		c1 := str[i]
		c2 := 0x30

		if i+1 < l {
			c2 = int(str[i+1])
		}

		c1 = (c1 - 0x30) % 10
		c2 = (c2 - 0x30) % 10
		num := (c1 << 4) | uint8(c2)

		if _, err := stream.Write([]byte{num}); err != nil {
			return fmt.Errorf(prefix, err)
		}
	}

	return nil
}

func WriteOID(stream io.Writer, oid []byte) error {
	t12 := oid[0]*40 + oid[1]
	prefix := "per: write oid: %w"

	// Write length
	if _, err := stream.Write([]byte{0x5}); err != nil {
		return fmt.Errorf(prefix, err)
	}

	// Writing first two tuples
	if _, err := stream.Write([]byte{t12}); err != nil {
		return fmt.Errorf(prefix, err)
	}

	if _, err := stream.Write(oid[2:6]); err != nil {
		return fmt.Errorf(prefix, err)
	}

	return nil
}

func WriteChoice(stream io.Writer, choice byte) error {
	if _, err := stream.Write([]byte{choice}); err != nil {
		return fmt.Errorf("per: write choice: %w", err)
	}

	return nil
}

func WriteSelection(stream io.Writer, s byte) error {
	if _, err := stream.Write([]byte{s}); err != nil {
		return fmt.Errorf("pre: write selection: %w", err)
	}

	return nil
}

// PER Reading functions

func ReadLength(stream io.Reader) (int, error) {
	var b uint8

	prefix := "per: read length: %w"

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return 0, fmt.Errorf(prefix, err)
	}

	if b&0x80 != 0 {
		length := int(b&^0x80) << 8

		var bb uint8

		if err := core.ReadSingleAny(stream, &bb, binary.BigEndian); err != nil {
			return 0, fmt.Errorf(prefix, err)
		}

		return length + int(bb), nil
	}

	return int(b), nil
}

func ReadChoice(stream io.Reader) (byte, error) {
	var b uint8

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return 0x0, fmt.Errorf("per: read choice: %w", err)
	}

	return b, nil
}

func ReadOID(stream io.Reader) ([]byte, error) {
	length, err := ReadLength(stream)
	prefix := "per: read oid: %w"

	if err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	if length != 5 {
		return []byte{}, fmt.Errorf(prefix, fmt.Errorf("not valid oid"))
	}

	oid := make([]byte, 6)

	var t12 uint8

	if err := core.ReadSingleAny(stream, &t12, binary.BigEndian); err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	oid[0] = t12 / 40
	oid[1] = t12 % 40

	if _, err := stream.Read(oid[2:6]); err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	return oid, nil
}

func ReadNumberOfSet(stream io.Reader) (int, error) {
	var b uint8

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return 0, fmt.Errorf("per: read numb of set: %w", err)
	}

	return int(b), nil
}

func ReadOctetString(stream io.Reader, minLen int) ([]byte, error) {
	prefix := "per: read oct string: %w"
	length, err := ReadLength(stream)

	if err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	buff := make([]byte, length+minLen)

	if _, err := stream.Read(buff); err != nil {
		return []byte{}, fmt.Errorf(prefix, err)
	}

	return buff, nil
}

func ReadInteger(stream io.Reader) (uint32, error) {
	prefix := "per: read int: %w"
	length, err := ReadLength(stream)

	if err != nil {
		return 0, fmt.Errorf(prefix, err)
	}

	switch {
	case length == 0:
		return 0, nil
	case length == 1:
		var b uint8

		if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
			return 0, fmt.Errorf(prefix, err)
		}

		return uint32(b), nil
	case length == 2:
		b, err := ReadInteger16(stream, 0)

		if err != nil {
			return 0, fmt.Errorf(prefix, err)
		}

		return uint32(b), fmt.Errorf(prefix, err)
	}

	return 0, fmt.Errorf(prefix, fmt.Errorf("invalid int length"))
}

func ReadInteger16(stream io.Reader, minV uint16) (uint16, error) {
	prefix := "per: read int16: %w"
	var b uint16

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return 0, fmt.Errorf(prefix, err)
	}

	if b > 0xffff-minV {
		return 0, fmt.Errorf("per: read int16: uint16 invalid value %0#x > %0#x", b, 0xffff-minV)
	}

	return b + minV, nil
}

func ReadEnumerated(stream io.Reader) (uint8, error) {
	var b uint8

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return 0, fmt.Errorf("per: read enumerated: %w", err)
	}

	return b, nil
}
