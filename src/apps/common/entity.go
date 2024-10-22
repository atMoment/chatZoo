package common

import "net"

type IEntity interface {
	GetNetConn() net.Conn
	GetID() string
}

type EntityInfo struct {
	entityID string
	conn     net.Conn
}

func (e *EntityInfo) GetNetConn() net.Conn {
	return e.conn
}

func (e *EntityInfo) GetID() string {
	return e.entityID
}
