package main

import (
	"fmt"
	"net/http"
	"sync"
)

/*
http 服务器, 职责如下
接收客户端login请求, 交换公钥, 存储数据库, 发送 gate端口
问题：login时存数据库会不会被爆冲？即一堆http请求来注册 (正常游戏有sdk验证,但是我没有)
*/

const (
	LoginListen = "127.0.0.1:7897" // 走配置
)

type handleHttpFunc func(w http.ResponseWriter, r *http.Request)

var DefaultHandleHttpFuncList map[string]handleHttpFunc

func init() {
	DefaultHandleHttpFuncList["/register"] = handleRegister
	DefaultHandleHttpFuncList["/login"] = handleLogin
}

func main() {
	if len(DefaultHandleHttpFuncList) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	for key, val := range DefaultHandleHttpFuncList {
		http.HandleFunc(key, val)
	}
	go func() {
		defer wg.Done()
		if err := http.ListenAndServe(LoginListen, nil); err != nil {
			fmt.Println("login listen err: ", err)
		}
	}()
	wg.Wait() // 也应该监听 kill信号
}

func handleRegister(w http.ResponseWriter, r *http.Request) {

}

func handleLogin(w http.ResponseWriter, r *http.Request) {

}
