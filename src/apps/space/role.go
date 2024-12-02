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

func (r *_User) CRPC_JoinRoom(roomID string) string {
	room, err := roomMgr.GetEntity(roomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	room.JoinRoom(r.GetEntityID())
	r.joinRoomID = roomID
	r.joinRoomType = room.GetType()
	fmt.Printf("CRPC_JoinRoom success  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CRPC_CreateRoom(typ int, roomID string, limit int) string {
	room, err := roomMgr.AddEntity(roomID, typ, limit)
	if err != nil {
		fmt.Printf("%v CRPC_CreateRoom add entity err:%v \n", r.GetEntityID(), err)
		return "failed"
	}
	room.JoinRoom(r.GetEntityID())
	r.joinRoomID = roomID
	r.joinRoomType = typ
	fmt.Printf("CRPC_CreateRoom success  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CRPC_GetRecommendRoom() (string, []string) {
	return "success", roomMgr.TravelRoom()
}
func (r *_User) CRPC_QuitRoom(roomID string) string {
	return "success"
}

func (r *_User) CRPC_ChainRoomReady() string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	ir, ok := room.(IChainRoom)
	if !ok {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	ir.Ready(r.GetEntityID())
	fmt.Printf("CRPC_GuessRoomReady success userid:%v, roomid:%v \n ", r.GetEntityID(), r.joinRoomID)
	return "success"
}

func (r *_User) CRPC_ChainRoomSendMsg(content string) string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	ir, ok := room.(IChainRoom)
	if !ok {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	ir.Collect(r.GetEntityID(), content)
	fmt.Printf("CRPC_GuessRoomSendMsg success  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "success"
}

func (r *_User) CRPC_ChatRoomSendMsg(content string) string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_ChatRoomSendMsg get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	ir, ok := room.(IChatRoom)
	if !ok {
		fmt.Printf("%v CRPC_ChatRoomSendMsg trans failed\n ", r.GetEntityID())
		return "failed"
	}
	err = ir.Chat(r.GetEntityID(), content)
	if err != nil {
		fmt.Printf("%v CRPC_ChatRoomSendMsg chat err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	fmt.Printf("CRPC_ChatRoomSendMsg success  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "success"
}
