package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
)

type _User struct {
	stopCh chan struct{}
	wg     *sync.WaitGroup
	common.IEntityInfo
	module *_Module
}

type _Module struct {
	*ChainModule
	*ChatModule
	*CalculateModule
}

func NewModule(user *_User) *_Module {
	return &_Module{
		NewChainModule(user),
		NewChatModule(user),
		NewCalculateModule(user),
	}
}

func NewUser(entityID string, conn net.Conn) *_User {
	user := &_User{
		wg:          &sync.WaitGroup{},
		IEntityInfo: common.NewEntityInfo(entityID, conn),
	}
	user.module = NewModule(user)
	user.SetRpc(user)
	go user.Loop()
	return user
}

func (u *_User) play() {
	// 试过wg.Add(1) 放到子协程开始, 但是主协程可能等不到子协程开始就执行wg.Wait(),然后就结束程序了
	// 只有一个协程犯得着用wg吗？这里是不是还考虑其他的协程处理别的逻辑？
	u.wg.Add(1)
	go u.sendLoop()
	u.wg.Wait()
	u.stopCh <- struct{}{}

	fmt.Println("play over")
}

func (u *_User) destroy() {
	u.GetNetConn().Close()
}

// 震惊！ conn直接复制可行
// sendLoop 持续从标准输入中读取, 并发送给服务器
func (u *_User) sendLoop() {
	defer func() { u.wg.Done(); fmt.Println("sendLoop over") }()

	moduleMethod := u.getPlayerInputModuleName()
	moduleMethod.Call([]reflect.Value{})
}

// //// 客户端表现模块
// getModuleName 获取玩家模块输入
func (u *_User) getPlayerInputModuleName() reflect.Value {
	moduleList := make([]string, 0)
	t := reflect.TypeOf(u.module)
	for i := 0; i < t.NumMethod(); i++ {
		moduleList = append(moduleList, strings.ToLower(t.Method(i).Name))
	}
	var method reflect.Value
	for {
		v := reflect.ValueOf(u.module)
		cmd := showGameHall(u.GetEntityID(), moduleList)
		cmd = strings.ToLower(cmd)
		cmd = fmt.Sprintf("%s%s", strings.ToUpper(string(cmd[0])), cmd[1:])
		method = v.MethodByName(cmd)

		if method.Kind() == reflect.Func && !method.IsNil() { // 进入下一流程
			return method
		}
		fmt.Println("模块名不对, 请重新输入 ", cmd)
	}
}
