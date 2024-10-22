package main

import (
	"ChatZoo/common"
	net2 "ChatZoo/common/net"
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"time"
)

// s.conn.LocalAddr() 这是服务器自己的地址  s.conn.RemoteAddr() 是客户端的ip

const (
	TickerInterval = 3 * time.Second
	chCacheSize    = 1024
)

type _Session struct {
	sessionID   string
	conn        net.Conn
	wg          *sync.WaitGroup // 通知我的父协程我结束了
	ticker      *time.Ticker
	ctx         context.Context      // 用于接收父协程结束的信号
	cancel      context.CancelFunc   // 暂时没用
	receiveCh   chan common.IMessage // 套接字收到的消息都放这里
	receiveOver chan struct{}        // 套接字消息全部读取完毕
	//sendCh    chan common.IMessage // 需要往套接字里写的消息都放这里
}

func NewSession(appCtx context.Context, conn net.Conn, wg *sync.WaitGroup) *_Session {
	ctx, cancel := context.WithCancel(appCtx)
	return &_Session{
		wg:        wg,
		ctx:       ctx,
		cancel:    cancel,
		conn:      conn,
		ticker:    time.NewTicker(TickerInterval),
		receiveCh: make(chan common.IMessage, chCacheSize),
		//sendCh:    make(chan common.IMessage, chCacheSize),
		receiveOver: make(chan struct{}),
	}
}

func (s *_Session) procLoop() {
	s.wg.Add(1)
	defer func() {
		fmt.Println("goroutine : session proc exit", s.conn.RemoteAddr())
		s.conn.Close()
		s.wg.Done() // 给父亲发信号
	}()
	// 为什么不用 writeConn(), 因为目前有些问题没法解决
	go s.readConn() // 只管从网卡读取数据放到缓冲区去
	go s.procMsg()  // 处理缓冲区的客户端的消息
	<-s.receiveOver // 监控 readConn, procMsg 子协程全部结束了才返回

	//go s.sendHeartbeat()
}

// readConn 从套接字里持续不断地读
func (s *_Session) readConn() {
	defer close(s.receiveCh)
	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("session handleConnect receive exit signal")
			return
		default:
			msg, err := net2.ReadFromConn(s.conn)
			if err != nil { // 如果客户端关闭了, 这是err!=nil, 但是receiveCh里面还有数据, 也要处理完哦
				fmt.Println("session handleConnect conn read err ", err)
				return
			}
			// 假如是满的,会阻塞在这里, 也不会走到 s.ctx.Done()
			// 但没问题。正常情况下, 只会卡一会, 等procMsg 取出消息处理后就能写进去了, 下次for循环,会先走到 s.ctx.Done()
			// 异常情况下, bug或者cpu/内存满导致procMsg很慢, 就算走了 ctx.Done(), procMsg处理剩下的消息也慢。
			// 当前这种情况需要想办法解决
			s.receiveCh <- msg
		}
	}
}

// writeConn 从套接字里持续不断地写, 有实际问题没法解决
func (s *_Session) writeConn() {
	/*
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
	*/
}

// procMsg 处理客户端信息 一直不断使用conn 读并且处理逻辑回复
func (s *_Session) procMsg() {
	for msg := range s.receiveCh {
		before := time.Now()
		switch m := msg.(type) { // 又是反射, 迄今为止,所有的卡点都是反射
		case *common.MsgCmdReq:
			s.rpcReq(m, s.conn)
		case *common.MsgUserLogin:
			s.rpcUserLogin(m)
		default:
			fmt.Println("unknown msg ", m)
		}
		after := time.Now()
		if after.Sub(before).Milliseconds() > 500 { // 预警 todo 太慢的情况下不能一直卡住。直接把玩家踢了？
			fmt.Println("procMsg too slow msgID: ", msg.GetID())
		}
	}
	s.receiveOver <- struct{}{}
}

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
// tcp 一来一回的怎么做到等待回的？

func (s *_Session) rpcReq(msg *common.MsgCmdReq, conn net.Conn) {
	userID := msg.UserID
	if len(userID) == 0 {
		fmt.Println("rpcReq uerID is empty ")
		return
	}
	entity, entityErr := entityMgr.AddOrGetEntity(userID, conn)
	if entityErr != nil {
		fmt.Println("rpcReq AddOrGetEntity err ", entityErr)
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
		fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
		return
	}

	// 注册
	// 名字查重, 存db, 有重复的就用userID
	// 没有重复的就随机生成userID
	// 创建一个entity ( 查重, userID 有没有重复)
	fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
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
4. reflect.Call 并发不安全的, 需要一个并发安全的消息通道   ok (整了一个channel保证顺序)
5. 当有网络消息来的时候, 都是make一个新的[]byte去接收消息, 并发量大的
   情况下, 垃圾回收很慢, 可能导致短时间内存上涨至宕机。
   用pool 解决
6. 聊天逻辑
*/
