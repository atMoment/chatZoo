package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
	"reflect"
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
}

func NewModule(user *_User) *_Module {
	return &_Module{
		NewChainModule(user),
	}
}

func NewUser(entityID string, conn net.Conn) *_User {
	user := &_User{
		wg:          &sync.WaitGroup{},
		IEntityInfo: common.NewEntityInfo(entityID, conn),
	}
	user.module = NewModule(user)
	user.SetRpc(user)
	return user
}

func (u *_User) play() {
	// 试过wg.Add(1) 放到子协程开始, 但是主协程可能等不到子协程开始就执行wg.Wait(),然后就结束程序了
	u.wg.Add(2)
	go u.receiveLoop()
	go u.dealLoop()
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

// receiveLoop 持续接收来自服务器的消息
func (u *_User) receiveLoop() {
	defer func() { u.wg.Done(); fmt.Println("receiveLoop over") }()
	for {
		err := u.GetRpc().ReceiveConn()
		if err != nil {
			fmt.Println("common.ReadFromConn err", err)
			return
		}
	}
}

func (u *_User) dealLoop() {
	for {
		select {
		case <-u.stopCh:
			return
		default:
			u.GetRpcQueue().Pop() // todo 需要等待处理完
		}
	}
}

// //// 客户端表现模块
// getModuleName 获取玩家模块输入
func (u *_User) getPlayerInputModuleName() reflect.Value {
	moduleList := make([]string, 0)
	t := reflect.TypeOf(u.module)
	for i := 0; i < t.NumMethod(); i++ {
		moduleList = append(moduleList, t.Method(i).Name)
	}
	var method reflect.Value
	for {
		v := reflect.ValueOf(u.module)
		cmd := showGameHall(u.GetEntityID(), moduleList)
		method = v.MethodByName(cmd)

		if method.Kind() == reflect.Func && !method.IsNil() { // 进入下一流程
			return method
		}
		fmt.Println("模块名不对, 请重新输入 ", cmd)
	}
}

/*
想要的结果, 显示推荐房间、创建房间、加入房间
只允许玩家发送1,2,3
玩家选择了2 和 3, 发到服务器成功了, 进入到下一个阶段
                发到服务器失败了, 再次显示这个

下一个阶段, 提示玩家准备
只允许玩家发送5, 发到服务器成功了, 进入到下一个阶段
下一个阶段, 等待通知开始后, 输入 6 内容
*/
