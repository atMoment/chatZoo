package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 解码主要依靠 reflect.Set(), 需要传入具体的类型

// Decode  解码, 把二进制数据(字节流)转化为具体类型 [容器类会先读取长度再读取到容器]
func Decode(data []byte, obj interface{}) error {
	// 这里buf只是起辅助作用,提供装字节流的容器
	buf := bytes.NewBuffer(data)
	err := decode(reflect.Indirect(reflect.ValueOf(obj)), buf)
	return err
}

func decode(v reflect.Value, buf *bytes.Buffer) error {
	switch v.Kind() {
	case reflect.Int8:
		n, err := readInt8(buf)
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetInt(int64(n))
	case reflect.Uint8:
		n, err := readUint8(buf)
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetUint(uint64(n))
	case reflect.Int32:
		n, err := readInt32(buf)
		if err != nil {
			return err
		}
		if !v.CanSet() {
			return errors.New("[decode] the value can't be set")
		}
		v.SetInt(int64(n))
	case reflect.String:
		s, err := readString(buf)
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
			err := decode(v.Field(i), buf)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New(fmt.Sprintf("%s, %d", "not support this type ", v.Kind()))
	}
	return nil
}

func readInt8(buf *bytes.Buffer) (int8, error) {
	n, err := buf.ReadByte()
	return int8(n), err
}

func readUint8(buf *bytes.Buffer) (uint8, error) {
	return buf.ReadByte()
}

func readInt32(buf *bytes.Buffer) (int32, error) {
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

func readUint32(buf *bytes.Buffer) (uint32, error) {
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

func readBytes(buf *bytes.Buffer) ([]byte, error) {
	// []byte 的size一定是放在前面的4个字节
	n, err := readUint32(buf)
	if err != nil {
		return nil, err
	}

	if n < 0 {
		return nil, errors.New("read buf failed type is []byte in decode")
	}

	bytesBuff := make([]byte, n)
	rn, err := buf.Read(bytesBuff)
	if err != nil || rn != int(n) {
		return nil, errors.New("read buf failed type is []byte in decode")
	}
	return bytesBuff, nil
}

func readString(buf *bytes.Buffer) (string, error) {
	str, err := readBytes(buf)
	if err != nil {
		return "", err
	}

	return string(str), err
}
