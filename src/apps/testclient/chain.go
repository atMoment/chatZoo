package main

import (
	"ChatZoo/common"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 接收服务器回调函数
const (
	_ = iota
	CMDID_CreateRoom
	CMDID_JoinRoom
	CMDID_GetRecommendRoom
)

const (
	RoomType_None = iota
	RoomType_Chat
	RoomType_Chain
)

const (
	ChainStep_NotBegin = iota
	ChainStep_JoinRoom
	ChainStep_Ready
	ChainStep_GameBegin
)

type ChainModule struct {
	nextStage int
	key       string
	user      *_User
}

func NewChainModule(user *_User) *ChainModule {
	return &ChainModule{
		user: user,
	}
}

func showRoomInfo() {
	fmt.Println("this is chain module, 输入指令后回车换行结束 ")
	fmt.Printf("创建空房间请输入 [1 房间名字 房间最大人数] 示例：1 roomname 2 \n")
	fmt.Printf("加入已有房间请输入 [2 房间名字] 示例：2 roomname  \n")
	fmt.Printf("查看推荐房间请输入 [3] 示例：3 \n")
}

func showChainReady() {
	fmt.Println("this is chain module, 输入指令后回车换行结束 ")
	fmt.Printf("您已加入接龙房间, 游戏即将开始, 准备好请输入 5 \n")
}

func (u *ChainModule) Chain() {
	for {
		u.selectRoom()
		if u.nextStage == ChainStep_Ready {
			break
		}
	}
}

// SetNextStage 设置阶段
func (u *ChainModule) setNextStage(methodName string) error {
	switch methodName {
	case "CRPC_GetRecommendRoom":
		u.nextStage = ChainStep_NotBegin
	case "CRPC_CreateRoom", "CRPC_JoinRoom":
		u.nextStage = ChainStep_JoinRoom
	case "CRPC_ChainRoomReady":
		u.nextStage = ChainStep_Ready
	case "CRPC_ChainRoomSendMsg":
		u.nextStage = ChainStep_GameBegin
	default:
		return fmt.Errorf("unsupported methodName%v, can't trans stage", methodName)
	}
	return nil
}

func (u *ChainModule) selectRoom() {
	for {
		showRoomInfo()
		err, cmds := waitPlayerInput(3)
		if err != nil {
			fmt.Printf("selectRoom input err:%v\n", err)
			continue
		}
		switch cmds[0] {
		case CreateRoom: // 创建房间
			roomLimit, AtoiErr := strconv.Atoi(cmds[2])
			if AtoiErr != nil {
				fmt.Println(ModuleNameChat, " create room limit not num ", cmds[2])
				continue
			}
			u.createRoom(cmds[1], roomLimit)
		case JoinRoom:
			u.joinRoom(cmds[1])
		case RecommendRoom:
			u.recommendRoom()
		}
	}
}

func (u *ChainModule) readyRoom() {
	// CRPC_ChainRoomReady
}

func (u *ChainModule) createRoom(roomid string, limit int) error {
	methodName := "CRPC_CreateRoom"
	ret := <-u.user.GetRpc().SendReq(methodName, RoomType_Chain, roomid, limit)
	return analyseRpcReqRet(ret)
}

func (u *ChainModule) joinRoom(roomid string) error {
	methodName := "CRPC_JoinRoom"
	ret := <-u.user.GetRpc().SendReq(methodName, roomid)
	return analyseRpcReqRet(ret)
}

func (u *ChainModule) recommendRoom() error {
	methodName := "CRPC_GetRecommendRoom"
	ret := <-u.user.GetRpc().SendReq(methodName)
	err := analyseRpcReqRet(ret)
	if err != nil {
		return err
	}
	recommendList, ok := ret.Rets[1].([]string)
	if !ok {
		return fmt.Errorf("not []string")
	}
	fmt.Println("recommend room: ", recommendList)
	return nil
}

func (u *ChainModule) chainStart(key string) {
	for {
		fmt.Println("guess game turn begin, 请以此字符串的结尾作为开头组成成语/俗语 ", key)
		fmt.Printf("请输入 接龙内容 \n")

		err, cmds := waitPlayerInput(1)
		if err != nil {
			fmt.Printf("chainStart input err:%v\n", err)
			continue
		}

		methodName := "CRPC_ChainRoomSendMsg"
		ret := <-u.user.GetRpc().SendReq(methodName, cmds[0])
		err = analyseRpcReqRet(ret)
		if err != nil {
			fmt.Printf("chainStart methodName:%v req err:%v  \n", methodName, err)
			continue
		}
	}
}

func (r *_User) SRPC_ChainGameTurnBegin(key string) {
	r.module.ChainModule.chainStart(key)
}

func (r *_User) SPRC_ChainGameOver() {
	fmt.Println("guess game over")
}

func analyseRpcReqRet(ret *common.CallRet) error {
	if ret.Err != nil {
		return ret.Err
	}
	if len(ret.Rets) == 0 {
		return errors.New("ret.Rets length is 0")
	}
	errStr, ok := ret.Rets[0].(string)
	if !ok {
		return errors.New("ret.Rets[0] not string")
	}
	if errStr != "success" {
		return fmt.Errorf("ret err:%v", errStr)
	}
	return nil
}

func waitPlayerInput(cmdLength int) (error, []string) {
	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		return inputErr, nil
	}

	// 把字符串中的\r\n筛选出来
	cmds := filterSeparator(input)
	if len(cmds) != cmdLength {
		return fmt.Errorf("输入长度不匹配"), nil
	}
	return nil, cmds
}

// filterSeparator 过滤分隔符"1 room \r\n" 转化为 "1 room", 再转化为[]string{"1", "room"}
func filterSeparator(input string) []string {
	separator := "\r\n"
	words := ""
	for _, v := range strings.Split(input, separator) {
		if v == separator {
			continue
		}
		words += v
	}
	return strings.Split(words, " ")
}
