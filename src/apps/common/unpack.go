package common

import (
	"fmt"
)

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
	default:
		return 0, fmt.Errorf("unpack failed, unsupported type %v", argType)
	}
}
