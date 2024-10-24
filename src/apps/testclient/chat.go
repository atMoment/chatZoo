package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	ModuleNameChat = "Chat"
	CreateRoom     = "1"
	JoinRoom       = "2"
	RecommendRoom  = "3"
	ChatRoom       = "4"
)

func (u *_User) Chat() (string, []interface{}) {
	fmt.Println("欢迎来到chat zoo, 输入指令后回车换行结束 ")
	fmt.Printf("创建空房间请输入 [1 房间名字] 示例：1 myroomname \n")
	fmt.Printf("加入已有房间请输入 [2 房间名字] 示例：2 joinroomname \n")
	fmt.Printf("查看推荐房间请输入 [3] 示例：3 \n")
	fmt.Printf("聊天请输入 [4 房间名字 聊天内容] 示例：4 roonname content \n")

	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println(ModuleNameChat, " os.stdin read err ", inputErr)
		return "", nil
	}

	// 把字符串中的\r\n筛选出来
	words := ""
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		words += v // todo 好像这种写法很消耗,有新的写法
	}

	cmds := strings.Split(words, " ")

	if len(cmds) == 0 {
		fmt.Println(ModuleNameChat, " 无有效输入 ")
		return "", nil
	}
	var methodName string
	out := make([]interface{}, 0)
	switch cmds[0] {
	case CreateRoom: // 创建房间
		if len(cmds[1]) == 0 {
			fmt.Println(ModuleNameChat, "创建房间需要输入房间名字")
			return "", nil
		} else {
			methodName = "CreateRoom"
			out = append(out, reflect.ValueOf(cmds[1]))
		}
	case JoinRoom:
		if len(cmds[1]) == 0 {
			fmt.Println(ModuleNameChat, "加入房间需要输入房间名字")
			return "", nil
		} else {
			methodName = "JoinRoom"
			out = append(out, reflect.ValueOf(cmds[1]))
		}
	case RecommendRoom:
		fmt.Println(ModuleNameChat, "开发中...")
		return "", nil
	case ChatRoom:
		if len(cmds) != 3 {
			fmt.Println(ModuleNameChat, "chatroom 参数不对")
			return "", nil
		} else {
			methodName = "ChatRoom"
			out = append(out, reflect.ValueOf(cmds[1]), reflect.ValueOf(cmds[2]))
		}
	default:
		fmt.Println(ModuleNameChat, " 参数不对 ")
		return "", nil
	}
	return methodName, out
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
