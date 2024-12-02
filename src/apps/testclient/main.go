package main

import (
	"ChatZoo/common/cfg"
	"ChatZoo/common/encrypt"
	"ChatZoo/common/hhttp"
	"ChatZoo/common/login"
	mmsg "ChatZoo/common/msg"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
读数据原始方案
stocMsg := make([]byte, 1024)
_, err := conn.Read(stocMsg)

如果超过1024长度, 套接字中的剩余信息会被丢掉; 而且如果是很少的信息, 1024长度浪费
解决方案如下:
如果要提前知道这个包的数据长度好提前make出来接受, 那么需要先发一个包, 包里只放后一个包的数据长度
但是, 长度和数据分两个包发, 如果丢包了则会出问题

更新方案
利用io.ReadFull分段式读取, 一个包里放本次数据长度和数据

_, err = io.ReadFull(c.conn, sizeBuf)
size := uint32(0)
err = common.Decode(sizeBuf, &size)
infoBuf := make([]byte, size)
_, err = io.ReadFull(c.conn, infoBuf)
*/

// todo 一些通用的数据结构需要跨进程使用

const (
	Code_Success = 1
	Code_Failed  = 2
)

func main() {
	client := NewClient()
	isRegister, isLogin, isVisitor, openID, pwd := client.userLoginInput()
	client.isVisitor = isVisitor
	err := client.sendLoginHttp(isRegister, isLogin, isVisitor, openID, pwd)
	if err != nil {
		fmt.Println("send login http err ", err)
		return
	}

	conn, err := net.Dial("tcp", client.gateAddr)
	if err != nil {
		fmt.Printf("net.Dial err:%v gateAddr:%v\n", err, client.gateAddr)
		return
	}

	err = client.sendLoginTcp(conn)
	if err != nil {
		fmt.Println("send login tcp err ", err)
		return
	}
	overTimer := time.NewTicker(5 * time.Second)
Loop:
	for {
		select {
		case <-overTimer.C:
			fmt.Println("wait login rsp over time")
			overTimer.Stop()
			return
		default:
			err = waitLoginResp(conn)
			overTimer.Stop()
			if err != nil {
				fmt.Println("wait login rsp err ", err)
				return
			} else {
				break Loop
			}
		}
	}
	//等待login返回成功再创建
	user := NewUser(client.openID, conn)
	if user == nil {
		fmt.Println("user 登录失败")
		conn.Close()
		return
	}
	defer user.destroy()
	go user.play()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}

func waitLoginResp(conn net.Conn) error {
	msg, err := mmsg.ReadFromConn(conn)
	if err != nil {
		fmt.Println("session handleConnect conn read err ", err)
		return err
	}
	if msg.GetID() != mmsg.MsgID_UserLoginResp {
		return fmt.Errorf("before login, receive other type msg:%v", msg.GetID())
	}
	ret, _ := msg.(*mmsg.MsgUserLoginResp)
	if len(ret.Err) != 0 {
		return errors.New(ret.Err)
	}
	return nil
}

func showGameHall(openid string, moduleList []string) string {
	fmt.Printf("%s welcome to chatZoo, this is game hall. we support some game ", openid)
	fmt.Println("[ > w < ]. [ * v * ]. [ /// - /// ]. [ ` ~ ` ]. [ :) ] ")
	var moduleName string
	for _, v := range moduleList {
		fmt.Printf(" [%v] ", v)
	}
	fmt.Printf("\n")
	fmt.Println("请输入选择的模块")
	fmt.Scanln(&moduleName)
	return moduleName
}

type _Client struct {
	isVisitor              bool
	openID                 string
	gateAddr               string
	clientPublicKey        string
	communicationSecretKey string
	cfg                    cfg.IChatZooServerConfig
}

func NewClient() *_Client {
	return &_Client{
		cfg: cfg.NewServerConfig(),
	}
}

// userLoginInput 玩家登录输入  返回值:是否注册, 是否登录, 是否游客, openID, pwd
func (s *_Client) userLoginInput() (bool, bool, bool, string, string) {
	var (
		register = "register"
		login    = "login"
		visitor  = "visitor"
	)

	// Login 登录/注册账号或者游客登录
	var input, openID, pwd string
	var isRegister, isLogin, isVisitor bool
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
		for {
			fmt.Println("请输入注册的账号名字, 不允许带空格")
			fmt.Scanln(&openID) // 从标准控制中输入,以空格分隔
			if len(openID) != 0 {
				fmt.Println("请输入注册的账号密码, 不允许带空格")
				fmt.Scanln(&pwd) // 从标准控制中输入,以空格分隔,
				isRegister = true
				break
			}
			fmt.Printf("当前账号名字输入为：%v, 此为无效输入\n", openID)
		}
	case login:
		for {
			fmt.Println("请输入登录的账号名字")
			fmt.Scanln(&openID)
			if len(openID) != 0 {
				fmt.Println("请输入登录的账号密码, 不允许带空格")
				fmt.Scanln(&pwd) // 从标准控制中输入,以空格分隔
				isLogin = true
				break
			}
			fmt.Printf("当前账号名字输入为：%v, 此为无效输入\n", openID)
		}
	case visitor: // 游客是随机生成的openID
		isVisitor = true
		openID = encrypt.NewGUID()
	}
	return isRegister, isLogin, isVisitor, openID, pwd
}

func (s *_Client) sendLoginHttp(isRegister, isLogin, isVisitor bool, openID, pwd string) error {
	url := "http://" + s.cfg.GetAppListenAddr(cfg.AppTypeLogin)
	var info []byte
	var err error

	clientPrivateKey, clientPublicKey := encrypt.Pair()

	if isRegister {
		url += "/register"
		req := login.RegisterReq{
			Name:      openID,
			Pwd:       pwd,
			PublicKey: clientPublicKey.String(),
		}
		if info, err = json.Marshal(req); err != nil {
			return err
		}
	}

	if isVisitor {
		url += "/login"
		req := login.LoginReq{
			ID:        openID,
			IsVisitor: true,
			PublicKey: clientPublicKey.String(),
		}
		if info, err = json.Marshal(req); err != nil {
			return err
		}
	}

	if isLogin {
		url += "/login"
		req := login.LoginReq{
			ID:        openID,
			Pwd:       pwd,
			PublicKey: clientPublicKey.String(),
		}
		if info, err = json.Marshal(req); err != nil {
			return err
		}
	}

	var ret []byte
	if ret, err = hhttp.HttpPostBodyWithToken(url, hhttp.MD5(""), info); err != nil {
		return err
	}
	fmt.Println("send login http success ", url)

	resp := login.LoginResp{}
	// 解析参数 使用 json.Unmarshal
	if err = json.Unmarshal(ret, &resp); err != nil {
		return err
	}

	if resp.Code != Code_Success {
		return fmt.Errorf("resp code not success , err:%v", resp.Err)
	}

	serverPublicKey, ok := big.NewInt(0).SetString(resp.PublicKey, 0)
	if !ok {
		return errors.New("can't trans server public key")
	}

	s.gateAddr = resp.GateAddr
	s.clientPublicKey = clientPublicKey.String()
	s.communicationSecretKey = encrypt.Key(clientPrivateKey, serverPublicKey).String()
	s.openID = openID
	return nil
}

func (s *_Client) sendLoginTcp(conn net.Conn) error {
	err := mmsg.WriteToConn(conn, &mmsg.MsgUserLogin{
		OpenID:    s.openID,
		IsVisitor: s.isVisitor,
		PublicKey: s.clientPublicKey,
	})
	return err
}
