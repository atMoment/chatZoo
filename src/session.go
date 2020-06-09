// 网路底层函数
package main

import (
	"context"
	"fmt"
	"net"
)

type ISession interface {
	Send(message Message) error

	AddMessageHandler(func(message Message) error)

	//GetID() string
}


type _TcpSession struct{
	conn  net.Conn
	ctx   context.Context
	cancelCtx context.CancelFunc
	uid  string
	reviceCh chan *Message
	sendCh chan  *Message
}



func NewSession(conn net.Conn, uid string) *_TcpSession {
	sess := &_TcpSession{
		conn: conn,
		uid: uid,
		reviceCh: make(chan *Message),
		sendCh: make(chan *Message),
	}
	sess.ctx, sess.cancelCtx = context.WithCancel(context.Background())
	return sess
}

func (s *_TcpSession) GetSessionUid() string {
	return s.uid
}

func (s *_TcpSession) Send(message *Message) {}

// 每次无线循环读数据，是阻塞式等待读
func (s *_TcpSession) read_routine() {

	loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		default:

/*
			// 读长度
			buf := make([]byte, 1)
			_, err := io.ReadFull(s.conn, buf)
			if err != nil {
				fmt.Println("read data length failed")
				continue
			}
			// 把数据长度存一下
			buf_reader := bytes.NewReader(buf)
			var size int32
			err = binary.Read(buf_reader, binary.LittleEndian, &size)
			if err != nil {
				fmt.Println("save data length failed")
				continue
			}
*/
			// 读数据
			data_buf := make([]byte, 1024)
			//_, err := io.ReadFull(s.conn, data_buf)           // 这个函数很奇怪
			_, err1 := s.conn.Read(data_buf)
			if err1 != nil {
				fmt.Println("read data failed err =", err1)
				break loop
			}

			// 解码
			msg, err2 := Decode(data_buf)
			if err2 != nil {
				fmt.Println("decode data failed")
				break loop
			}
			s.reviceCh <- msg
		}
	}

}

// 无限循环写数据，等待阻塞写
func (s *_TcpSession) write_routine() {
loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		case msg:= <-s.sendCh:
			// 写数据
			data, err := Encode(msg)
			if err != nil {
				fmt.Println("encode failed err is ", err)
				break loop
			}
			_, err2 := s.conn.Write(data)
			if err2 != nil {
				fmt.Println("conn.Write failed err is ", err)
				break loop
			}
		}
	}
}




