package common

import "net"

type IEntityInfo interface {
	GetNetConn() net.Conn
	GetEntityID() string
}

type EntityInfo struct {
	EntityID string
	Conn     net.Conn
}

func (e *EntityInfo) GetNetConn() net.Conn {
	return e.Conn
}

func (e *EntityInfo) GetEntityID() string {
	return e.EntityID
}
