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
	case int:
		s.WriteUint8(argTypeInt)
		s.WriteInt(m)
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
	case []byte:
		s.WriteUint8(argTypeBytes)
		s.WriteBytes(m)
	default:
		kind := reflect.TypeOf(m).Kind()
		if kind == reflect.Slice {
			s.WriteUint8(argTypeAnySlice)
			p.WriteAnySlice(m, s)
		} else {
			return fmt.Errorf("pack failed, unsupported type:%v typeName:%v, msg:%v", m, reflect.TypeOf(msg).Name(), m)
		}
	}
	return nil
}

func (p *_AnyMsgPacker) WriteAnySlice(msg interface{}, s IByteStream) error {
	arr := reflect.ValueOf(msg)
	l := arr.Len()
	s.WriteUint32(uint32(l)) // 最多4G
	if l > 0 {

		for i := 0; i < l; i++ {
			if err := p.Pack(arr.Index(i).Interface(), s); err != nil {
				return err
			}
		}
	} else {

		//if nil slice,write a example elem
		elType := reflect.TypeOf(msg).Elem()
		examEl := NewPackValueType(elType)
		err := p.Pack(examEl, s)
		if err != nil {
			return fmt.Errorf("write nil slice example error:%v", err)
		}
	}
	// todo 就算长度为0, 也得放个模型进去
	return nil
}

// NewPackValueType 反射创建对象。基础类型用实例，struct必须是指针
func NewPackValueType(elType reflect.Type) interface{} {
	if elType.Kind() == reflect.Ptr {
		elType = elType.Elem()
		v := reflect.New(elType).Interface()
		return v
	}
	v := reflect.New(elType).Elem().Interface()
	return v
}

func (p *_AnyMsgPacker) WriteAnyMap(msg interface{}, s IByteStream) error {
	mmap := reflect.ValueOf(msg)
	l := mmap.Len()
	if err := p.Pack(uint32(l), s); err != nil {
		return err
	}
	if l <= 0 {
		return nil
	}
	it := mmap.MapRange()
	for it.Next() {
		k := it.Key().Interface()
		v := it.Value().Interface()
		if err := p.Pack(k, s); err != nil {
			return err
		}
		if err := p.Pack(v, s); err != nil {
			return err
		}
	}
	return nil
}
