package main

import (
	"ChatZoo/common"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const TickerInterval = 3 * time.Second

var (
	ErrDecodeFail    = errors.New("decode fail")
	ErrConnWriteFail = errors.New("conn write fail")
)

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
			ctosMsg, err := common.ReadFromConn(s.conn)
			if err != nil {
				fmt.Println("session handleConnect conn read err ", err)
				return
			}

			err = s.replayClient(ctosMsg)
			if err != nil {
				if errors.Is(err, ErrDecodeFail) {
					continue
				}
				return
			}
		}
	}
}

// replayClient 根据从客户端中读到的信息, 解释分析再发送结果给客户端
func (s *_Session) replayClient(ctosMsg string) error {
	var stocMsg string
	// s.conn.LocalAddr() 这是服务器自己的地址
	fmt.Printf("I'm server, I read client from:%v content:%v\n ", s.conn.RemoteAddr(), string(ctosMsg))
	result, err := calculate(string(ctosMsg))
	if err != nil {
		stocMsg = fmt.Sprintf("session server calculate, err:%v", err)
	} else {
		stocMsg = fmt.Sprintf("%s = %d", string(ctosMsg), result)
	}
	common.WriteToConn(s.conn, stocMsg)
	return nil
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
