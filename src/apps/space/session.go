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

// s.conn.LocalAddr() 这是服务器自己的地址  s.conn.RemoteAddr() 是客户端的ip

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
		s.rpcReq(m)
	case *common.MsgUserLogin:
		s.rpcUserLogin(m)
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

// 怎么做到客户端等待服务器返回值的？
// tcp 一来一回的怎么做到等待回的？

func (s *_Session) rpcReq(msg *common.MsgCmdReq) {
	userID := msg.UserID
	if len(userID) == 0 {
		fmt.Println("rpcReq uerID is empty ")
		return
	}
	entity, entityErr := entityMgr.GetEntity(userID)
	if entityErr != nil {
		// 报错返回
		return
	}
	v := reflect.ValueOf(entity.user)
	method := v.MethodByName(msg.MethodName)
	args, unpackErr := common.UnpackArgs(msg.Args)
	if unpackErr != nil {
		fmt.Println("session handleConnect unpackArgs err ", unpackErr)
		return
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	method.Call(in) // todo 这是并发不安全的, 需要改一下
}

func (s *_Session) rpcUserLogin(msg *common.MsgUserLogin) {
	openID := msg.OpenID
	if len(openID) == 0 {
		fmt.Println("rcpUserLogin, openID is empty ")
		return
	}
	// entity 内存中如果在了, 要先销毁后新创建, 因为session肯定不是原来的session 了
	if msg.IsVisitor {
		// 查重, 内存中没有这个entity
		// 不存数据库, 创建一个entity
		// 成功失败都要返回客户端消息
		return
	}

	// 注册
	// 名字查重, 存db, 有重复的就用userID
	// 没有重复的就随机生成userID
	// 创建一个entity ( 查重, userID 有没有重复)
}

func (s *_Session) rpcUserLogout(msg *common.MsgUserLogout) {
	userID := msg.UserID
	if len(userID) == 0 {
		// 报错返回
		return
	}
	entityMgr.DeleteEntity(userID)
}

// 如果客户端断线了, 也相当于玩家下线。 需要实时监控客户端状态

/*
本期要做的事情有：
1. 存db
2. 接收rpc请求并发送回复 怎么做
3. session层和每一个客户端互通心跳
4. reflect.Call 并发不安全的, 需要一个并发安全的消息通道
5. 当有网络消息来的时候, 都是make一个新的[]byte去接收消息, 并发量大的
   情况下, 垃圾回收很慢, 可能导致短时间内存上涨至宕机。
   用pool 解决
6. 聊天逻辑
*/
