package msg

import (
	"fmt"
	"reflect"
)

// 【消息参数打包文件 ...interface 转化为 []byte】

// 目前支持  bool,string,int,nil
// 函数参数无法使用decode, 因为无法事先知道这个函数参数的类型
// 想要同一个消息中变化调用函数名字和参数就可以调用不同的函数。 因此同一个消息过来无法知道要承接的函数的参数
// 其实也可以知道, 但这种做法更好

// 原理: 先把参数的类型放在最前面,后面跟参数的值的二进制数据

// PackArgs 将具体类型的参数们转化为[]byte
func PackArgs(args ...interface{}) ([]byte, error) {
	s := newByteStream(make([]byte, 8))
	for _, arg := range args {
		err := msgPacker.Pack(arg, s)
		if err != nil {
			return nil, err
		}
	}
	return s.GetUsedSlice(), nil
}

// Pack 先把类型写进去, 再把内容写进去  需要自动生成这段代码吗？
func (p *_AnyMsgPacker) Pack(msg interface{}, s IByteStream) error {
	if msg == nil {
		s.WriteUint8(argTypeNil)
		return nil
	}

	switch m := msg.(type) {
	case int8:
		s.WriteUint8(argTypeInt8)
		s.WriteInt8(m)
	case uint8:
		s.WriteUint8(argTypeUint8)
		s.WriteUint8(m)
	case uint16:
		s.WriteUint8(argTypeUint16)
		s.WriteUint16(m)
	case int16:
		s.WriteUint8(argTypeInt16)
		s.WriteInt16(m)
	case uint32:
		s.WriteUint8(argTypeUint32)
		s.WriteUint32(m)
	case int32:
		s.WriteUint8(argTypeInt32)
		s.WriteInt32(m)
	case uint64:
		s.WriteUint8(argTypeUint64)
		s.WriteUint64(m)
	case int64:
		s.WriteUint8(argTypeInt64)
		s.WriteInt64(m)
	case bool:
		s.WriteUint8(argTypeBool)
		s.WriteBool(m)
	case string:
		s.WriteUint8(argTypeString)
		s.WriteString(m)
	default:
		return fmt.Errorf("pack failed, unsupported type:%v typeName:%v, msg:%v", m, reflect.TypeOf(msg).Name(), m)
	}
	return nil
}
