// 网络底层函数  与上层接洽
// 处理与客户端连接以及消息接洽
package main

import (
	"context"
	"fmt"
	"net"
)
type Service struct {
	listen net.Listener
	ctx    context.Context
	cancelCtx context.CancelFunc
	on_message func(msg *Message, sess *_TcpSession)
	on_connect func (sess *_TcpSession)
	on_disconnect func (sess *_TcpSession)
	count int
}

// create service socket
func NewService(host string, connect_type string) *Service {
	l, err := net.Listen(connect_type, host)
	if err != nil {
		fmt.Println("net.listen err = ", err)
		return nil
	}
	s :=  &Service{
		listen: l,
		count: 0,
		on_message: nil,
		on_connect: nil,
		on_disconnect: nil,
	}
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	return s
}

func (s *Service) RegisterMessageHandle(handler func(msg *Message,  sess *_TcpSession)) {
	s.on_message = handler
}

func (s *Service) RegisterConnectHandle(handler func(sess *_TcpSession)) {
	s.on_connect = handler
}

func (s *Service) RegisterDisConnectHandle(handler func(sess *_TcpSession)) {
	s.on_disconnect = handler
}


func (s *Service) AcceptConn() {
	defer s.Destroy()
	loop:
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			fmt.Println("conn.Accept err ", err)
			continue
		}
		s.count++
		fmt.Println("count is ", s.count, conn.RemoteAddr())
		select {
		case <-s.ctx.Done():
			break loop
		default:
		}

		go s.ConnectHandle(conn, conn.RemoteAddr())
	}
}

func (s *Service) Destroy() {
	s.cancelCtx()
	if s.listen != nil {
		_ = s.listen.Close()
		s.listen = nil
	}
	fmt.Println("net service closing")
}

func (s *Service) ConnectHandle(conn net.Conn, uid net.Addr) {
	session  := NewSession(conn, uid.String())

	defer conn.Close()
	go session.read_routine()
	go session.write_routine()
	if (s.on_connect != nil) {
		s.on_connect(session)
	}

	// 或者在这里开 session处理函数的协程，就不用下面的操作了。

loop:
	for {
		select {
		case msg := <-session.reviceCh:                                       // 如果多个地方写如通道，将从通道里面去取东西，再依次回复，将导致阻塞
			if session.reviceCh != nil {
				s.on_message(msg, session)
			}
			case <-session.ctx.Done() :
				go s.on_disconnect(session)
				break loop
		}
	}
}