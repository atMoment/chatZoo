package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// Encode 编码 把对象变成二进制数据(字节流)  [会先把信息长度写入,再写入信息]
func Encode(obj interface{}) (int, []byte, error) {
	buff := new(bytes.Buffer)

	err := encode(reflect.Indirect(reflect.ValueOf(obj)), buff)
	if err != nil {
		return 0, nil, err
	}
	return buff.Len(), buff.Bytes(), nil
}

func encode(v reflect.Value, buff *bytes.Buffer) error {
	switch v.Kind() {
	case reflect.Int32:
		writeInt32(int32(v.Int()), buff)
	case reflect.String:
		writeString(v.String(), buff)
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			err := encode(v.Field(i), buff)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New(fmt.Sprintf("%s, %d", "not support this type", v.Kind()))
	}
	return nil
}

func writeInt32(b int32, buff *bytes.Buffer) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(b))
	buff.Write(buf)
}

func writeUInt32(b uint32, buff *bytes.Buffer) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, b)
	buff.Write(buf)
}

func writeBytes(bs []byte, buff *bytes.Buffer) {
	// 先写长度，再写数据
	writeUInt32(uint32(len(bs)), buff)
	buff.Write(bs)
}

func writeString(s string, buff *bytes.Buffer) {
	// 先写长度，再写数据
	writeUInt32(uint32(len(s)), buff)
	buff.Write([]byte(s))
}
