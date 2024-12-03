package common

import "net"

type IEntityInfo interface {
	GetNetConn() net.Conn
	GetEntityID() string
	GetRpc() IEntityRpc
	SetRpc(entity IEntityInfo)
	GetRpcQueue() IEntityRpcQueue
}

type EntityInfo struct {
	entityID string
	conn     net.Conn
	rpc      IEntityRpc
	rpcQueue IEntityRpcQueue
}

func NewEntityInfo(entityID string, conn net.Conn) *EntityInfo {
	self := &EntityInfo{
		entityID: entityID,
		conn:     conn,
		rpcQueue: NewRpcQueue(),
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
