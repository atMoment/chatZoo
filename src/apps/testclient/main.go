package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return
	}
	defer conn.Close()
	fmt.Println("已连接计算服务器,请输入你的四则运算公式, 空格分割, \\n 为结束符, 例如 [3 * 3 + 9]")
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
		// 故而这种方式不好, 服务器在还没有接收客户端消息之前, 怎么会先知道长度呢？
		// 除非分两个包, 第一个包先传长度
		stocMsg := make([]byte, 2048) //如果超过长度, 套接字中的剩余信息会被丢掉
		_, err = conn.Read(stocMsg)
		if err != nil {
			fmt.Println("conn.Read failed ", err)
			return
		}
		fmt.Println("receive server msg ", string(stocMsg))
		ret := strings.Split(string(stocMsg), " ")
		fmt.Println(ret)
		time.Sleep(5 * time.Second)
	}
}

func receiveFromStdin() {

}
