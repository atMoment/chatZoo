package common

import (
	"fmt"
	"net"
	"reflect"
	"sync"
)

type IEntityInfo interface {
	GetNetConn() net.Conn
	GetEntityID() string
	GetRpc() IEntityRpc
	SetRpc(entity IEntityInfo)
	GetRpcQueue() IEntityRpcQueue
	GetComponent(name string) reflect.Value
	AddComponent(name string, c interface{})
	Destroy()
	Loop()
}

type EntityInfo struct {
	entityID   string
	conn       net.Conn
	rpc        IEntityRpc
	rpcQueue   IEntityRpcQueue
	components map[string]reflect.Value
}

func NewEntityInfo(entityID string, conn net.Conn) *EntityInfo {
	self := &EntityInfo{
		entityID:   entityID,
		conn:       conn,
		rpcQueue:   NewRpcQueue(),
		components: make(map[string]reflect.Value),
	}
	return self
}

func (e *EntityInfo) SetRpc(entity IEntityInfo) {
	e.rpc = NewEntityRpc(entity)
}

func (e *EntityInfo) GetNetConn() net.Conn {
	return e.conn
}

func (e *EntityInfo) GetEntityID() string {
	return e.entityID
}

func (e *EntityInfo) GetRpc() IEntityRpc {
	return e.rpc
}

func (e *EntityInfo) GetRpcQueue() IEntityRpcQueue {
	return e.rpcQueue
}

func (e *EntityInfo) GetComponent(name string) reflect.Value {
	return e.components[name]
}

func (e *EntityInfo) AddComponent(name string, c interface{}) {
	e.components[name] = reflect.ValueOf(c)
}

func (e *EntityInfo) Destroy() {
	e.conn.Close()
	fmt.Printf("entity destroy id:%v", e.entityID)
}

func (e *EntityInfo) Loop() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go e.push(wg)
	go e.pop(wg)
	wg.Wait()
	fmt.Printf("entity loop exit id:%v\n", e.entityID)
}

func (e *EntityInfo) push(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := e.GetRpc().ReceiveConn()
		// 如果客户端关闭了, 这是err!=nil, 但是消息队列中里面还有数据, 也要处理完哦
		// 可是处理消息的时候需要向客户端发送消息, conn 客户端关闭后必然也会出错
		if err != nil {
			e.GetRpcQueue().Close()
			fmt.Printf("receive conn err:%v, gate user destroy, id:%v \n", err, e.GetEntityID())
			return
		}
	}
}

func (e *EntityInfo) pop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		rets, index, isClose := e.GetRpcQueue().Pop()
		if isClose {
			return
		}
		if index == HeartBeatIndex {

		}
		if index != 0 {
			out := make([]interface{}, len(rets))
			for i, ret := range rets {
				out[i] = ret.Interface()
			}
			err := e.GetRpc().SendRsp(index, out...)
			if err != nil {
				fmt.Println("send rsp err:", err)
			}
		}
	}
}
