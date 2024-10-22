package msg

import (
	"encoding/binary"
	"errors"
)

var msgPacker = &_AnyMsgPacker{}

type _AnyMsgPacker struct{}

type IByteStream interface {
	WriteUint8(v uint8) error
	WriteInt8(v int8) error
	WriteBool(v bool) error
	WriteUint16(v uint16) error
	WriteInt16(v int16) error
	WriteUint32(v uint32) error
	WriteInt32(v int32) error
	WriteUint64(v uint64) error
	WriteInt64(v int64) error
	ReadUint8() (uint8, error)
	ReadInt8() (int8, error)
	ReadBool() (bool, error)
	ReadUint16() (uint16, error)
	ReadInt16() (int16, error)
	ReadUint32() (uint32, error)
	ReadInt32() (int32, error)
	ReadUint64() (uint64, error)
	ReadInt64() (int64, error)
	WriteString(v string) error
	ReadString() (string, error)
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

func (bs *ByteStream) ReadLoop(f func(s IByteStream) (interface{}, error)) ([]interface{}, error) {
	ret := make([]interface{}, 0)
	for i := uint32(0); i < uint32(len(bs.data)); {
		arg, err := f(bs)
		if err != nil {
			return nil, err
		}
		i = bs.readPos
		ret = append(ret, arg)
	}
	return ret, nil
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

////////////////////// 工具函数 //////////////////

func (bs *ByteStream) WriteUint8(v uint8) error {
	if err := bs.writeCheck(1); err != nil {
		return err
	}
	bs.data[bs.writePos] = v //int8/uint8 只占1个字节。可表示十进制 [-128, 127] [0, 255]
	bs.writePos++
	return nil
}

func (bs *ByteStream) ReadUint8() (uint8, error) {
	if bs.readPos >= uint32(len(bs.data)) {
		return 0, errors.New("length not enough")
	}
	v := bs.data[bs.readPos]
	bs.readPos++
	return v, nil
}

func (bs *ByteStream) WriteInt8(v int8) error {
	return bs.WriteUint8(byte(v))
}

func (bs *ByteStream) ReadInt8() (int8, error) {
	v, err := bs.ReadUint8()
	return int8(v), err
}

// WriteBool bool, 1代表true, 0 代表false
func (bs *ByteStream) WriteBool(v bool) error {
	if v {
		return bs.WriteUint8(1)
	}
	return bs.WriteUint8(0)
}

func (bs *ByteStream) ReadBool() (bool, error) {
	v, err := bs.ReadUint8()
	if v == 1 {
		return true, err
	}
	return false, err
}

func (bs *ByteStream) WriteUint16(v uint16) error {
	if err := bs.writeCheck(2); err != nil {
		return err
	}
	binary.LittleEndian.PutUint16(bs.data[bs.writePos:bs.writePos+2], v)
	bs.writePos += 2
	return nil
}

func (bs *ByteStream) ReadUint16() (uint16, error) {
	if bs.readPos+2 > uint32(len(bs.data)) {
		return 0, errors.New("length not enough")
	}
	v := binary.LittleEndian.Uint16(bs.data[bs.readPos : bs.readPos+2])
	bs.readPos += 2
	return v, nil
}

func (bs *ByteStream) WriteInt16(v int16) error {
	return bs.WriteUint16(uint16(v))
}

func (bs *ByteStream) ReadInt16() (int16, error) {
	v, err := bs.ReadUint16()
	return int16(v), err
}

func (bs *ByteStream) WriteUint32(v uint32) error {
	if err := bs.writeCheck(4); err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(bs.data[bs.writePos:bs.writePos+4], v)
	bs.writePos += 4
	return nil
}

func (bs *ByteStream) ReadUint32() (uint32, error) {
	if bs.readPos+4 > uint32(len(bs.data)) {
		return 0, errors.New("length not enough")
	}
	v := binary.LittleEndian.Uint32(bs.data[bs.readPos : bs.readPos+4])
	bs.readPos += 4
	return v, nil
}

func (bs *ByteStream) WriteInt32(v int32) error {
	return bs.WriteUint32(uint32(v))
}

func (bs *ByteStream) ReadInt32() (int32, error) {
	v, err := bs.ReadUint32()
	return int32(v), err
}

func (bs *ByteStream) WriteUint64(v uint64) error {
	size := uint32(8)
	if err := bs.writeCheck(size); err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(bs.data[bs.writePos:bs.writePos+size], v)
	bs.writePos += size
	return nil
}

func (bs *ByteStream) ReadUint64() (uint64, error) {
	size := uint32(8)
	if bs.readPos+size > uint32(len(bs.data)) {
		return 0, errors.New("length not enough")
	}
	v := binary.LittleEndian.Uint64(bs.data[bs.readPos : bs.readPos+size])
	bs.readPos += size
	return v, nil
}

func (bs *ByteStream) WriteInt64(v int64) error {
	return bs.WriteUint64(uint64(v))
}

func (bs *ByteStream) ReadInt64() (int64, error) {
	v, err := bs.ReadUint64()
	return int64(v), err
}

func (bs *ByteStream) WriteString(v string) error {
	var err error // 长度为uint32够用了, 4G的数据
	if err = bs.WriteUint32(uint32(len(v))); err != nil {
		return err
	}
	for _, s := range v {
		if err = bs.WriteUint8(byte(s)); err != nil {
			return err
		}
	}
	return nil
}

func (bs *ByteStream) ReadString() (string, error) {
	length, err := bs.ReadUint32()
	if err != nil {
		return "", err
	}
	ret := make([]byte, length)
	for i := uint32(0); i < length; i++ {
		v, readUint8Err := bs.ReadUint8()
		if readUint8Err != nil {
			return "", readUint8Err
		}
		ret[i] = v
	}
	return string(ret), nil
}
