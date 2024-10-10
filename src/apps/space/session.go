package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

const TickerInterval = 3 * time.Second

type _Session struct {
	conn   net.Conn
	wg     *sync.WaitGroup
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(appCtx context.Context, conn net.Conn, wg *sync.WaitGroup) *_Session {
	ctx, cancel := context.WithCancel(appCtx)
	return &_Session{
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
		conn:   conn,
		ticker: time.NewTicker(TickerInterval),
	}
}

func (s *_Session) proc() {
	s.wg.Add(1)
	defer func() {
		fmt.Println("goroutine : session proc exit", s.conn.RemoteAddr())
		s.wg.Done() // 给父亲发信号
	}()
	s.handleConnect()
	s.conn.Close()
	//go s.sendHeartbeat()
}

func (s *_Session) handleConnect() {
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("session handleConnect receive exit signal")
			return
		default: // 一直不断使用conn 读并且回复
			ctosMsg := make([]byte, 1024) // 为什么要用1024个0占位？ 不能长度为0, 容量为 1024吗？
			_, err := s.conn.Read(ctosMsg)
			if err != nil {
				fmt.Println("session handleConnect conn read err ", err)
				return
			}
			var stocMsg string
			// s.conn.LocalAddr() 这是服务器自己的地址
			fmt.Printf("I'm server, I read client from:%v content:%v\n ", s.conn.RemoteAddr(), string(ctosMsg))
			result, err := calculate(string(ctosMsg))
			if err != nil {
				stocMsg = fmt.Sprintf("session server calculate, err:%v", err)
			} else {
				//fmt.Println(string(ctosMsg))
				//stocMsg = fmt.Sprintf("%s = %d", "hh\r\n hello", result)
				//_ = result
				stocMsg = fmt.Sprintf("%s = %d", string(ctosMsg), result)
			}
			fmt.Println(stocMsg)
			_, err = s.conn.Write([]byte(stocMsg))
			if err != nil {
				fmt.Println("session handleConnect conn write err ", err)
				return
			}
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
