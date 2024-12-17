package common

import (
	mmsg "ChatZoo/common/msg"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type IEntityRpc interface {
	SendNotify(methodName string, arg ...interface{}) error
	SendReq(methodName string, methodArgs ...interface{}) chan *CallRet
	SingleCall(methodName string, args ...interface{}) error
	ReceiveConn() error
	SendRsp(index int32, methodRets ...interface{}) error
}

type _EntityRpc struct {
	waitRetRpc sync.Map // 如果是entity层面的, 可以用普通map
	msgIndex   int32    // 不用atomic, 但是需要保证entity的所有函数都是顺序执行的
	entity     IEntityInfo
}

func NewEntityRpc(entity IEntityInfo) IEntityRpc {
	return &_EntityRpc{
		entity: entity,
	}
}

type CallRet struct {
	Rets []interface{}
	//RetData []byte
	Err     error
	Done    chan *CallRet // 为什么像俄罗斯套娃？ 因为sync要保存结构体, 在rsp来之前要把 timeout写进去, 仅存 channel 做不了
	Timeout *time.Timer   // 不能无限等待返回
}

func (s *_EntityRpc) String() string {
	return fmt.Sprintf("[entityRpc :%v]", s.entity.GetEntityID())
}

func (s *_EntityRpc) SendNotify(methodName string, arg ...interface{}) error {
	args, err := mmsg.PackArgs(arg...)
	if err != nil {
		return fmt.Errorf("pack rets err %v", err)
	}
	msg := &mmsg.MsgNotify{
		MethodName: methodName,
		Args:       args,
	}
	err = mmsg.WriteToConn(s.entity.GetNetConn(), msg)
	return err
}

func (s *_EntityRpc) SendReq(methodName string, methodArgs ...interface{}) chan *CallRet {
	ret := &CallRet{
		Done: make(chan *CallRet, 1), // make(chan *CallRet) 就没有返回, 为什么？
	}

	args, err := mmsg.PackArgs(methodArgs...)
	if err != nil {
		ret.Err = errors.New("pack args err " + err.Error())
		ret.Done <- ret // 如果是同步的话, 因为外部没有接收的,会阻塞在这里
		close(ret.Done)
		return ret.Done
	}

	index := atomic.AddInt32(&s.msgIndex, 1)
	msg := &mmsg.MsgCmdReq{
		MethodName: methodName,
		Args:       args,
		Index:      index,
	}

	ret.Timeout = time.AfterFunc(3*time.Second, func() {
		s.waitRetRpc.Delete(index)
		ret.Err = errors.New("over time")
		ret.Done <- ret
		close(ret.Done)
	})
	// 先存再发送套接字, 收rsp是另外的协程, 防止收rsp时在store的前面,这样就找不到了
	s.waitRetRpc.Store(index, ret)

	err = mmsg.WriteToConn(s.entity.GetNetConn(), msg)
	if err != nil {
		s.waitRetRpc.Delete(index)
		ret.Err = errors.New("write to conn err " + err.Error())
		ret.Done <- ret
		close(ret.Done)
		return ret.Done
	}
	return ret.Done
}

func (s *_EntityRpc) ReceiveConn() error {
	msg, err := mmsg.ReadFromConn(s.entity.GetNetConn())
	if err != nil {
		fmt.Println(s, " ReadFromConn ", err)
		return err
	}
	before := time.Now()
	switch m := msg.(type) { // 又是反射, 迄今为止,所有的卡点都是反射
	case *mmsg.MsgCmdReq:
		reqErr := s.receiveReq(m)
		if reqErr != nil {
			fmt.Println("receiveReq err ", reqErr)
		}
	case *mmsg.MsgCmdRsp:
		reqErr := s.receiveRsp(m)
		if reqErr != nil {
			fmt.Println("receiveRsp err ", reqErr)
		}
	case *mmsg.MsgNotify:
		reqErr := s.receiveNotify(m)
		if reqErr != nil {
			fmt.Println("receiveNotify err ", reqErr)
		}
	default:
		fmt.Println("unsupported msg ", msg.GetID())
		return nil
	}
	after := time.Now()
	if after.Sub(before).Milliseconds() > 500 { // 预警 todo 太慢的情况下不能一直卡住。直接把玩家踢了？
		fmt.Println("procMsg too slow msgID: ", msg.GetID())
	}
	return nil
}

func (s *_EntityRpc) SingleCall(methodName string, args ...interface{}) error {
	method := s.analyseMethodName(methodName)
	if method.Kind() != reflect.Func || method.IsNil() {
		return fmt.Errorf("can't find methodName %v", methodName)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	s.entity.GetRpcQueue().Push(0, method, in) // 并发安全
	return nil
}

func (s *_EntityRpc) SendRsp(index int32, methodRets ...interface{}) error {
	args, err := mmsg.PackArgs(methodRets...)
	if err != nil {
		return fmt.Errorf("pack rets err %v", err)
	}
	msg := &mmsg.MsgCmdRsp{
		Rets:  args,
		Index: index,
	}
	err = mmsg.WriteToConn(s.entity.GetNetConn(), msg)
	if err != nil {
		return fmt.Errorf("write to conn err %v", err)
	}
	return nil
}

func (s *_EntityRpc) receiveRsp(msg *mmsg.MsgCmdRsp) error {
	r, ok := s.waitRetRpc.Load(msg.Index)
	if !ok {
		return errors.New("过期的回复")
	}
	ret, _ := r.(*CallRet)
	ret.Timeout.Stop()
	args, unpackErr := mmsg.UnpackArgs(msg.Rets)
	if unpackErr != nil {
		ret.Err = fmt.Errorf("unpack arg err:%v", unpackErr)
	} else {
		ret.Rets = args
	}

	ret.Done <- ret
	close(ret.Done)
	return nil
}

func (s *_EntityRpc) receiveNotify(msg *mmsg.MsgNotify) error {
	v := reflect.ValueOf(s.entity)
	method := v.MethodByName(msg.MethodName)
	if method.Kind() != reflect.Func || method.IsNil() {
		return fmt.Errorf("can't find methodName %v", msg.MethodName)
	}

	args, unpackErr := mmsg.UnpackArgs(msg.Args)
	if unpackErr != nil {
		return fmt.Errorf("session handleConnect unpackArgs err :%v", unpackErr)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	s.entity.GetRpcQueue().Push(0, method, in) // 并发安全
	//method.Call(in) // todo 这是并发不安全的, 需要改一下
	return nil
}

func (s *_EntityRpc) receiveReq(msg *mmsg.MsgCmdReq) error {
	v := reflect.ValueOf(s.entity)
	method := v.MethodByName(msg.MethodName)
	if method.Kind() != reflect.Func || method.IsNil() {
		return fmt.Errorf("can't find methodName %v", msg.MethodName)
	}
	args, unpackErr := mmsg.UnpackArgs(msg.Args)
	if unpackErr != nil {
		return fmt.Errorf("ReceiveReq unpackArgs err :%v", unpackErr)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Call panic, methodName:%v err:%v, arg:%v\n", msg.MethodName, err, args)
		}
	}()
	s.entity.GetRpcQueue().Push(msg.Index, method, in) // 并发安全
	return nil

	/*
		rets := method.Call(in) // todo 这是并发不安全的, 需要改一下
		out := make([]interface{}, len(rets))
		for i, ret := range rets {
			out[i] = ret.Interface()
		}
		return s.sendRsp(msg.Index, out...)
	*/
}

func (s *_EntityRpc) analyseMethodName(methodName string) reflect.Value {
	rets := strings.Split(methodName, ".")
	if len(rets) == 1 {
		return reflect.ValueOf(s.entity).MethodByName(methodName)
	}
	component := s.entity.GetComponent(rets[0])
	if component.IsNil() {
		panic(fmt.Sprintf("component name illegal  :%v", rets[0]))
	}
	return reflect.ValueOf(component).MethodByName(rets[1])
}
