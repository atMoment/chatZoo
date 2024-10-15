package common

import "errors"

var msgPacker = &_AnyMsgPacker{}

type _AnyMsgPacker struct{}

type IByteStream interface {
	WriteUint8(v uint8) error
	WriteInt8(v int8) error
	WriteBool(v bool) error
}

type ByteStream struct {
	data     []byte
	readPos  uint32
	writePos uint32
}

func newByteStream(data []byte) *ByteStream {
	return &ByteStream{
		data: data,
	}
}

func (bs *ByteStream) GetUsedSlice() []byte {
	if bs.data == nil || bs.writePos == 0 {
		return nil
	}
	return bs.data[0:bs.writePos]
}

const (
	_                = iota
	argTypeInt8      = 1
	argTypeUint8     = 2
	argTypeInt16     = 3
	argTypeUint16    = 4
	argTypeInt32     = 5
	argTypeUint32    = 6
	argTypeInt64     = 7
	argTypeUint64    = 8
	argTypeFloat32   = 9
	argTypeFloat64   = 10
	argTypeString    = 11
	argTypeBytearray = 12
	argTypeBool      = 13
	argTypeNil       = 14
	argTypeError     = 15

	argStreamMsg    = 20
	argProtoBuffMsg = 21

	argTypeAnySlice = 100
	argTypeAnyMap   = 101
)

func (bs *ByteStream) writeCheck(c uint32) error {
	if bs.data == nil {
		return errors.New("data is nil")
	}
	if int(bs.writePos+c) > len(bs.data) {
		// 不够了要扩容了撒
		bs.growByteBuffer(len(bs.data) * 2)
	}
	return nil
}

func (bs *ByteStream) growByteBuffer(needSize int) error {
	/*
		动态增长，最大MAX_STREAM_SIZE
		最小一次增加STREAM_GROW_MIN_SIZE
		todo 这里曾经内存泄露过, 需要慎重考虑扩容算法, 我先写个最简单的
	*/
	newData := make([]byte, needSize)
	copy(newData, bs.data)
	bs.data = newData
	return nil
}

func (bs *ByteStream) WriteUint8(v uint8) error {
	if err := bs.writeCheck(1); err != nil {
		return err
	}
	bs.data[bs.writePos] = v //int8/uint8 只占1个字节。可表示十进制 [-128, 127] [0, 255]
	bs.writePos = bs.writePos + 1
	return nil
}

func (bs *ByteStream) WriteInt8(v int8) error {
	return bs.WriteUint8(byte(v)) //todo ? 这么做合适嘛？
}

// WriteBool bool, 1代表true, 0 代表false
func (bs *ByteStream) WriteBool(v bool) error {
	if v {
		return bs.WriteUint8(1)
	}
	return bs.WriteUint8(0)
}
