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

func (r *_User) Calculate(expression string) string {
	var ret string
	result, err := calculate(expression)
	if err != nil {
		ret = fmt.Sprintf("session server calculate, err:%v", err)
	} else {
		ret = fmt.Sprintf("%s = %d", expression, result)
	}
	fmt.Printf("Calculate success expression:%v, ret:%v\n ", expression, ret)
	return ret
}

func (r *_User) JoinRoom(roomID string) string {
	room, err := roomMgr.AddOrGetEntity(roomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return "failed"
	}
	room.joinRoom(r.GetEntityID())
	r.joinRoomID = roomID
	fmt.Printf("Room JoinRoom  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CreateRoom(roomID string) {

}

func (r *_User) QuitRoom(roomID string) {

}

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
