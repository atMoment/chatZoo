package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

const TickerInterval = 3 * time.Second

type _Session struct {
	conn   net.Conn
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(appCtx context.Context, conn net.Conn) *_Session {
	ctx, cancel := context.WithCancel(appCtx)
	return &_Session{
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
		ticker: time.NewTicker(TickerInterval),
	}
}

func (s *_Session) proc() {
	go s.handleConnect()
	go s.sendHeartbeat()
}

func (s *_Session) handleConnect() {
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("session handleConnect exit")
			return
		default:
			// 一直不断使用conn 读
		}
	}

}

// sendHeartbeat 实时给客户端发送心跳, 告诉他我还活着
func (s *_Session) sendHeartbeat() {
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("session sendHeartbeat exit")
			return
		case <-s.ticker.C:
			data := "i'm live"
			s.conn.Write([]byte(data))
		}
	}
}
