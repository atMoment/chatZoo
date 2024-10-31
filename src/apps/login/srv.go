package main

import (
	"ChatZoo/common/cfg"
	"ChatZoo/common/db"
	"ChatZoo/common/encrypt"
	"ChatZoo/common/hhttp"
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

/*
CREATE TABLE `User`  (

`ID` char(16) NOT NULL,  // 定长,不管中文还是英文,只能装16个

`Data` mediumblob NULL,  // 用3字节存储, 1字节=8bit, 用24bit存储, 2^24 = 16M, 能存16M数据

PRIMARY KEY (`ID`) USING BTREE

) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = Dynamic;

*/

// todo 最好还是有一个 id 作为唯一ID, 假如不取名呢？ 假如允许名字重复呢？
// 假如不取名也允许做一些操作呢？

// 两个不同的进程, 通过redis传递信息
const (
	mysqlTableUser = "User"

	HttpName_Register = "/register"
	HttpName_Login    = "/login"

	Code_Success = 1 // todo 应该有跨进程使用的统一的数据结构 protobuf?
	Code_Failed  = 2
)

type handleHttpFunc func(w http.ResponseWriter, r *http.Request)

type _Srv struct {
	handleHttpFuncList map[string]handleHttpFunc
	storeUtil          db.IStoreUtil
	cacheUtil          db.ICacheUtil
	cfg                cfg.IChatZooServerConfig
	typ                string // 什么类型的进程
}

func NewSrv() *_Srv {
	ret := &_Srv{
		handleHttpFuncList: make(map[string]handleHttpFunc),
		typ:                cfg.AppTypeLogin,
		cfg:                cfg.NewServerConfig(),
	}
	ret.start()
	return ret
}

func (s *_Srv) start() {
	storeUtil, err := db.NewStoreUtil(s.cfg.GetMysqlCfg().MysqlUser, s.cfg.GetMysqlCfg().MysqlPwd, s.cfg.GetMysqlCfg().MysqlAddr, s.cfg.GetMysqlCfg().MysqlDataBase, time.Duration(s.cfg.GetMysqlCfg().MysqlCmdTimeoutSec)*time.Second)
	if err != nil {
		panic(err)
	}
	s.storeUtil = storeUtil
	cacheUtil, err := db.NewCacheUtil(s.cfg.GetRedisCfg().RedisAddr, s.cfg.GetRedisCfg().RedisPwd, s.cfg.GetRedisCfg().RedisDB, time.Duration(s.cfg.GetRedisCfg().RedisCmdTimeoutSec)*time.Second)
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
		if err := http.ListenAndServe(s.cfg.GetAppListenAddr(s.typ), nil); err != nil {
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
		fmt.Printf("handleRegister isSuccess:%v rsp:%+v \n", rsp.Code == Code_Success, rsp)
		hhttp.WriteJsonRsp(w, rsp) // 结构体也可以
	}()

	err := hhttp.ParseJsonReq(r, &req) // 得指针
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
	if len(req.Name) > 16 {
		rsp.Err = fmt.Sprintf("name is toolong")
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
	rsp.GateAddr = s.cfg.GetAppListenAddr(cfg.AppTypeGate)
	rsp.Code = Code_Success
}

func (s *_Srv) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := login.LoginReq{}
	rsp := login.LoginResp{}
	defer func() {
		fmt.Printf("handleLogin isSuccess:%v rsp:%+v \n", rsp.Code == Code_Success, rsp)
		hhttp.WriteJsonRsp(w, rsp) // 结构体也可以
	}()
	err := hhttp.ParseJsonReq(r, &req) // 得指针
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
	rsp.GateAddr = s.cfg.GetAppListenAddr(cfg.AppTypeGate)
	rsp.Code = Code_Success
}
