package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 这里buf只是起辅助作用，提供装字节流的容器
var buf *bytes.Buffer

// 编码 将字节流转化为消息
// 这里将字节流读到obj里
func Decode(data []byte, obj interface{}) error{
	buf = bytes.NewBuffer(data)
	err := decode(reflect.Indirect(reflect.ValueOf(obj)))
	return err
}

func decode(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int8:
		n, err := readInt8()
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetInt(int64(n))
	case reflect.Uint8:
		n, err := readUint8()
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetUint(uint64(n))
	case reflect.Int32:
		n, err := readInt32()
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetInt(int64(n))
	case reflect.String:
		s, err := readString()
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetString(s)
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			err := decode(v.Field(i))
			if err != nil {
				return err
			}
	    }
	default:
		return errors.New(fmt.Sprintf("%s, %d", "not support this type ", v.Kind()))
	}
	return nil
}

func readInt8() (int8, error) {
	n, err := buf.ReadByte()
	return int8(n), err
}

func readUint8()(uint8, error) {
	return buf.ReadByte()
}

func readInt32()(int32, error) {
	buff := make([]byte, 4)
	n, err := buf.Read(buff)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("read buf failed type is int32 in decode")
	}
	return int32(binary.LittleEndian.Uint32(buff)), nil
}

func readUint32()(uint32, error) {
	buff := make([]byte, 4)
	n, err := buf.Read(buff)
	if err != nil {
		return 0, err
	}

	if n != 4 {
		return 0, errors.New("read buf failed type is uint32 in decode")
	}
	return binary.LittleEndian.Uint32(buff), nil
}

func readBytes() ([]byte, error) {
	// []byte 的size一定是放在前面的4个字节
	n, err := readUint32()
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("read buf failed type is []byte in decode")
	}

	buff := make([]byte, n)
	rn, err := buf.Read(buff)
	if err != nil || rn != int(n) {
		return nil, errors.New("read buf failed type is []byte in decode")
	}
	return buff, nil
}

func readString() (string, error) {
	str, err := readBytes()
	if err != nil {
		return "", err
	}

	return string(str), err
}

