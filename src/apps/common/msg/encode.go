package msg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// 【消息结构体解码文件  (msg.go 定义的消息结构体) struct 转化为 []byte】

// 目前支持 bool, string, []byte, struct, int

// Encode 编码 把对象变成二进制数据(字节流)  [容器类会先把信息长度写入,再写入信息]
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
	case reflect.Bool:
		writeBool(v.Bool(), buff)
	case reflect.Int8:
		writeInt8(int8(v.Int()), buff)
	case reflect.Uint8:
		writeUint8(uint8(v.Int()), buff)
	case reflect.Int16:
		writeInt16(int16(v.Int()), buff)
	case reflect.Uint16:
		writeUint16(uint16(v.Int()), buff)
	case reflect.Int32, reflect.Int:
		writeInt32(int32(v.Int()), buff)
	case reflect.Uint32, reflect.Uint:
		writeUint32(uint32(v.Int()), buff)
	case reflect.Int64:
		writeInt64(int64(v.Int()), buff)
	case reflect.Uint64:
		writeUint64(uint64(v.Int()), buff)
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
	case reflect.Slice: // todo 怎么知道slice 的元素类型？ 万一不是 []byte呢？ 而是 []string 或者 []int64? 会报错
		writeBytes(v.Bytes(), buff)
	default:
		return errors.New(fmt.Sprintf("%s, %d", "not support this type", v.Kind()))
	}
	return nil
}

func writeBool(ok bool, buff *bytes.Buffer) {
	if ok { // 为真写1, 为0写0
		writeInt8(1, buff)
	} else {
		writeInt8(0, buff)
	}
}

func writeInt8(b int8, buff *bytes.Buffer) {
	buff.WriteByte(byte(b))
	//buf := make([]byte, 1)
	//buf[0] = byte(b)
	//buff.Write(buf)
}

func writeUint8(b uint8, buff *bytes.Buffer) {
	buff.WriteByte(byte(b))
	//buf := make([]byte, 1)
	//buf[0] = byte(b)
	//buff.Write(buf)
}

func writeInt16(b int16, buff *bytes.Buffer) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(b))
	buff.Write(buf)
}

func writeUint16(b uint16, buff *bytes.Buffer) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, b)
	buff.Write(buf)
}

func writeInt32(b int32, buff *bytes.Buffer) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(b))
	buff.Write(buf)
}

func writeUint32(b uint32, buff *bytes.Buffer) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, b)
	buff.Write(buf)
}

func writeInt64(b int64, buff *bytes.Buffer) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(b))
	buff.Write(buf)
}

func writeUint64(b uint64, buff *bytes.Buffer) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, b)
	buff.Write(buf)
}

func writeBytes(bs []byte, buff *bytes.Buffer) {
	// 先写长度，再写数据
	writeUint32(uint32(len(bs)), buff)
	buff.Write(bs)
}

func writeString(s string, buff *bytes.Buffer) {
	// 先写长度，再写数据
	writeUint32(uint32(len(s)), buff)
	buff.Write([]byte(s))
}
