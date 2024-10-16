package common

import (
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrReadFull  = errors.New("io.ReadFull err")
	ErrConnWrite = errors.New("conn write err")
	ErrDecode    = errors.New("decode err")
	ErrEncode    = errors.New("encode err")
)

// ReadFromConn 当前只能处理string
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

// WriteToConn msg 必须是实际的结构, 不能是interface
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

var AllMsgMap map[int32]IMessage

func init() {
	AllMsgMap = make(map[int32]IMessage)
	AllMsgMap[MsgID_CmdReq] = &MsgCmdReq{}
	AllMsgMap[MsgID_CmdRsp] = &MsgCmdRsp{}
}

type IMessage interface {
	GetID() int32
}

const (
	MsgID_None = iota
	MsgID_CmdReq
	MsgID_CmdRsp
)

type MsgCmdReq struct {
	MethodName string
	RoleID     string
	Args       []byte
}

func (m *MsgCmdReq) GetID() int32 {
	return MsgID_CmdReq
}

type MsgCmdRsp struct {
	Arg string
}

func (m *MsgCmdRsp) GetID() int32 {
	return MsgID_CmdRsp
}
