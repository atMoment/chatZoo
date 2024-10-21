package main

import (
	"fmt"
	"net"
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

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7788")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return
	}

	user := NewUser(conn)
	defer user.destroy()

	if !user.login() {
		fmt.Println("user 登录失败")
		return
	}
	var moduleFunc ModuleFunc
	var ok bool
	for {
		if moduleFunc, ok = moduleFuncList[showGameHall()]; ok {
			break
		}
		fmt.Println("模块名不对, 请重新输入 ", err)
	}

	user.play(moduleFunc)
}

func showGameHall() string {
	fmt.Println("welcome to chatZoo, this is game hall. we support some game ")
	fmt.Println("[ > w < ]. [ * v * ]. [ /// - /// ]. [ ` ~ ` ]. [ :) ] ")
	var moduleName string
	fmt.Println("[chat] [guess]")
	fmt.Println("请输入选择的模块")
	fmt.Scanln(&moduleName)
	return moduleName
}

/*
先登录/游客
进入游戏大厅
选择游戏
退出游戏,回到大厅
退出账号
*/
