package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
	"reflect"
	"sync"
)

type _User struct {
	wg *sync.WaitGroup
	common.IEntityInfo
	module *_Module
}

type _Module struct{}

func NewUser(entityID string, conn net.Conn) *_User {
	user := &_User{
		wg:          &sync.WaitGroup{},
		IEntityInfo: common.NewEntityInfo(entityID, conn),
		module:      &_Module{},
	}
	user.SetRpc(user)
	return user
}

func (u *_User) play() {
	// 试过wg.Add(1) 放到子协程开始, 但是主协程可能等不到子协程开始就执行wg.Wait(),然后就结束程序了
	u.wg.Add(2)
	go u.receiveLoop()
	go u.sendLoop()
	u.wg.Wait()
}

func (u *_User) destroy() {
	u.GetNetConn().Close()
}

// 震惊！ conn直接复制可行
// sendLoop 持续从标准输入中读取, 并发送给服务器
func (u *_User) sendLoop() {
	defer func() { u.wg.Done(); fmt.Println(" receiveFromStdinAndWrite over") }()
	moduleList := make([]string, 0)
	t := reflect.TypeOf(u.module)
	for i := 0; i < t.NumMethod(); i++ {
		moduleList = append(moduleList, t.Method(i).Name)
	}

	// 玩家输入指令
	var method reflect.Value
	for {
		v := reflect.ValueOf(u.module)
		cmd := showGameHall(u.GetEntityID(), moduleList)
		method = v.MethodByName(cmd)

		if method.Kind() == reflect.Func && !method.IsNil() { // 进入下一流程
			break
		}
		fmt.Println("模块名不对, 请重新输入 ", cmd)
	}

	for {
		args := method.Call([]reflect.Value{})
		if len(args) < 2 {
			fmt.Println("模块失败, 参数数量错误 ")
			continue
		}
		// reflect 全是卡点。 把 reflect.Value 转化成 string
		methodName := args[0].Interface().(string)
		if len(methodName) == 0 {
			fmt.Println("methodName empty ")
			continue
		}
		// 根据消息类型ID反解析参数, 解析出来怎么确定参数顺序
		// map[index]interface

		// 将 reflect.Value 还原成实际的类型
		in := make([]interface{}, len(args)-1)
		var j int
		for i := 1; i < len(args); i++ {
			in[j] = args[i].Interface()
		}
		// 想要声明一个函数, 函数的返回值是 ...interface, 方便传入 SendReq中。 返回值是真实的类型而不是 真实类型转化后的interface类型
		ret := <-u.GetRpc().SendReq(methodName, in...)
		if ret.Err != nil {
			fmt.Println("ret.Err  ", ret.Err)
			continue
		}
		fmt.Println("get ret ", ret.Rets)
	}
}

// receiveLoop 持续接收来自服务器的消息
func (u *_User) receiveLoop() {
	defer func() { u.wg.Done(); fmt.Println(" receiveLoop over") }()
	for {
		err := u.GetRpc().ReceiveConn()
		if err != nil {
			fmt.Println("common.ReadFromConn err", err)
			return
		}
	}
}
