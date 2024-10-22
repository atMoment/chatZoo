package common

import (
	mmsg "ChatZoo/common/msg"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type _RpcMgr struct {
	IEntityMgr
	waitRetRpc sync.Map
	msgIndex   int32
}

type IRpcMgr interface {
	EntitySendReq(entityID string, methodName string, methodArgs ...interface{}) chan *_CallRet
	EntitySendRsp(entityID string, index int32, methodRets ...interface{}) error
	EntityReceive(entityID string) error
}

type _CallRet struct {
	Rets []interface{}
	//RetData []byte
	Err     error
	Done    chan *_CallRet // 为什么像俄罗斯套娃？ 因为sync要保存结构体, 在rsp来之前要把 timeout写进去, 仅存 channel 做不了
	Timeout *time.Timer    // 不能无限等待返回
}

func (s *_RpcMgr) EntitySendReq(entityID string, methodName string, methodArgs ...interface{}) chan *_CallRet {
	conn, _ := s.getEntityConn(entityID)
	ret := &_CallRet{
		Done: make(chan *_CallRet),
	}

	args, err := mmsg.PackArgs(methodArgs)
	if err != nil {
		ret.Err = errors.New("pack args err " + err.Error())
		ret.Done <- ret
		close(ret.Done)
		return ret.Done
	}

	index := atomic.AddInt32(&s.msgIndex, 1)
	msg := &mmsg.MsgCmdReq{
		MethodName: methodName,
		Args:       args,
		Index:      index,
	}
	err = mmsg.WriteToConn(conn, msg)
	if err != nil {
		ret.Err = errors.New("write to conn err " + err.Error())
		ret.Done <- ret
		close(ret.Done)
		return ret.Done
	}
	ret.Timeout = time.AfterFunc(3*time.Second, func() {
		s.waitRetRpc.Delete(index)
		ret.Err = errors.New("over time")
		ret.Done <- ret
		close(ret.Done)
	})
	s.waitRetRpc.Store(index, ret)
	return ret.Done
}

func (s *_RpcMgr) EntitySendRsp(entityID string, index int32, methodRets ...interface{}) error {
	conn, _ := s.getEntityConn(entityID)
	args, err := mmsg.PackArgs(methodRets)
	if err != nil {
		return fmt.Errorf("pack rets err %v", err)
	}
	msg := &mmsg.MsgCmdRsp{
		Rets:  args,
		Index: index,
	}
	err = mmsg.WriteToConn(conn, msg)
	if err != nil {
		return fmt.Errorf("write to conn err %v", err)
	}
	return nil
}

func (s *_RpcMgr) EntityReceive(entityID string) error {
	entity, err := s.GetEntity(entityID)
	if err != nil {
		return err
	}
	conn := entity.GetNetConn()
	msg, err := mmsg.ReadFromConn(conn)
	if err != nil {
		return err
	}
	switch m := msg.(type) {
	case *mmsg.MsgCmdReq:
		return s.receiveReq(entity, m)
	case *mmsg.MsgCmdRsp:
		return s.receiveRsp(m)
	default:
		return fmt.Errorf("unspport type:%v", m.GetID())
	}
}

func (s *_RpcMgr) getEntityConn(entityID string) (net.Conn, error) {
	info, err := s.GetEntity(entityID)
	if err != nil {
		panic(fmt.Sprintf("entity not exist! %v", entityID))
	}
	return info.GetNetConn(), nil
}

func (s *_RpcMgr) receiveReq(entity interface{}, msg *mmsg.MsgCmdReq) error {
	v := reflect.ValueOf(entity)
	method := v.MethodByName(msg.MethodName)
	args, unpackErr := mmsg.UnpackArgs(msg.Args)
	if unpackErr != nil {
		return fmt.Errorf("session handleConnect unpackArgs err :%v", unpackErr)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	method.Call(in) // todo 这是并发不安全的, 需要改一下
	return nil
}

func (s *_RpcMgr) receiveRsp(msg *mmsg.MsgCmdRsp) error {
	r, ok := s.waitRetRpc.Load(msg.Index)
	if !ok {
		return errors.New("过期的回复")
	}
	ret, _ := r.(*_CallRet)
	ret.Timeout.Stop()
	args, unpackErr := mmsg.UnpackArgs(msg.Rets)
	if unpackErr != nil {
		ret.Err = errors.New("unpack arg err")
	} else {
		ret.Rets = args
	}

	ret.Done <- ret
	close(ret.Done)
	return nil
}

func (s *_RpcMgr) SendNotifyToEntity(userid string, arg string) error {
	entity, err := s.GetEntity(userid)
	if err != nil {
		return err
	}
	args, err := mmsg.PackArgs(arg)
	if err != nil {
		return fmt.Errorf("pack rets err %v", err)
	}
	msg := &mmsg.MsgNotify{
		Args: args,
	}
	err = mmsg.WriteToConn(entity.conn, msg)
	return err
}

func (s *_RpcMgr) SendNotifyToEntityList(userIds map[string]struct{}, arg ...interface{}) {
	f := func(key, value any) bool {
		id, keyOk := key.(string)
		if !keyOk {
			fmt.Println("RpcToEntityList key not string")
			return false
		}
		_, find := userIds[id]
		if !find {
			return true
		}
		info, valOk := value.(*EntityInfo)
		if !valOk {
			fmt.Println("RpcToEntityList value not _EntityInfo")
			return false
		}

		args, err := mmsg.PackArgs(arg)
		if err != nil {
			fmt.Printf("pack rets err %v\n", err)
			return false
		}
		msg := &mmsg.MsgNotify{
			Args: args}
		err = mmsg.WriteToConn(info.conn, msg)
		if err != nil {
			fmt.Println("RpcToEntityList wrtie conn err ", err)
		}

		return true
	}
	s.TravelMgr(f)
}
