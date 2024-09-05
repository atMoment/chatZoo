package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type _App struct {
	exitChan  chan struct{}
	appCtx    context.Context
	appCancel context.CancelFunc
}

func NewApp() *_App {
	app := &_App{
		exitChan: make(chan struct{}, 1), // 长度为1或者不为1都一样
	}
	app.appCtx, app.appCancel = context.WithCancel(context.Background())
	return app
}

func (a *_App) run() {
	ln, err := net.Listen("tcp", "127.0.0.1:7788") // 必须外部配置
	if err != nil {
		fmt.Println("net listen err ", err)
		return
	}
	go a.acceptHandler(ln)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	ln.Close()
	<-a.exitChan
}

// acceptHandler ln 可以直接复制!
func (a *_App) acceptHandler(ln net.Listener) {
	for {
		conn, acceptErr := ln.Accept()
		if acceptErr != nil {
			fmt.Println("ln accept err ", acceptErr)
			a.appCancel()
			return
		}
		session := NewSession(a.appCtx, conn)
		go session.proc() // todo 怎么等他们都结束才完全退出
	}
}

/* lcc 写的错误写法
func (a *_App) accept() {
	fmt.Println("app run")
	ln, err := net.Listen("tcp", "127.0.0.1:7788") // 必须外部配置
	if err != nil {
		fmt.Println("net listen err ", err)
		return
	}
	defer ln.Close()
	for {
		select {
		case <-a.exitChan:
			fmt.Println("app exit")
			a.appCancel()
			return
		default:
			conn, acceptErr := ln.Accept() // 主协程退出以后阻塞在这里, 也不会走上面的case, 也不会走 ln.Close(), 协程泄露
			if acceptErr != nil {
				fmt.Println("ln accept err ", acceptErr)
				return
			}
			session := NewSession(a.appCtx, conn)
			go session.proc()
		}
	}
}
*/
