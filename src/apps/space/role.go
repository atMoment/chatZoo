package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
)

// 设计上一个玩家只能加入一个房间
type _User struct {
	*common.EntityInfo
	joinRoomID   string
	joinRoomType int
}

func NewUser(entityID string, conn net.Conn) (*_User, error) {
	user := &_User{EntityInfo: common.NewEntityInfo(entityID, conn)}
	user.SetRpc(user)

	return user, nil
}

func (r *_User) CRPC_JoinRoom(typ int, roomID string, limit int) string {
	room, err := roomMgr.AddOrGetEntity(roomID, typ, limit)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return "failed"
	}
	room.JoinRoom(r.GetEntityID())
	r.joinRoomID = roomID
	r.joinRoomType = typ
	fmt.Printf("Room JoinRoom  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CRPC_CreateRoom(typ int, roomID string, limit int) string {

}

func (r *_User) QuitRoom(roomID string) {

}

func (r *_User) GuessRoomReady() string {
	if len(r.joinRoomID) == 0 {
		return "not in room"
	}
	ir, err := getGuessRoom(r.joinRoomID)
	if err != nil {
		fmt.Printf("GuessRoomReady get guess room err:%v\n", err)
		return err.Error()
	}
	ir.Ready(r.GetEntityID())
	fmt.Printf("GuessRoomReady userid:%v, roomid:%v \n ", r.GetEntityID(), r.joinRoomID)
	return "ChatRoom success"
}

func (r *_User) GuessRoomSendMsg(content string) string {
	if len(r.joinRoomID) == 0 {
		return "not in room"
	}
	ir, err := getGuessRoom(r.joinRoomID)
	if err != nil {
		fmt.Printf("get guess room err:%v\n", err)
		return err.Error()
	}
	ir.Collect(r.GetEntityID(), content)
	fmt.Printf("Room chat  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "ChatRoom success"
}

func getGuessRoom(joinRoomID string) (IChainRoom, error) {
	if len(joinRoomID) == 0 {
		return nil, fmt.Errorf("ChatRoom failed, roomID is empty")
	}
	room, err := roomMgr.GetEntity(joinRoomID)
	if err != nil {
		return nil, fmt.Errorf("ChatRoom failed")
	}
	ir, ok := room.(IChainRoom)
	if !ok {
		return nil, fmt.Errorf("trans type illegal")
	}
	return ir, nil
}

/*
func (r *_User) ChatRoom(content string) string {
	if len(r.joinRoomID) == 0 {
		return "ChatRoom failed, roomID is empty"
	}
	room, err := roomMgr.AddOrGetEntity(r.joinRoomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)

		return "ChatRoom failed"
	}
	room.chat(r.GetEntityID(), r.GetEntityID(), content)
	fmt.Printf("Room chat  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "ChatRoom success"
}
*/
