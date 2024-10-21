package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
	"sync"
	"time"
)

type _User struct {
	userID string
	conn   net.Conn
	wg     *sync.WaitGroup
	ticker *time.Ticker
}

func NewUser(conn net.Conn) *_User {
	return &_User{
		conn:   conn,
		wg:     &sync.WaitGroup{},
		ticker: time.NewTicker(30 * time.Second),
	}
}

// login 是否登录成功
func (u *_User) login() bool {
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
		openID = u.conn.LocalAddr().String()
		isVisitor = true
	}
	if len(openID) == 0 {
		fmt.Println("openID is empty ")
		return false
	}
	if err := u.send(&common.MsgUserLogin{
		OpenID:    openID,
		IsVisitor: isVisitor,
	}); err != nil {
		fmt.Println("user login send failed ", err)
		return false
	}
	// todo 最好能等到服务器返回结果才算真正登录成功失败, 现在还不知道怎么做！ 开始头痛起来
	fmt.Printf("user login success  openID:%v isVisitor:%v\n", openID, isVisitor)
	return true
}

func (u *_User) play(moduleFunc ModuleFunc) {
	// 试过wg.Add(1) 放到子协程开始, 但是主协程可能等不到子协程开始就执行wg.Wait(),然后就结束程序了
	u.wg.Add(2)
	go u.receiveLoop()
	go u.sendLoop(moduleFunc)
	u.wg.Wait()
}

func (u *_User) destroy() {
	u.conn.Close()
}

// 震惊！ conn直接复制可行
// sendLoop 持续从标准输入中读取, 并发送给服务器
func (u *_User) sendLoop(moduleFunc ModuleFunc) {
	defer func() { u.wg.Done(); fmt.Println(" receiveFromStdinAndWrite over") }()
	for {
		msg := moduleFunc(u.userID)
		if msg == nil {
			continue
		}
		writeErr := u.send(msg)
		if writeErr != nil {
			fmt.Println("conn.Write failed ", writeErr)
			continue
		}
	}
}

// receiveLoop 持续接收来自服务器的消息
func (u *_User) receiveLoop() {
	defer func() { u.wg.Done(); fmt.Println(" receiveLoop over") }()
	for {
		msg, err := u.receive()
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

func (u *_User) send(msg common.IMessage) error {
	return common.WriteToConn(u.conn, msg)
}

func (u *_User) receive() (common.IMessage, error) {
	msg, err := common.ReadFromConn(u.conn)
	if err != nil {
		fmt.Println("common.ReadFromConn err", err)
		return nil, err
	}
	return msg, nil
}
