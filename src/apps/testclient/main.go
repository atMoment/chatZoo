package main

import (
	"ChatZoo/common"
	"ChatZoo/testclient/logic"
	"fmt"
	"net"
	"sync"
	"time"
)

/*
读数据原始方案
stocMsg := make([]byte, 1024)
_, err := conn.Read(stocMsg)

如果超过1024长度, 套接字中的剩余信息会被丢掉; 而且如果是很少的信息, 1024长度浪费
解决方案如下:
如果要提前知道这个包的数据长度好提前make出来接受, 那么需要先发一个包, 包里只放后一个包的数据长度
但是, 长度和数据分两个包发, 如果丢包了则会出问题

更新方案
利用io.ReadFull分段式读取, 一个包里放本次数据长度和数据

_, err = io.ReadFull(c.conn, sizeBuf)
size := uint32(0)
err = common.Decode(sizeBuf, &size)
infoBuf := make([]byte, size)
_, err = io.ReadFull(c.conn, infoBuf)
*/

type _Client struct {
	conn   net.Conn
	wg     *sync.WaitGroup
	ticker *time.Ticker
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return
	}
	defer conn.Close()
	fmt.Println("已连接计算服务器,请输入你的四则运算公式, 空格分割, \\n 为结束符, 例如 [3 * 3 + 9]")

	client := &_Client{
		conn:   conn,
		wg:     &sync.WaitGroup{},
		ticker: time.NewTicker(30 * time.Second),
	}
	// 试过wg.Add(1) 放到子协程开始, 但是主协程可能等不到子协程开始就执行wg.Wait(),然后就结束程序了
	client.wg.Add(2)
	go client.receiveFromStdinAndWrite()
	go client.read()
	client.wg.Wait()
}

// 震惊！ conn直接复制可行
// receiveFromStdin 持续从标准输入中读取, 并发送给服务器
func (c *_Client) receiveFromStdinAndWrite() {
	defer func() { c.wg.Done(); fmt.Println(" receiveFromStdinAndWrite over") }()
	for {
		msg := logic.Chat(c.conn.LocalAddr().String())
		if msg == nil {
			continue
		}
		err := common.WriteToConn(c.conn, msg)
		if err != nil {
			fmt.Println("conn.Write failed ", err)
			continue
		}
	}
}

// read 持续从服务器中读取消息
func (c *_Client) read() {
	defer func() { c.wg.Done(); fmt.Println(" read over") }()
	for {
		msg, err := common.ReadFromConn(c.conn)
		if err != nil {
			fmt.Println("common.ReadFromConn err", err)
			return
		}
		switch m := msg.(type) {
		case *common.MsgCmdRsp:
			fmt.Println("client read context ", m.Arg)
		default:
			fmt.Println("client receive msg illegal ", msg.GetID())
		}
	}
}
