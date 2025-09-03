package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/code-by-meal/go-rdp/log"
)

type Tag string

const (
	OrderTag Tag = "order"
)

// Serializing of any structures to slice of bytes
func Serialize(obj any) ([]byte, error) {
	var buff bytes.Buffer

	v := reflect.ValueOf(obj)

	if !v.IsValid() {
		log.Dbg("Not <e>valid</> reflect type.")

		return buff.Bytes(), fmt.Errorf("value is not valid")
	}

	// Check if value is pointer and not nil
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			log.Dbg("Pointer is <e>nil</>.")

			return buff.Bytes(), fmt.Errorf("nil pointer")
		}

		v = v.Elem()
	}

	t := v.Type()

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldV := v.Field(i)
			fieldT := t.Field(i)

			// tags proccessing
			order := _GetOrder(fieldT)

			switch fieldV.Kind() {
			// uint`s
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
				if _, err := buff.Write([]byte(fieldV.String() + "\r\n")); err != nil {
					return buff.Bytes(), err
				}
			default:
				log.Dbg("Try serialize <e>unexpected</> type..")
			}
		}
	default:
		log.Dbg("Try to <d>se</>rialize reflect type <e>non structure</>.")
	}

	return buff.Bytes(), nil
}

// Try unserialize slice of bytes to struct
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
		log.Dbg(fmt.Sprintf("order: %v", order))

		if fieldV.Kind() == reflect.Pointer {
			fieldV = fieldV.Elem()
		}

		switch fieldV.Kind() {
		case reflect.Uint8:
			b, err := buff.ReadByte()

			if err != nil {
				return err
			}

			fieldV.SetUint(uint64(b))
		case reflect.Uint16:
			var tmp [2]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return err
			}

			fieldV.SetUint(uint64(order.Uint16(tmp[:])))
		case reflect.Uint32:
			var tmp [4]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return err
			}

			fieldV.SetUint(uint64(order.Uint32(tmp[:])))
		case reflect.Uint64:
			var tmp [8]byte

			if _, err := buff.Read(tmp[:]); err != nil {
				return err
			}

			fieldV.SetUint(order.Uint64(tmp[:]))
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
