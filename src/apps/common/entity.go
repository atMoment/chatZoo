package common

import (
	"net"
	"reflect"
)

type IEntityInfo interface {
	GetNetConn() net.Conn
	GetEntityID() string
	GetRpc() IEntityRpc
	SetRpc(entity IEntityInfo)
	GetRpcQueue() IEntityRpcQueue
	GetComponent(name string) reflect.Value
	AddComponent(name string, c interface{})
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
