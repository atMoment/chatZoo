package main

import (
	"ChatZoo/common"
	mmsg "ChatZoo/common/msg"
	"context"
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
	TickerInterval = 3 * time.Second
	chCacheSize    = 1024
)

type _Session struct {
	userID    string
	conn      net.Conn
	wg        *sync.WaitGroup // 通知我的父协程我结束了
	ticker    *time.Ticker
	ctx       context.Context    // 用于接收父协程结束的信号
	cancel    context.CancelFunc // 暂时没用
	receiveCh chan mmsg.IMessage // 套接字收到的消息都放这里
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
		receiveCh: make(chan mmsg.IMessage, chCacheSize),
		//sendCh:    make(chan common.IMessage, chCacheSize),
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
	s.readConn() // 只管从网卡读取数据放到缓冲区去

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
			msg, err := mmsg.ReadFromConn(s.conn)
			if err != nil { // 如果客户端关闭了, 这是err!=nil, 但是receiveCh里面还有数据, 也要处理完哦
				fmt.Println("session handleConnect conn read err ", err)
				return
			}
			before := time.Now()
			switch m := msg.(type) { // 又是反射, 迄今为止,所有的卡点都是反射
			case *mmsg.MsgUserLogin:
				s.rpcUserLogin(m)
			case *mmsg.MsgUserLogout:
				s.rpcUserLogout()
			default:
				if len(s.userID) > 0 {
					common.DefaultSrvEntity.ReceiveMsg(s.userID, msg) // 我感觉这代码像屎一样...
				}
			}
			after := time.Now()
			if after.Sub(before).Milliseconds() > 500 { // 预警 todo 太慢的情况下不能一直卡住。直接把玩家踢了？
				fmt.Println("procMsg too slow msgID: ", msg.GetID())
			}
		}
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

func (s *_Session) rpcUserLogin(msg *mmsg.MsgUserLogin) {
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
		common.DefaultSrvEntity.AddEntity(openID, NewUser(openID, s.conn))
		s.userID = openID
		fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
		return
	}

	// 注册
	// 名字查重, 存db, 有重复的就用userID
	// 没有重复的就随机生成userID
	// 创建一个entity ( 查重, userID 有没有重复)
	common.DefaultSrvEntity.AddEntity(openID, NewUser(openID, s.conn))
	s.userID = openID
	fmt.Printf("rpcUserLogin success userID:%v isVisitor:%v\n", msg.OpenID, msg.IsVisitor)
}

func (s *_Session) rpcUserLogout() {
	common.DefaultSrvEntity.DeleteEntity(s.userID)
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
