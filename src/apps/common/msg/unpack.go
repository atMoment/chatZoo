package msg

import (
	"fmt"
	"reflect"
)

// 【消息参数解包文件  []byte 转化为 []interface】

// 先读放在最前面的类型, 再根据类型占多少字节读后面的数据

// UnpackArgs 将[]byte 转化具体类型再转化为为[]interface{}
func UnpackArgs(data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	s := newByteStream(data)
	return s.ReadLoop(msgPacker.Unpack)
}

// Unpack 先把类型写进去, 再把内容写进去  需要自动生成这段代码吗？
func (p *_AnyMsgPacker) Unpack(s IByteStream) (interface{}, error) {
	argType, err := s.ReadUint8()
	if err != nil {
		return nil, err
	}
	switch argType {
	case argTypeNil:
		return nil, nil
	case argTypeInt8:
		return s.ReadInt8()
	case argTypeUint8:
		return s.ReadUint8()
	case argTypeInt16:
		return s.ReadInt16()
	case argTypeUint16:
		return s.ReadUint16()
	case argTypeInt32:
		return s.ReadInt32()
	case argTypeUint32:
		return s.ReadUint32()
	case argTypeInt64:
		return s.ReadInt64()
	case argTypeUint64:
		return s.ReadUint64()
	case argTypeBool:
		return s.ReadBool()
	case argTypeString:
		return s.ReadString()
	case argTypeBytes:
		return s.ReadBytes()
	case argTypeAnySlice:
		return p.ReadAnySlice(s)
	default:
		return 0, fmt.Errorf("unpack failed, unsupported type %v", argType)
	}
}

func (p *_AnyMsgPacker) ReadAnySlice(s IByteStream) (interface{}, error) {
	length, err := s.ReadUint32()
	if err != nil {
		return nil, err
	}
	// 先解一个出来看看是什么
	example, err := p.Unpack(s)
	if err != nil {
		return nil, err
	}
	rSlice := make([]reflect.Value, 0)
	if length > 0 {
		rSlice = append(rSlice, reflect.ValueOf(example))
		for i := 1; i < int(length); i++ {
			example, err = p.Unpack(s)
			if err != nil {
				return nil, err
			}
			rSlice = append(rSlice, reflect.ValueOf(example))
		}
	}

	typeOfElement := reflect.TypeOf(example)             // 得到元素类型
	typeOfElementSlice := reflect.SliceOf(typeOfElement) // 得到[]元素类型
	elementSlice := reflect.MakeSlice(typeOfElementSlice, 0, 0)
	elementSliceInterface := elementSlice.Interface() // 这是啥
	el := reflect.ValueOf(&elementSliceInterface).Elem()
	val_arr1 := reflect.Append(elementSlice, rSlice...)
	el.Set(val_arr1)
	return elementSliceInterface, nil
}
