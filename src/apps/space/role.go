package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
)

type _User struct {
	*common.EntityInfo
}

func NewUser(entityID string, conn net.Conn) *_User {
	user := &_User{common.NewEntityInfo(entityID, conn)}
	user.SetRpc(user)
	return user
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

func (r *_User) JoinRoom(roomID string) {
	room, err := roomMgr.AddOrGetEntity(roomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return
	}
	room.joinRoom(r.GetEntityID())
	fmt.Printf("Room JoinRoom  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
}

func (r *_User) CreateRoom(roomID string) {

}

func (r *_User) QuitRoom(roomID string) {

}

func (r *_User) ChatRoom(roomID, content string) {
	room, err := roomMgr.AddOrGetEntity(roomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return
	}
	room.chat(r.GetEntityID(), r.GetEntityID(), content)
	fmt.Printf("Room chat  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), roomID, content)
}
