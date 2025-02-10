package main

import (
	"ChatZoo/common"
	mmsg "ChatZoo/common/msg"
	"fmt"
	"net"
	"sync"
	"time"
)

/////////////////  gate  user  ////////////////

type IGateUser interface {
	Destroy()
}

type _GateUser struct {
	//*_User
	*common.EntityInfo
	entityID          string
	conn              net.Conn
	wg                *sync.WaitGroup
	heartbeatOverTime *time.Ticker
	receiveHeartbeat  chan struct{}
	isDestroy         chan struct{}
}

func NewGateUser(conn net.Conn, entityID string) *_GateUser {
	gateUser := &_GateUser{
		EntityInfo:        common.NewEntityInfo(entityID, conn),
		entityID:          entityID,
		conn:              conn,
		heartbeatOverTime: time.NewTicker(TickerInterval),
		receiveHeartbeat:  make(chan struct{}, 1),
		isDestroy:         make(chan struct{}, 1),
		wg:                &sync.WaitGroup{},
	}
	spaceUser := &_User{EntityInfo: common.NewEntityInfo(entityID, conn)}
	gateUser.SetRpc(spaceUser)
	go gateUser.start()
	return gateUser
}

func (u *_GateUser) start() {
	u.wg.Add(2)
	go u.receive()
	go u.loop()
	u.wg.Wait()
}

func (u *_GateUser) Destroy() {
	u.destroy()
}

// 调用destroy 可以使协程结束吗？
func (u *_GateUser) destroy() {
	u.GetRpcQueue().Close()
	u.isDestroy <- struct{}{}
	u.conn.Close()
}

// receive  本质上是 conn.Read, 并放到队列中去, 失败了表示与客户端断联
func (u *_GateUser) receive() {
	defer u.wg.Done()
	for {
		err := u.GetRpc().ReceiveConn()
		// 如果客户端关闭了, 这是err!=nil, 但是消息队列中里面还有数据, 也要处理完哦
		// 可是处理消息的时候需要向客户端发送消息, conn 客户端关闭后必然也会出错
		if err != nil {
			u.destroy()
			fmt.Printf("receive conn err:%v, gate user destroy, id:%v \n", err, u.GetEntityID())
			return
		}
	}
}

// loop 本质上是 处理队列中的事件
func (u *_GateUser) loop() {
	defer u.wg.Done()
	for {
		rets, index, isClose := u.GetRpcQueue().Pop()
		if isClose {
			return
		}
		if index == common.HeartBeatIndex {

		}
		if index != 0 {
			out := make([]interface{}, len(rets))
			for i, ret := range rets {
				out[i] = ret.Interface()
			}
			err := u.GetRpc().SendRsp(index, out...)
			if err != nil {
				fmt.Println("send rsp err:", err)
			}
		}
	}
}

func (u *_GateUser) receiveHeartBeat() {
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

func (u *_GateUser) sendHeartBeat() {
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
