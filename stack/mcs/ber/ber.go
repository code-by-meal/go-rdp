package ber

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/code-by-meal/go-rdp/core"
)

// BER Classes

type Class uint8

const (
	ClassMask Class = 0xC0
	ClassUniv Class = 0x00
	ClassAppl Class = 0x40
	ClassCtxt Class = 0x80
	ClassPriv Class = 0xC0
)

// BER Types

type Type uint8

const (
	PcMask    Type = 0x20
	Primitive Type = 0x00
	Construct Type = 0x20
)

// BER Tags

type Tag uint8

const (
	TagMask             Tag = 0x1F
	TagBool             Tag = 0x01
	TagInteger          Tag = 0x02
	TagString           Tag = 0x03
	TagOctetString      Tag = 0x04
	TagObjectIdentifier Tag = 0x06
	TagEnumerated       Tag = 0x0A
	TagSequence         Tag = 0x10
	TagSequenceOf       Tag = 0x10
)

func Pc(pc bool) Type {
	if pc {
		return Construct
	}

	return Primitive
}

// BER Write functions

func _WriteUniversalTag(stream io.Writer, tag Tag, pc bool) error {
	if _, err := stream.Write([]byte{(byte(ClassUniv) | byte(Pc(pc))) | (byte(TagMask) & byte(tag))}); err != nil {
		return fmt.Errorf("ber: write univ tag: %w", err)
	}

	return nil
}

func _WriteLength(stream io.Writer, l int) error {
	msg := "ber: write length: %w"

	switch {
	case l > 0xff:
		if _, err := stream.Write([]byte{0x80 ^ 2}); err != nil {
			return fmt.Errorf(msg, err)
		}

		if err := binary.Write(stream, binary.BigEndian, uint16(l)); err != nil {
			return fmt.Errorf(msg, err)
		}
	case l > 0x7f:
		if _, err := stream.Write([]byte{0x80 ^ 1, byte(l)}); err != nil {
			return fmt.Errorf(msg, err)
		}
	default:
		if _, err := stream.Write([]byte{byte(l)}); err != nil {
			return fmt.Errorf(msg, err)
		}
	}

	return nil
}

func WriteInteger(stream io.Writer, v int) error {
	msg := "ber: write int: %w"

	if err := _WriteUniversalTag(stream, TagInteger, false); err != nil {
		return fmt.Errorf(msg, err)
	}

	switch {
	case v <= 0xff:
		if err := _WriteLength(stream, 1); err != nil {
			return fmt.Errorf(msg, err)
		}

		if _, err := stream.Write([]byte{byte(v & 0xff)}); err != nil {
			return fmt.Errorf(msg, err)
		}
	case v <= 0xffff:
		if err := _WriteLength(stream, 2); err != nil {
			return fmt.Errorf(msg, err)
		}

		if err := binary.Write(stream, binary.BigEndian, uint16(v&0xffff)); err != nil {
			return fmt.Errorf(msg, err)
		}
	default:
		if err := _WriteLength(stream, 4); err != nil {
			return fmt.Errorf(msg, err)
		}

		if err := binary.Write(stream, binary.BigEndian, uint32(v)); err != nil {
			return fmt.Errorf(msg, err)
		}
	}

	return nil
}

func WriteBool(stream io.Writer, v bool) error {
	msg := "ber: write bool: %w"

	if err := _WriteUniversalTag(stream, TagBool, false); err != nil {
		return fmt.Errorf(msg, err)
	}

	if err := _WriteLength(stream, 1); err != nil {
		return fmt.Errorf(msg, err)
	}

	vv := byte(0x0)

	if v {
		vv = byte(0xff)
	}

	if _, err := stream.Write([]byte{vv}); err != nil {
		return fmt.Errorf(msg, err)
	}

	return nil
}

func WriteOctetString(stream io.Writer, v string) error {
	msg := "ber: write octet string: %w"

	if err := _WriteUniversalTag(stream, TagOctetString, false); err != nil {
		return fmt.Errorf(msg, err)
	}

	if err := _WriteLength(stream, len(v)); err != nil {
		return fmt.Errorf(msg, err)
	}

	if _, err := stream.Write([]byte(v)); err != nil {
		return fmt.Errorf(msg, err)
	}

	return nil
}

func WriteDomainParameters(stream io.Writer, v []byte) error {
	msg := "ber: write dp: %w"

	if err := _WriteUniversalTag(stream, TagSequence, true); err != nil {
		return fmt.Errorf(msg, err)
	}

	if err := _WriteLength(stream, len(v)); err != nil {
		return fmt.Errorf(msg, err)
	}

	if _, err := stream.Write(v); err != nil {
		return fmt.Errorf(msg, err)
	}

	return nil
}

func WriteApplicationTag(stream io.Writer, tag Tag, v []byte) error {
	msg := "ber: write appl: %w"
	b := []byte{}

	if tag > 30 {
		b = append(b, byte(ClassAppl)|byte(Construct)|byte(TagMask))
		b = append(b, byte(tag))
	} else {
		b = append(b, byte(ClassAppl)|byte(Construct)|byte(TagMask&tag))
	}

	if err := binary.Write(stream, binary.BigEndian, b); err != nil {
		return fmt.Errorf(msg, err)
	}

	if err := _WriteLength(stream, len(v)); err != nil {
		return fmt.Errorf(msg, err)
	}

	if _, err := stream.Write(v); err != nil {
		return fmt.Errorf(msg, err)
	}

	return nil
}

// BER Read functions

func _ReadUniversalTag(stream io.Reader, tag Tag, pc bool) error {
	var b uint8

	msg := "ber: read ut: %w"

	if err := core.ReadSingleAny(stream, &b, binary.BigEndian); err != nil {
		return fmt.Errorf(msg, err)
	}

	if tag != Tag(b) {
		return fmt.Errorf("ber: read ut: invalid tag")
	}

	return nil
}

func _ReadLength(stream io.Reader) (int, error) {
	return 0, nil
}
