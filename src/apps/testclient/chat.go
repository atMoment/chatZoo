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
	user       *_User
	joinRoomID string // 加入的房间ID
}

func NewChatModule(user *_User) *ChatModule {
	return &ChatModule{
		user: user,
	}
}

func (u *ChatModule) String() string {
	return ModuleNameChat
}

func (u *ChatModule) Chat() {
	for {
		err := u.chat()
		if err != nil {
			fmt.Println("chat err:%v", err)
			return
		}
	}

}

func (u *ChatModule) chat() error {
	if u.joinRoomID == "" {
		showRoomInfo(u.String())
	} else {
		showChat()
	}

	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println(ModuleNameChat, " os.stdin read err ", inputErr)
		return inputErr
	}

	cmds := strings.Split(dealInput4(input), " ")
	if len(cmds) == 0 {
		fmt.Println(ModuleNameChat, " 无有效输入, 长度不对 ", len(cmds))
		return fmt.Errorf("input is empty")
	}
	var err error
	switch cmds[0] {
	case CreateRoom: // 创建房间
		if len(cmds) < 2 {
			return fmt.Errorf("joinroom, len(cmds) < 2")
		}
		if len(cmds[1]) == 0 {
			return fmt.Errorf("joinroom, roomid is empty")
		}
		err = u.createRoom(cmds[1], 100)
	case JoinRoom:
		if len(cmds) < 2 {
			return fmt.Errorf("joinroom, len(cmds) < 2")
		}
		if len(cmds[1]) == 0 {
			return fmt.Errorf("joinroom, roomid is empty")
		}
		err = u.joinRoom(cmds[1])
	case RecommendRoom:
		err = u.recommendRoom()
	case ChatRoom:
		if len(cmds) < 2 {
			return fmt.Errorf("joinroom, len(cmds) < 2")
		}
		if len(cmds[1]) == 0 {
			return fmt.Errorf("joinroom, roomid is empty")
		}
		if len(u.joinRoomID) == 0 {
			return fmt.Errorf("未加入房间")
		}
		u.chatRoom(cmds[1])
	default:
		fmt.Println(ModuleNameChat, " 参数不对 ")
	}
	return err
}

func (u *ChatModule) createRoom(roomid string, limit int) error {
	methodName := "CRPC_CreateRoom"
	ret := <-u.user.GetRpc().SendReq(methodName, RoomType_Chat, roomid, limit)
	err := analyseRpcReqRet(ret)
	if err != nil {
		return err
	}
	u.joinRoomID = roomid
	fmt.Printf("create room %v success\n", roomid)
	return nil
}

func (u *ChatModule) joinRoom(roomid string) error {
	methodName := "CRPC_JoinRoom"
	ret := <-u.user.GetRpc().SendReq(methodName, roomid)
	err := analyseRpcReqRet(ret)
	if err != nil {
		return err
	}
	u.joinRoomID = roomid
	fmt.Printf("join room %v success\n", roomid)
	return nil
}

func (u *ChatModule) recommendRoom() error {
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

func (u *ChatModule) chatRoom(msg string) error {
	methodName := "CRPC_Chat"
	ret := <-u.user.GetRpc().SendReq(methodName, u.joinRoomID, msg)
	err := analyseRpcReqRet(ret)
	if err != nil {
		return err
	}
	return nil
}

func showChat() {
	fmt.Printf("您已加入房间, 聊天请输入 [101 聊天内容]  \n")
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
