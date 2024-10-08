package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type _App struct {
	exitChan  chan struct{}
	appCtx    context.Context
	appCancel context.CancelFunc
	wg        *sync.WaitGroup
}

func NewApp() *_App {
	app := &_App{
		exitChan: make(chan struct{}, 1), // 长度为1或者不为1都一样吗？
		wg:       &sync.WaitGroup{},
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
	// 等待结束信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	ln.Close()   // 发通知让子协程acceptHandler 退出
	<-a.exitChan // 阻塞住, 为的是等待所有子协程都退出主协程才退出
}

// acceptHandler ln 可以直接复制!
func (a *_App) acceptHandler(ln net.Listener) {
	for {
		conn, acceptErr := ln.Accept()
		if acceptErr != nil {
			fmt.Println("groutine: ln accept err ", acceptErr)
			a.appCancel()            // 1 通知 多, 通知多个子协程退出
			a.wg.Wait()              // 1 等待 多(多通知1), sync.WaitGroup等他们都结束才完全退出
			a.exitChan <- struct{}{} // 给父协程发信号, 让老父亲别等了, 我也要销毁了
			return
		}
		session := NewSession(a.appCtx, conn, a.wg)
		go session.proc()
	}
}

/* lcc 写的错误写法
func (a *_App) acceptHandler() {
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
			// 还没来得及走下一个循环到 走到exitChan
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
