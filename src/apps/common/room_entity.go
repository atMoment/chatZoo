package common

import (
	"fmt"
	"reflect"
	"strings"
)

// 没有读写套接字的能力,只有单线程的能力

type IRoomEntity interface {
	AddComponent(name string, c interface{})
	SingleCall(methodName string, args ...interface{}) error
}

type RoomEntity struct {
	entityID   string
	rpcQueue   IEntityRpcQueue
	components map[string]reflect.Value
	stopCh     chan struct{}
}

func NewRoomEntity(entityID string) *RoomEntity {
	self := &RoomEntity{
		entityID:   entityID,
		rpcQueue:   NewRpcQueue(),
		components: make(map[string]reflect.Value),
	}
	go self.loop()
	return self
}

func (e *RoomEntity) loop() {
	for {
		_, _, isClose := e.rpcQueue.Pop()
		if isClose {
			return
		}
	}
}

func (e *RoomEntity) Destroy() {
	e.stopCh <- struct{}{}
}

func (e *RoomEntity) AddComponent(name string, c interface{}) {
	e.components[name] = reflect.ValueOf(c)
}

func (e *RoomEntity) SingleCall(methodName string, args ...interface{}) error {
	method := e.analyseMethodName(methodName)
	if method.Kind() != reflect.Func || method.IsNil() {
		return fmt.Errorf("can't find methodName %v", methodName)
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	e.rpcQueue.Push(0, method, in) // 并发安全
	return nil
}

func (e *RoomEntity) analyseMethodName(methodName string) reflect.Value {
	rets := strings.Split(methodName, ".")
	if len(rets) == 1 {
		panic(fmt.Sprintf("len(ret1) == 1"))
	}
	component, ok := e.components[rets[0]]
	if !ok {
		panic(fmt.Sprintf("can't find component, name:%v", rets[0]))
	}
	if component.IsNil() {
		panic(fmt.Sprintf("component name illegal  :%v", rets[0]))
	}
	return component.MethodByName(rets[1])
}
