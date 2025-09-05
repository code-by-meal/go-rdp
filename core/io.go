package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/code-by-meal/go-rdp/log"
)

type Tag string

const (
	OrderTag Tag = "order"
)

func Serialize(obj any) ([]byte, error) {
	var buff bytes.Buffer

	v := reflect.ValueOf(obj)

	if !v.IsValid() {
		return buff.Bytes(), fmt.Errorf("value is not valid")
	}

	// Check if value is pointer and not nil
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return buff.Bytes(), fmt.Errorf("nil pointer")
		}

		v = v.Elem()
	}

	t := v.Type()

	switch v.Kind() { // nolint
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldV := v.Field(i)
			fieldT := t.Field(i)

			// tags proccessing
			order := _GetOrder(fieldT)

			switch fieldV.Kind() { // nolint
			case reflect.Uint8:
				if err := buff.WriteByte(byte(fieldV.Uint())); err != nil {
					return buff.Bytes(), err
				}
			case reflect.Uint16:
				tmp := make([]byte, 2)
				order.PutUint16(tmp, uint16(fieldV.Uint()))

				if _, err := buff.Write(tmp); err != nil {
					return buff.Bytes(), err
				}
			case reflect.Uint32:
				tmp := make([]byte, 4)
				order.PutUint32(tmp, uint32(fieldV.Uint()))

				if _, err := buff.Write(tmp); err != nil {
					return buff.Bytes(), err
				}
			case reflect.Uint64:
				tmp := make([]byte, 8)
				order.PutUint64(tmp, fieldV.Uint())

				if _, err := buff.Write(tmp); err != nil {
					return buff.Bytes(), err
				}
			case reflect.String:
				if _, err := buff.Write([]byte(fieldV.String())); err != nil {
					return buff.Bytes(), err
				}
			case reflect.Struct:
				strcBytes, err := Serialize(fieldV.Interface())

				if err != nil {
					return buff.Bytes(), err
				}

				if _, err := buff.Write(strcBytes); err != nil {
					return buff.Bytes(), err
				}
			default:
				log.Dbg("Try serialize <e>unexpected</> type..")
			}
		}
	default:
		log.Dbg("Trying to <d>se</>rialize reflect type <e>non structure</>.")
	}

	return buff.Bytes(), nil
}

func Unserialize(buff *bytes.Buffer, dst any) error {
	v := reflect.ValueOf(dst)

	if !v.IsValid() {
		return fmt.Errorf("reflect value is not valid")
	}

	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("reflect value must non-nil pointer")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("reflect pointer must by to struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldV := v.Field(i)
		fieldT := t.Field(i)

		// tags proccessing
		order := _GetOrder(fieldT)

		if fieldV.Kind() == reflect.Pointer {
			fieldV = fieldV.Elem()
		}

		switch fieldV.Kind() { // nolint
		case reflect.Uint8:
			b, err := buff.ReadByte()

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fieldV.SetUint(uint64(b))
		case reflect.Uint16:
			var tmp [2]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return fmt.Errorf("%w", err)
			}

			fieldV.SetUint(uint64(order.Uint16(tmp[:])))
		case reflect.Uint32:
			var tmp [4]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return fmt.Errorf("%w", err)
			}

			fieldV.SetUint(uint64(order.Uint32(tmp[:])))
		case reflect.Uint64:
			var tmp [8]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return fmt.Errorf("%w", err)
			}

			fieldV.SetUint(order.Uint64(tmp[:]))
		case reflect.String:
			// implement
		case reflect.Struct:
			deepPtr := reflect.New(fieldV.Type())

			if err := Unserialize(buff, deepPtr.Interface()); err != nil {
				return fmt.Errorf("%w", err)
			}

			fieldV.Set(deepPtr.Elem())
		default:
			log.Dbg("Try to <d>UN</>serialize reflect type <e>non structure</>.")
		}
	}

	return nil
}

func _GetOrder(field reflect.StructField) binary.ByteOrder {
	// b - BigEndian
	// l - LittleEndian
	tag := field.Tag.Get(string(OrderTag))

	if len(tag) == 0 || tag == "b" {
		return binary.BigEndian
	}

	return binary.LittleEndian
}

// Single objects reading
func ReadSingleAny(stream io.Reader, ptr any, order binary.ByteOrder) error {
	v := reflect.ValueOf(ptr)

	if !v.IsValid() {
		return fmt.Errorf("io: read single: not valid ptr")
	}

	if v.IsNil() || v.Kind() == reflect.Struct {
		return fmt.Errorf("io: read single: ptr is nil or struct")
	}

	v = v.Elem()
	msg := "io: read single (type: %s): %w"

	switch v.Kind() {
	case reflect.Uint8:
		tmp := make([]byte, 1)

		if _, err := stream.Read(tmp); err != nil {
			return fmt.Errorf(msg, v.Type(), err)
		}

		v.SetUint(uint64(tmp[0]))
	case reflect.Uint16:
		tmp := make([]byte, 2)

		if _, err := stream.Read(tmp); err != nil {
			return fmt.Errorf(msg, v.Type(), err)
		}

		v.SetUint(uint64(order.Uint16(tmp)))
	case reflect.Uint32:
		tmp := make([]byte, 4)

		if _, err := stream.Read(tmp); err != nil {
			return fmt.Errorf(msg, v.Type(), err)
		}

		v.SetUint(uint64(order.Uint32(tmp)))
	case reflect.Uint64:
		tmp := make([]byte, 8)

		if _, err := stream.Read(tmp); err != nil {
			return fmt.Errorf(msg, v.Type(), err)
		}

		v.SetUint(order.Uint64(tmp))
	default:
		log.Info(fmt.Sprintf("<e>[UNKNOWN SINGLE TYPE]</> <d>%s</>", v.Type()))
	}

	return nil
}
