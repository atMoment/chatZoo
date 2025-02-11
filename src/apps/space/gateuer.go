package main

import (
	"ChatZoo/common"
	"net"
	"sync"
	"time"
)

/////////////////  gate  user  ////////////////

type IGateUser interface {
	Destroy()
}

type _User struct {
	*common.EntityInfo
	entityID          string
	conn              net.Conn
	wg                *sync.WaitGroup
	heartbeatOverTime *time.Ticker
	receiveHeartbeat  chan struct{}
	isDestroy         chan struct{}

	joinRoomID string // 以后改为组件
}

func NewGateUser(conn net.Conn, entityID string) *_User {
	user := &_User{
		EntityInfo:        common.NewEntityInfo(entityID, conn),
		entityID:          entityID,
		conn:              conn,
		heartbeatOverTime: time.NewTicker(TickerInterval),
		receiveHeartbeat:  make(chan struct{}, 1),
		isDestroy:         make(chan struct{}, 1),
		wg:                &sync.WaitGroup{},
	}
	user.SetRpc(user)
	go user.Loop()
	return user
}

/*

func (u *_User) receiveHeartBeat() {
	for {
		select {
		case <-u.receiveHeartbeat:
			u.heartbeatOverTime.Reset(TickerInterval)
		case <-u.heartbeatOverTime.C:
			u.destroy()
			return
		case <-u.isDestroy:
			return
		}
	}
}

func (u *_User) sendHeartBeat() {
	for {
		select {
		case <-time.Tick(3 * time.Second):
			msg := &mmsg.HeatBeat{}
			err := mmsg.WriteToConn(u.conn, msg)
			if err != nil {
				u.destroy()
				return
			}
		}
	}
}
*/
