package main

import (
	"ChatZoo/common/db"
	"ChatZoo/common/encrypt"
	"ChatZoo/common/login"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 两个不同的进程, 通过redis传递信息
const (
	LoginListen        = "127.0.0.1:7897" // 走配置
	MysqlDataBase      = "chatZoo"        //"happytest" //数据库名字
	MysqlUser          = "root"
	MysqlPwd           = "111111"
	MysqlAddr          = "127.0.0.1:3306"
	MysqlCmdTimeoutSec = 3 * time.Second
	mysqlTableUser     = "User"
	GateListenAddr     = "127.0.0.1:7788"

	redisAddr          = "127.0.0.1:6379"
	redisPassword      = ""
	redisDB            = 8
	redisCmdTimeoutSec = 3 * time.Second // redis 操作超时时间

	HttpName_Register = "/register"
	HttpName_Login    = "/login"
)

type handleHttpFunc func(w http.ResponseWriter, r *http.Request)

type _Srv struct {
	handleHttpFuncList map[string]handleHttpFunc
	storeUtil          db.IStoreUtil
	cacheUtil          db.ICacheUtil
}

func NewSrv() *_Srv {
	ret := &_Srv{
		handleHttpFuncList: make(map[string]handleHttpFunc),
	}
	ret.start()
	return ret
}

func (s *_Srv) start() {
	storeUtil, err := db.NewStoreUtil(MysqlUser, MysqlPwd, MysqlAddr, MysqlDataBase, MysqlCmdTimeoutSec)
	if err != nil {
		panic(err)
	}
	s.storeUtil = storeUtil
	cacheUtil, err := db.NewCacheUtil(redisAddr, redisPassword, redisDB, redisCmdTimeoutSec)
	if err != nil {
		panic(err)
	}
	s.cacheUtil = cacheUtil
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
	req := login.RegisterReq{}
	rsp := login.LoginResp{}
	defer func() {
		WriteJsonRsp(w, rsp) // 结构体也可以
	}()

	err := ParseJsonReq(r, &req) // 得指针
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}
	if len(req.PublicKey) == 0 {
		rsp.Err = fmt.Sprintf("publicKey is empty")
		rsp.Code = Code_Failed
		return
	}
	clientPublicKey, ok := big.NewInt(0).SetString(req.PublicKey, 0)
	if !ok {
		rsp.Err = fmt.Sprintf("publicKey is illegal")
		rsp.Code = Code_Failed
		return
	}

	privateKey, publicKey := encrypt.Pair()
	secret := encrypt.Key(privateKey, clientPublicKey).String()

	err = login.SaveLoginToken(s.cacheUtil, req.PublicKey, secret)
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}

	// 如果先save 后insert 失败了怎么办？ 没关系, redis 等会就过期了, 让玩家重新登

	sqlStr := fmt.Sprintf("insert into %s(ID, Data) values (?,?)  ", mysqlTableUser)
	err = s.storeUtil.InsertData(sqlStr, req.Name, []byte{})
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}

	rsp.PublicKey = publicKey.String()
	rsp.GateAddr = GateListenAddr
	rsp.Code = Code_Success
}

func (s *_Srv) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := login.LoginReq{}
	rsp := login.LoginResp{}
	defer func() {
		WriteJsonRsp(w, rsp) // 结构体也可以
	}()
	err := ParseJsonReq(r, &req) // 得指针
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}
	if len(req.PublicKey) == 0 {
		rsp.Err = fmt.Sprintf("publicKey is empty")
		rsp.Code = Code_Failed
		return
	}
	if !req.IsVisitor { // 非游客检查注册信息
		sqlStr := fmt.Sprintf("select ID from %s where ID = ? ", mysqlTableUser)
		var id string
		err = s.storeUtil.SelectData(sqlStr, &id, req.ID)
		if err != nil {
			rsp.Err = err.Error()
			rsp.Code = Code_Failed
			return
		}
		if id != req.ID {
			rsp.Err = errors.New("id not match").Error()
			rsp.Code = Code_Failed
			return
		}
	}

	clientPublicKey, ok := big.NewInt(0).SetString(req.PublicKey, 0)
	if !ok {
		rsp.Err = fmt.Sprintf("publicKey is illegal")
		rsp.Code = Code_Failed
		return
	}
	privateKey, publicKey := encrypt.Pair()
	secret := encrypt.Key(privateKey, clientPublicKey).String()

	err = login.SaveLoginToken(s.cacheUtil, req.PublicKey, secret)
	if err != nil {
		rsp.Err = err.Error()
		rsp.Code = Code_Failed
		return
	}
	rsp.PublicKey = publicKey.String()
	rsp.GateAddr = GateListenAddr
	rsp.Code = Code_Success
}
