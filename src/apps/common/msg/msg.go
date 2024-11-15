package msg

// 【消息定义文件】

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
	AllMsgMap[MsgID_UserLoginResp] = &MsgUserLoginResp{}
	AllMsgMap[MsgID_UserLogout] = &MsgUserLogout{}
	AllMsgMap[MsgID_Notify] = &MsgNotify{}
}

type IMessage interface {
	GetID() int32
}

const (
	MsgID_None = iota
	MsgID_CmdReq
	MsgID_CmdRsp
	MsgID_UserLogin
	MsgID_UserLoginResp
	MsgID_UserLogout
	MsgID_Notify
)

type MsgCmdReq struct {
	MethodName string
	Args       []byte // 参数
	Index      int32  // 消息号
}

func (m *MsgCmdReq) GetID() int32 {
	return MsgID_CmdReq
}

type MsgCmdRsp struct {
	Index int32  // 消息号
	Rets  []byte // 返回值
}

func (m *MsgCmdRsp) GetID() int32 {
	return MsgID_CmdRsp
}

type MsgUserLogin struct {
	OpenID    string
	IsVisitor bool   // 是否游客
	PublicKey string // 客户端公钥
}

func (m *MsgUserLogin) GetID() int32 {
	return MsgID_UserLogin
}

type MsgUserLoginResp struct {
	Err string // 错误为nil 表示成功, 不为nil 表示失败
}

func (m *MsgUserLoginResp) GetID() int32 {
	return MsgID_UserLoginResp
}

type MsgUserLogout struct {
	UserID string
}

func (m *MsgUserLogout) GetID() int32 {
	return MsgID_UserLogout
}

type MsgNotify struct {
	MethodName string
	Args       []byte
}

func (m *MsgNotify) GetID() int32 {
	return MsgID_Notify
}
