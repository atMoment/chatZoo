package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	LoginListen       = "127.0.0.1:7897" // 走配置
	HttpName_Register = "/register"
	HttpName_Login    = "/login"
)

type handleHttpFunc func(w http.ResponseWriter, r *http.Request)

type _Srv struct {
	handleHttpFuncList map[string]handleHttpFunc
}

func NewSrv() *_Srv {
	ret := &_Srv{
		handleHttpFuncList: make(map[string]handleHttpFunc),
	}
	//StoreUtil, err = db.NewStoreUtil(MysqlUser, MysqlPwd, MysqlAddr, MysqlDataBase, MysqlCmdTimeoutSec)

	ret.start()
	return ret
}

func (s *_Srv) start() {
	s.handleHttpFuncList[HttpName_Register] = s.handleRegister
	s.handleHttpFuncList[HttpName_Login] = s.handleLogin
}

func (s *_Srv) run() {
	over := make(chan struct{})
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-c
		close(over)
	}()

	for key, val := range s.handleHttpFuncList {
		http.HandleFunc(key, val)
	}
	go func() {
		if err := http.ListenAndServe(LoginListen, nil); err != nil {
			close(over)
			fmt.Println("login listen err: ", err)
		}
	}()
	<-over
	fmt.Println("srv run over")
}

func (s *_Srv) handleRegister(w http.ResponseWriter, r *http.Request) {
	req := RegisterReq{}
	rsp := LoginResp{}
	defer func() {
		WriteJsonRsp(w, rsp) // 结构体也可以
	}()

	err := ParseJsonReq(r, &req) // 得指针
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}

}

func (s *_Srv) handleLogin(w http.ResponseWriter, r *http.Request) {

}
