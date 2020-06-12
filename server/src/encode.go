package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

var my_buff *bytes.Buffer


// 解码，将消息转化为字节流
// 这里把data放到buff里去
func Encode(buff *bytes.Buffer, data interface{})(int, error) {
	my_buff = new(bytes.Buffer)

	err := encode(reflect.Indirect(reflect.ValueOf(data)))
	if err != nil {
		return 0, err
	}
	buff = my_buff
	return buff.Len(), nil
}

func encode(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int32:
		writeInt32(int32(v.Int()))
	case reflect.String:
		writeString(v.String())
	case reflect.Struct:
		l := v.NumField()
		for i := 0; i < l; i++ {
			err := encode(v.Field(i))
			if err != nil {
				return err
			}
		}
	default:
		return errors.New(fmt.Sprintf("%s, %d", "not support this type", v.Kind()))
	}
	return nil
}

func writeInt32(b int32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(b))
	my_buff.Write(buf)
}

func writeUInt32(b uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, b)
	my_buff.Write(buf)
}

func writeBytes(bs []byte) {
	// 先写长度，再写数据
	writeUInt32(uint32(len(bs)))
	my_buff.Write(bs)
}

func writeString(s string) {
	// 先写长度，再写数据
	writeUInt32(uint32(len(s)))
	my_buff.Write([]byte(s))
}