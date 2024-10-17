package main

import (
	"fmt"
)

type _User struct {
	userID string
}

func (r *_User) Calculate(expression string) {
	var ret string
	result, err := calculate(expression)
	if err != nil {
		ret = fmt.Sprintf("session server calculate, err:%v", err)
	} else {
		ret = fmt.Sprintf("%s = %d", expression, result)
	}
	RpcToEntity(r.userID, ret)
	fmt.Printf("Calculate success expression:%v, ret:%v\n ", expression, ret)
}

func (r *_User) JoinRoom(roomID string) {
	room, err := roomMgr.AddOrGetEntity(roomID)
	if err != nil {
		fmt.Println("chat get entity err ", err)
		return
	}
	room.joinRoom(r.userID)
	fmt.Printf("Room JoinRoom  userid:%v, roomid:%v \n ", r.userID, roomID)
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
	room.chat(r.userID, roomID, content)
	fmt.Printf("Room chat  userid:%v, roomid:%v content:%v\n ", r.userID, roomID, content)
}
