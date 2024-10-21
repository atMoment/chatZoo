package common

import (
	"errors"
)

var (
	ErrReadFull  = errors.New("io.ReadFull err")
	ErrConnWrite = errors.New("conn write err")
	ErrDecode    = errors.New("decode err")
	ErrEncode    = errors.New("encode err")
)

var AllMsgMap map[int32]IMessage

func init() {
	AllMsgMap = make(map[int32]IMessage)
	AllMsgMap[MsgID_CmdReq] = &MsgCmdReq{}
	AllMsgMap[MsgID_CmdRsp] = &MsgCmdRsp{}
	AllMsgMap[MsgID_UserLogin] = &MsgUserLogin{}
	AllMsgMap[MsgID_UserLogout] = &MsgUserLogout{}
}

type IMessage interface {
	GetID() int32
}

const (
	MsgID_None = iota
	MsgID_CmdReq
	MsgID_CmdRsp
	MsgID_UserLogin
	MsgID_UserLogout
)

type MsgCmdReq struct {
	MethodName string
	UserID     string
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

type MsgUserLogin struct {
	OpenID    string
	IsVisitor bool // 是否游客
}

func (m *MsgUserLogin) GetID() int32 {
	return MsgID_UserLogin
}

type MsgUserLogout struct {
	UserID string
}

func (m *MsgUserLogout) GetID() int32 {
	return MsgID_UserLogout
}
