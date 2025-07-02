package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ModuleNameChat = "Chat"
	CreateRoom     = "1"
	JoinRoom       = "2"
	RecommendRoom  = "3"
	RoomGuessReady = "4"
	ChatRoom       = "101"
)

type ChatModule struct {
	user *_User
}

func NewChatModule(user *_User) *ChatModule {
	return &ChatModule{
		user: user,
	}
}

func (u *ChatModule) String() string {
	return ModuleNameChat
}

func (u *ChatModule) Chat() (string, string) {
	showRoomInfo(u.String())

	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println(ModuleNameChat, " os.stdin read err ", inputErr)
		return "", ""
	}

	// 把字符串中的\r\n筛选出来
	var builder strings.Builder
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		builder.WriteString(v)
	}

	cmds := strings.Split(builder.String(), " ")
	if len(cmds) == 0 {
		fmt.Println(ModuleNameChat, " 无有效输入, 长度不对 ", len(cmds))
		return "", ""
	}
	var methodName string
	var ret string
	switch cmds[0] {
	case CreateRoom: // 创建房间
		fmt.Println(ModuleNameChat, "暂不支持, sorry")
		return "", ""
	case JoinRoom:
		methodName = "JoinRoom"
	case RecommendRoom:
		fmt.Println(ModuleNameChat, "开发中...")
		return "", ""
	case ChatRoom:
		methodName = "ChatRoom"
	case RoomGuessReady:
		methodName = "GuessRoomReady"
	default:
		fmt.Println(ModuleNameChat, " 参数不对 ")
		return "", ""
	}
	ret = cmds[1]
	return methodName, ret
}

//////////////     接收服务器回调函数    //////////////

func (r *_User) Notify_SToCMessage(msg string) {
	fmt.Println("receive msg ", msg)
}

/*
流程和剩下的安排
1. 登录/游客账号
2. 游戏大厅:
   聊天类： 你画我猜/聊天室/成语接龙/动物园里有什么/加字减字组成新的一句话/用其他语言描述此物品
   竞技类:  限时加减乘除24/限时2048, 排行榜展示
   复杂类： 来种地吧(具体需求未想好,可以当某人的孩子,出生就继承祖宅)
3. 选择一个游戏进入
4. 执行对应的游戏逻辑
5. 登出-自动关闭客户端 手动杀端
*/
