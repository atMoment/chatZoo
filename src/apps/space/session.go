package main

import (
	"ChatZoo/common"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"
)

const TickerInterval = 3 * time.Second

var (
	ErrDecodeFail    = errors.New("decode fail")
	ErrConnWriteFail = errors.New("conn write fail")
)

type _Session struct {
	sessionID string
	conn      net.Conn
	wg        *sync.WaitGroup
	ticker    *time.Ticker
	ctx       context.Context
	cancel    context.CancelFunc
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

func (s *_Session) procLoop() {
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
		default:
			s.procMsg()
		}
	}
}

// procMsg 处理客户端信息 一直不断使用conn 读并且处理逻辑回复
func (s *_Session) procMsg() {
	msg, err := common.ReadFromConn(s.conn)
	if err != nil {
		fmt.Println("session handleConnect conn read err ", err)
		return
	}
	switch m := msg.(type) { // 又是反射, 迄今为止,所有的卡点都是反射
	case *common.MsgCmdReq:
		// 注意不要直接使用客户端发的roleID就Add
		// s.conn.LocalAddr() 这是服务器自己的地址
		entity, entityErr := entityMgr.AddOrGetEntity(s.conn.RemoteAddr().String(), s.conn.RemoteAddr().String(), s.conn)
		if entityErr != nil {
			fmt.Println("session handleConnect AddOrGetEntity err ", entityErr)
			return
		}
		v := reflect.ValueOf(entity.role)
		method := v.MethodByName(m.MethodName)
		args, unpackErr := common.UnpackArgs(m.Args)
		if unpackErr != nil {
			fmt.Println("session handleConnect unpackArgs err ", unpackErr)
			return
		}
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			in[i] = reflect.ValueOf(arg)
		}
		method.Call(in)
	default:
		fmt.Println("unknown msg ", m)
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
