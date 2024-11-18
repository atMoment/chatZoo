package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
)

type _User struct {
	*common.EntityInfo
	joinRoomID string
}

func NewUser(entityID string, conn net.Conn) (*_User, error) {
	user := &_User{EntityInfo: common.NewEntityInfo(entityID, conn)}
	user.SetRpc(user)

	return user, nil
}

func (r *_User) JoinRoom(roomID string) string {
	room, err := roomMgr.AddOrGetEntity(roomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return "failed"
	}
	room.JoinRoom(r.GetEntityID())
	r.joinRoomID = roomID
	fmt.Printf("Room JoinRoom  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CreateRoom(roomID string) {

}

func (r *_User) QuitRoom(roomID string) {

}

func (r *_User) Ready() string {
	if len(r.joinRoomID) == 0 {
		return "ChatRoom failed, roomID is empty"
	}
	room, err := roomMgr.AddOrGetEntity(r.joinRoomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)

		return "ChatRoom failed"
	}
	room.(r.GetEntityID(), r.GetEntityID(), content)
	fmt.Printf("Room chat  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "ChatRoom success"
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
