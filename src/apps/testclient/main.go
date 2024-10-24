package main

import (
	mmsg "ChatZoo/common/msg"
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

	user := userLogin(conn)
	if user == nil {
		fmt.Println("user 登录失败")
		conn.Close()
		return
	}
	defer user.destroy()
	user.play()
}

func showGameHall() string {
	fmt.Println("welcome to chatZoo, this is game hall. we support some game ")
	fmt.Println("[ > w < ]. [ * v * ]. [ /// - /// ]. [ ` ~ ` ]. [ :) ] ")
	var moduleName string
	fmt.Printf("[%v] [%v]\n", ModuleNameChat, ModuleNameFourOperationCalculate)
	fmt.Println("请输入选择的模块")
	fmt.Scanln(&moduleName)
	return moduleName
}

func userLogin(conn net.Conn) *_User {
	var (
		register = "register"
		login    = "login"
		visitor  = "visitor"
	)

	// Login 登录/注册账号或者游客登录
	var input, openID string
	var isVisitor bool
	fmt.Println("登录账号请输入 login, 注册账号请输入 register, 游客登录请使用 visitor")
	for {
		fmt.Scanln(&input)
		if input == register || input == login || input == visitor {
			break
		}
		fmt.Printf("当前输入为：%v, 此为无效输入,请重新输入\n", input)
	}

	switch input {
	case register:
		fmt.Println("请输入注册的账号名字, 不允许带空格")
		fmt.Scanln(&openID) // 从标准控制中输入,以空格分隔
	case login:
		fmt.Println("请输入登录的账号名字")
		fmt.Scanln(&openID)
	case visitor:
		openID = conn.LocalAddr().String()
		isVisitor = true
	}
	if len(openID) == 0 {
		fmt.Println("openID is empty ")
		return nil
	}
	if err := mmsg.WriteToConn(conn, &mmsg.MsgUserLogin{
		OpenID:    openID,
		IsVisitor: isVisitor,
	}); err != nil {
		fmt.Println("user login send failed ", err)
		return nil
	}
	// todo 最好能等到服务器返回结果才算真正登录成功失败, 现在还不知道怎么做！
	// 假如客户端同时发了两条相同的请求, 服务器也对这两条消息进行了回复, 怎么知道谁是谁的回复。消息有唯一标识吗？
	fmt.Printf("user login success  openID:%v isVisitor:%v\n", openID, isVisitor)
	return NewUser(openID, conn)
}

/*
先登录/游客
进入游戏大厅
选择游戏
退出游戏,回到大厅
退出账号
*/
