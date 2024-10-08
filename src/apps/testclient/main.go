package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return
	}
	defer conn.Close()
	for {
		_, err = conn.Write([]byte("hello, server"))
		if err != nil {
			fmt.Println("conn.Write failed ", err)
			return
		}
		stocMsg := make([]byte, 1024)
		_, err = conn.Read(stocMsg)
		if err != nil {
			fmt.Println("conn.Read failed ", err)
			return
		}
		fmt.Println("receive server msg ", string(stocMsg))
		time.Sleep(5 * time.Second)
	}
}
