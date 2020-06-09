// 网络底层函数  与上层接洽
package main

import (
	"context"
	"fmt"
	"net"
	"sync"
)
type Service struct {
	sess   sync.Map
	listen net.Listener
	ctx    context.Context
	cancelCtx context.CancelFunc
	on_message func(msg *Message, s *Service)
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
		sess:   sync.Map{},
		listen: l,
	}
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	return s
}

func (s *Service) ReqMessageHandle(handler func(msg *Message, s *Service)) {
	s.on_message = handler
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
	s.sess.Store(session.GetSessionUid(), session)

	defer conn.Close()
	go session.read_routine()
	go session.write_routine()

	for {
		select {
		case msg := <-session.reviceCh:
			if session.reviceCh != nil {
				go s.on_message(msg, s)
			}
		}
	}
}

func (s *Service) SendMe(msg *Message, sess *_TcpSession ) {
	sess.sendCh <- msg
}

/*
func (s *Service) SendAll(msg *Message, list []string){
	fmt.Println("service Send all")
	s.TravelSess()
	for _, v := range list {
		v2, ok := s.sess.Load(v)
		if ok {
			v2.(*_TcpSession).sendCh <- msg
			fmt.Println("sendCh has some", msg)
		}
	}
}*/

func (s *Service) SendAll(msg *Message, list []string){
	f := func(k, v interface{}) bool {
		fmt.Println("range session k is ", k)
		v.(*_TcpSession).sendCh <- msg
		return true
	}
	s.sess.Range(f)
}

func (s *Service) TravelSess(){
	f := func(k, v interface{}) bool {
		fmt.Println("range session k is ", k)
		return true
	}
	s.sess.Range(f)
}