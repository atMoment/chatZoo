package main

import (
	"bufio"
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

type _ChainModule struct {
	nextStage int
}

func showRoomInfo() {
	fmt.Println("this is chain module, 输入指令后回车换行结束 ")
	fmt.Printf("创建空房间请输入 [1 房间名字 房间最大人数] 示例：1 roomname 2 \n")
	fmt.Printf("加入已有房间请输入 [2 房间名字] 示例：2 roomname  \n")
	fmt.Printf("查看推荐房间请输入 [3] 示例：3 \n")
}

func showChainInfo() {
	fmt.Println("this is chain module, 输入指令后回车换行结束 ")
	fmt.Printf("您已加入接龙房间, 准备好请输入 5 \n")
}

// Chain 想要返回任意数量,任意类型的参数. []interface不行, map[Any]interface{} 就行
func (u *_ChainModule) Chain() (string, map[int]interface{}) {
	switch u.nextStage {
	case ChainStep_NotBegin:
		showRoomInfo()
	case ChainStep_JoinRoom:
		showChainInfo()
	}
	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println(ModuleNameChat, " os.stdin read err ", inputErr)
		return "", nil
	}

	// 把字符串中的\r\n筛选出来
	cmds := filterSeparator(input)
	if len(cmds) == 0 {
		fmt.Println(ModuleNameChat, " 无有效输入, 长度不对 ", len(cmds))
		return "", nil
	}
	args := make(map[int]interface{})
	var methodName string
	switch cmds[0] {
	case CreateRoom: // 创建房间
		methodName = "CRPC_CreateRoom"
		args[0] = RoomType_Chain // 房间类型
		args[1] = cmds[1]        // 房间id
		roomLimit, err := strconv.Atoi(cmds[2])
		if err != nil {
			fmt.Println(ModuleNameChat, " create room limit not num ", cmds[2])
			return "", nil
		}
		args[2] = roomLimit // 房间最大人数
	case JoinRoom:
		methodName = "CRPC_JoinRoom"
		args[0] = cmds[1] // 房间名字
	case RecommendRoom:
		methodName = "CRPC_GetRecommendRoom"
	case RoomGuessReady:
		methodName = "CRPC_ChainRoomReady"
	case RoomGuessPlay:
		methodName = "CRPC_ChainRoomSendMsg"
		args[0] = cmds[1] // 接龙内容
	default:
		fmt.Println(" 参数不对 ")
		return "", nil
	}
	return methodName, args
}

func (r *_User) ChainGameTurnBegin(firstKey string) {
	fmt.Println("guess game turn begin, 请以此字符串的结尾作为开头组成成语/俗语 ", firstKey)
}

func (r *_User) ChainGameOver() {
	fmt.Println("guess game over")
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
