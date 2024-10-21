package common

import (
	"fmt"
	"io"
	"net"
)

type IConnect interface {
	WriteToConn(msg IMessage) error
	ReadFromConn() (IMessage, error)
	GetMyAddr() string
	Close()
}

type _Connect struct {
	net.Conn
}

// Send 消息结构
// | MsgSize | MsgID | MsgData |
// |    4    |   4   | ------- |
// MsgSize  = len(MsgData)

// WriteToConn msg 必须是实际的结构, 不能是interface
func (c *_Connect) WriteToConn(msg IMessage) error {
	_, msgID, err := Encode(msg.GetID())
	if err != nil {
		fmt.Printf("WriteToConn failed, encode msgID:%v err:%v \n", msg.GetID(), err)
		return ErrEncode
	}

	dataSize, msgData, err := Encode(msg)
	if err != nil {
		fmt.Printf("WriteToConn failed, encode msg err:%v \n", err)
		return ErrEncode
	}
	_, size, err := Encode(int32(dataSize))
	if err != nil {
		fmt.Printf("WriteToConn failed,  encode data size:%v err:%v \n", dataSize, err)
		return ErrEncode
	}
	allData := make([]byte, 0)
	allData = append(size, msgID...)
	allData = append(allData, msgData...) // 先放总数据长度, 再放msgID, 再放msgData
	_, err = c.Write(allData)
	if err != nil {
		fmt.Printf("WriteToConn failed,  conn write err:%v \n", err)
		return ErrConnWrite
	}
	return nil
}

// ReadFromConn 转化成具体的msg结构
func (c *_Connect) ReadFromConn() (IMessage, error) {
	//////////// 读 msgSize ///////////////
	sizeBuf := make([]byte, 4) // int32 占4个字节,2 ^ 32, 可容纳4G的数据量, 非常够用了
	var err error
	_, err = io.ReadFull(c, sizeBuf)
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read size err ", err)
		return nil, ErrReadFull
	}
	var size int32
	err = Decode(sizeBuf, &size)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode size err ", err)
		return nil, ErrDecode
	}

	//////////////  读 msgID ////////////
	msgIDBuf := make([]byte, 4)       // int32 占4个字节,2 ^ 32, 可容纳4G的数据量, 非常够用了
	_, err = io.ReadFull(c, msgIDBuf) // 读消息ID
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read size err ", err)
		return nil, ErrReadFull
	}
	var msgID int32
	err = Decode(msgIDBuf, &msgID)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode size err ", err)
		return nil, ErrDecode
	}

	/////////// 读 msg  /////////////

	infoBuf := make([]byte, size)
	_, err = io.ReadFull(c, infoBuf) // 按照总长度 读套接字所有数据
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read data info err ", err)
		return nil, ErrReadFull
	}

	m, ok := AllMsgMap[msgID]
	if !ok {
		fmt.Println("ReadFromConn failed, msgID illegal ", msgID)
		return nil, ErrDecode
	}

	err = Decode(infoBuf, m)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode msg err ", err)
		return nil, ErrReadFull
	}

	return m, nil
}

func (c *_Connect) Close() {
	c.Close()
}

func (c *_Connect) GetMyAddr() string {
	return c.LocalAddr().String()
}

// 临时过渡用

func WriteToConn(conn net.Conn, msg IMessage) error {
	_, msgID, err := Encode(msg.GetID())
	if err != nil {
		fmt.Printf("WriteToConn failed, encode msgID:%v err:%v \n", msg.GetID(), err)
		return ErrEncode
	}

	dataSize, msgData, err := Encode(msg)
	if err != nil {
		fmt.Printf("WriteToConn failed, encode msg err:%v \n", err)
		return ErrEncode
	}
	_, size, err := Encode(int32(dataSize))
	if err != nil {
		fmt.Printf("WriteToConn failed,  encode data size:%v err:%v \n", dataSize, err)
		return ErrEncode
	}
	allData := make([]byte, 0)
	allData = append(size, msgID...)
	allData = append(allData, msgData...) // 先放总数据长度, 再放msgID, 再放msgData
	_, err = conn.Write(allData)
	if err != nil {
		fmt.Printf("WriteToConn failed,  conn write err:%v \n", err)
		return ErrConnWrite
	}
	return nil
}

func ReadFromConn(conn net.Conn) (IMessage, error) {
	//////////// 读 msgSize ///////////////
	sizeBuf := make([]byte, 4) // int32 占4个字节,2 ^ 32, 可容纳4G的数据量, 非常够用了
	var err error
	_, err = io.ReadFull(conn, sizeBuf)
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read size err ", err)
		return nil, ErrReadFull
	}
	var size int32
	err = Decode(sizeBuf, &size)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode size err ", err)
		return nil, ErrDecode
	}

	//////////////  读 msgID ////////////
	msgIDBuf := make([]byte, 4)          // int32 占4个字节,2 ^ 32, 可容纳4G的数据量, 非常够用了
	_, err = io.ReadFull(conn, msgIDBuf) // 读消息ID
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read size err ", err)
		return nil, ErrReadFull
	}
	var msgID int32
	err = Decode(msgIDBuf, &msgID)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode size err ", err)
		return nil, ErrDecode
	}

	/////////// 读 msg  /////////////

	infoBuf := make([]byte, size)
	_, err = io.ReadFull(conn, infoBuf) // 按照总长度 读套接字所有数据
	if err != nil {
		fmt.Println("ReadFromConn failed, io.ReadFull read data info err ", err)
		return nil, ErrReadFull
	}

	m, ok := AllMsgMap[msgID]
	if !ok {
		fmt.Println("ReadFromConn failed, msgID illegal ", msgID)
		return nil, ErrDecode
	}

	err = Decode(infoBuf, m)
	if err != nil {
		fmt.Println("ReadFromConn failed, decode msg err ", err)
		return nil, ErrReadFull
	}

	return m, nil
}
