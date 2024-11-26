package main

import (
	"bufio"
	"fmt"
	"os"
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

func (u *_Module) Chain() (string, map[int]interface{}) {
	fmt.Println("this is chain module, 输入指令后回车换行结束 ")
	fmt.Printf("创建空房间请输入 [1 房间名字] 示例：1 myroomname \n")
	fmt.Printf("加入已有房间请输入 [2 房间名字 房间类型] 示例：2 joinroomname  \n")
	fmt.Printf("查看推荐房间请输入 [3] 示例：3 \n")

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
	var ret string
	switch cmds[0] {
	case CreateRoom: // 创建房间
		methodName = "CRPC_CreateRoom"
		args[0] = RoomType_Chain // 房间类型
		args[1] = cmds[1]        // 房间名字
		args[2] = cmds[2]        // 房间最大人数
	case JoinRoom:
		methodName = "CRPC_JoinRoom"
		args[0] = cmds[1]        // 房间名字
		args[1] = RoomType_Chain // 房间类型
	case RecommendRoom:
		methodName = "CRPC_GetRecommendRoom"
		args[0] = RoomType_Chain // 房间类型
	case RoomGuessReady:
		methodName = "CRPC_GuessRoomReady"
	case RoomGuessPlay:
		methodName = "CRPC_GuessRoomSendMsg"
	default:
		fmt.Println(" 参数不对 ")
		return "", nil
	}
	ret = cmds[1]
	return methodName, ret
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