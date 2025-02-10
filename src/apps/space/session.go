package main

import (
	"ChatZoo/common/db"
	"ChatZoo/common/login"
	mmsg "ChatZoo/common/msg"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// s.conn.LocalAddr() 这是服务器自己的地址  s.conn.RemoteAddr() 是客户端的ip
/*
放弃 开三个协程 (读套接字、处理消息、往套接字写)
原有如下：
1. 没法很好处理服务器关闭与写协程的关系
2. 读套接字协程读到数据后放到channel里, 处理消息协程从channel中拿。 channel有长度限制,满了会阻塞
3. 有的entity允许并发, 比如srv
*/

const (
	TickerInterval           = 3 * time.Second
	SessionWaitLoginDuration = 5 * time.Second
	chCacheSize              = 1024
)

type _Session struct {
	conn         net.Conn
	wg           *sync.WaitGroup // 通知我的父协程我结束了
	ticker       *time.Ticker
	ctx          context.Context    // 用于接收父协程结束的信号
	cancel       context.CancelFunc // 暂时没用
	cacheUtil    db.ICacheUtil
	loginFail    chan struct{} // 登录失败信号
	loginSuccess chan struct{} // 登录成功信号
	//sendCh    chan common.IMessage // 需要往套接字里写的消息都放这里
}

func NewSession(appCtx context.Context, conn net.Conn, wg *sync.WaitGroup, cacheUtil db.ICacheUtil) *_Session {
	ctx, cancel := context.WithCancel(appCtx)
	s := &_Session{
		wg:           wg,
		ctx:          ctx,
		cancel:       cancel,
		conn:         conn,
		ticker:       time.NewTicker(TickerInterval),
		cacheUtil:    cacheUtil,
		loginFail:    make(chan struct{}, 1),
		loginSuccess: make(chan struct{}, 1),
	}
	return s
}

// session 只是gate user登录之前的中间态
func (s *_Session) procLoop() {
	defer s.wg.Done()
	go s.waitLogin()
	select {
	case <-s.ctx.Done(): // 等待登录中收到关服信号
		s.conn.Close() // 子协程直接退出
		fmt.Println("session handleConnect receive exit signal")
	case <-s.loginFail: // 登录失败
		s.conn.Close() // 子协程直接退出
		fmt.Println("session login fail")
	}
}

func (s *_Session) waitLogin() {
	// 超时和收到客户端一条消息不是登录都算登录失败
	go func() {
		select {
		case <-s.loginSuccess: // 收到成功就返回不再计算超时时间
		case <-time.After(SessionWaitLoginDuration):
			close(s.loginFail)
		}
	}()

	// 服务器 conn close, err 为 use of closed network connection
	msg, err := mmsg.ReadFromConn(s.conn)
	if err != nil {
		fmt.Println("waitLogin read conn err ", err)
		return
	}
	switch m := msg.(type) { // 又是反射, 迄今为止,所有的卡点都是反射
	case *mmsg.MsgUserLogin:
		createErr := s.createUser(m) // 这是一条非常特殊的消息
		var errString string
		if createErr != nil {
			errString = createErr.Error()
		}
		if err = mmsg.WriteToConn(s.conn, &mmsg.MsgUserLoginResp{
			Err: errString,
		}); err != nil {
			fmt.Println("deal user login msg", err)
		} else {
			fmt.Printf("gate login isSuccess:%v, openID:%v isVisitor:%v err:%v\n", errString == "", m.OpenID, m.IsVisitor, errString)
		}
		close(s.loginSuccess)
	default:
		close(s.loginFail)
		fmt.Println("unknown msg ", msg.GetID())
	}
}

/*
// writeConn 从套接字里持续不断地写, 有实际问题没法解决
func (s *_Session) writeConn() {
		// todo 不管是走哪个case退出,都需要做标记, 不然上层接着往 sendCh中填充数据,造成不必要的内存泄露
		// 上层需要根据标记来做操作
		for {
			select {
			case <-s.ctx.Done():
				// todo 如果要关了, 走到这里来了, sendCh里还有数据怎么办呢？
				fmt.Println("session writeConn receive exit signal")
				return
			case msg := <-s.sendCh:
				// todo 如果是空的, 会一直阻塞在这里, 不会走 s.ctx.Done(), 怎么办？  没想到好办法
				err := common.WriteToConn(s.conn, msg)
				// 如果客户端关了, err != nil, 剩下的数据就不用发了,因为发的话也发不过去
				if err != nil {
					fmt.Println("session writeConn err ", err)
					return
				}
			default:
			}
		}
}
*/

// sendHeartbeat 实时给客户端发送心跳, 告诉他我还活着 (对于服务器来说,客户端断线了, readConn会报错,服务器会立即知道)
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
// tcp 一来一回的怎么做到等待回的？  用channel, 直到回的那条来了才放开

func (s *_Session) createUser(msg *mmsg.MsgUserLogin) error {
	openID := msg.OpenID
	if len(openID) == 0 {
		return errors.New("openID is empty")
	}
	_, err := login.VerifyLoginToken(openID, s.cacheUtil, msg.PublicKey)
	if err != nil {
		return fmt.Errorf("verify token err:%v", err)
	}
	// todo 后续的消息用这个加密解密   login消息不加密吗？
	// entity 内存中如果在了, 要先销毁后新创建, 因为session肯定不是原来的session 了
	if msg.IsVisitor {
		// 查重, 内存中没有这个entity
		// 不存数据库, 创建一个entity
		// 成功失败都要返回客户端消息
		gateuser := NewGateUser(s.conn, openID)
		DefaultGateSrvEntity.AddGateUser(openID, gateuser)
		fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
		return nil
	}

	// 注册
	// 名字查重, 存db, 有重复的就用userID
	// 没有重复的就随机生成userID
	// 创建一个entity ( 查重, userID 有没有重复)
	gateuser := NewGateUser(s.conn, openID)
	DefaultGateSrvEntity.AddGateUser(openID, gateuser)
	fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
	return nil
}
