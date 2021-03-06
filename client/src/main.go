package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)


func main() {
	fmt.Println("I'm client")
	conn, err := net.Dial("tcp", "127.0.0.1:3344")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	defer conn.Close()

	RequestJoin(conn)

	input_chan := make(chan string)
	socket_chan := make(chan string)
	done_chan := make(chan int)

	go ReadStdin(input_chan, done_chan)
	go ReadSocket(socket_chan, conn, done_chan)

	loop:
	for {
		select {
		case msg := <- input_chan :
			RequestChat(conn, msg)
		case msg2 := <- socket_chan:
			fmt.Println("[main] msg is ", msg2)
		case <- done_chan:
			fmt.Println("[main] done_chan is over")
			break loop
		}
	}
}

func ReadStdin(p chan string, d chan int) {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("read from stdin failed err is ", err)
			d <- 1
			return
		}
		//fmt.Printf("input read string is %v\n", input)
		p <- input
	}
}

func ReadSocket(p chan string, conn net.Conn, d chan int) {
	for {
		var str string
		err := AnalyzeMessage(conn, &str)
		if err != nil {
			fmt.Println("conn.read failed err is ", err)
			d <- 1
			return
		}
		p <- str
	}
}
func AnalyzeMessage(conn net.Conn, str *string) error{
	size_buf := make([]byte, 4)
	_, err := conn.Read(size_buf)
	if err != nil {
		fmt.Println("read data size failed err is ", err)
		return err
	}
	var size int32
	buf := bytes.NewReader(size_buf)
	err = binary.Read(buf, binary.LittleEndian, &size)
	if err != nil {
		fmt.Println("decode size failed err is ", err)
		return err
	}

	data_buf := make([]byte, size -4)     // 减去刚刚读的size字节
	_, err = conn.Read(data_buf)
	if err != nil {
		fmt.Println("read data failed err is ", err)
		return err
	}
	msg, err2 := DDecode(data_buf)
	if err2 != nil {
		fmt.Println("decode data failed err is ", err2)
		return err
	}

	switch msg.id {
	case Response_join:
		*str = msg.GetString().(*ResJoin).Words
	case Response_chat:
		*str = msg.GetString().(*ResChat).Words
	default:
		fmt.Println("[AnalyzeMessage] 现在还没有对应的消息处理")
	}

	return nil
}

func RequestJoin(conn net.Conn) {
	reqjoin := ReqJoin{1}
	msg := NewMessage(Request_join, reqjoin)
	data, err := EEncode(msg)

	if err != nil {
		fmt.Println("msg EEncode failed err is ", err)
		return
	}
	_, err2 := conn.Write(data)
	if err2 != nil {
		fmt.Println("conn.write data failed err is ", err2)
	}
}

func RequestChat(conn net.Conn, str string) {
	reqchat := ReqChat{str}
	msg := NewMessage(Request_chat, reqchat)
	data, err := EEncode(msg)

	if err != nil {
		return
	}
	_, err2 := conn.Write(data)
	if err2 != nil {
		fmt.Println("conn.write data failed err is ", err2)
	}
}

