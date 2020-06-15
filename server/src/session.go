// 网路底层函数
package main

import (
	"bytes"
	"context"
	"encoding/binary"
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

func (s *_TcpSession) Send(msg *Message) {
	s.sendCh <- msg
}

func (s *_TcpSession) Destroy() {
	s.cancelCtx()                             // 手动执行了cancel函数
	if s.conn != nil {
		_ = s.conn.Close()
	}
}

// 每次无线循环读数据，是阻塞式等待读
func (s *_TcpSession) read_routine() {

	loop:
	for {
		select {
		case <-s.ctx.Done():
			break loop
		default:
			msg, err := s._analyzeMessage(s.conn)
			if err != nil {
				fmt.Println("[read_routine] anlasize []byte failed err is", err)
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
			data, err := EEncode(msg)
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

func (s *_TcpSession) _analyzeMessage(conn net.Conn) (*Message, error){
	// 读数据
	size_buf := make([]byte, 4)
	_, err := s.conn.Read(size_buf)
	if err != nil {
		fmt.Println("conn.read data size failed err is ", err)
		return nil,err
	}

	var size int32
	buf_reader := bytes.NewReader(size_buf)
	err = binary.Read(buf_reader, binary.LittleEndian, &size)
	if err != nil {
		fmt.Println("decode size failed err is", err)
		return nil, err
	}

	data_buf := make([]byte, size -4)                    // 减去刚刚读过的size字节
	//_, err := io.ReadFull(s.conn, data_buf)           // 这个函数很奇怪
	_, err = s.conn.Read(data_buf)
	if err != nil {
		fmt.Println("conn.read data failed err =", err)
		return nil, err
	}

	// 解码
	msg, err := DDecode(data_buf)
	if err != nil {
		fmt.Println("decode data failed err is ", err)
		return nil, err
	}
	return msg, nil
}





