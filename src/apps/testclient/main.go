package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return
	}
	defer conn.Close()
	fmt.Println("您已进入监控范围, 现在请畅快聊天吧... 请输入, \\n 为分隔符")
	for {
		//fmt.Scanln(&word) // 从标准控制中输入,以空格分隔
		inputReader := bufio.NewReader(os.Stdin)
		input, inputErr := inputReader.ReadString('\n')
		if inputErr != nil {
			fmt.Println("os.stdin read err ", err)
			return
		}

		_, err = conn.Write([]byte(input))
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
