package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)


const Request_join = 1
const Response_join = 2
const Response_chat = 3
const Request_chat = 4

func main() {
	fmt.Println("I'm client1")
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
			fmt.Println("msg is ", msg2)
		case <- done_chan:
			fmt.Println("done_chan is over")
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
		p <- input
	}
}

func ReadSocket(p chan string, conn net.Conn, d chan int) {
	for {
		data_buf := make([]byte, 1024)
		_, err := conn.Read(data_buf)

		if err != nil {
			fmt.Println("conn.read failed err is ", err)
			d <- 1
			return
		}
		p <- string(data_buf[:])

	}
}


func RequestJoin(conn net.Conn) {
	msg := NewMessage(Request_join, []byte("1"))
	data, err := Encode(msg)

	if err != nil {
		return
	}
	_, err2 := conn.Write(data)
	if err2 != nil {
		fmt.Println("conn.write data failed err is ", err2)
	}
}

func RequestChat(conn net.Conn, str string) {
	msg := NewMessage(Request_chat, []byte(str))
	data, err := Encode(msg)

	if err != nil {
		return
	}
	_, err2 := conn.Write(data)
	if err2 != nil {
		fmt.Println("conn.write data failed err is ", err2)
	}
}

