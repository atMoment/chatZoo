package main

import (
	"fmt"
)

type Role struct {
	roleID string
}

func (r *Role) Calculate(expression string) {
	var ret string
	result, err := calculate(expression)
	if err != nil {
		ret = fmt.Sprintf("session server calculate, err:%v", err)
	} else {
		ret = fmt.Sprintf("%s = %d", expression, result)
	}
	Rpc(r.roleID, ret)
	fmt.Printf("Calculate success expression:%v, ret:%v\n ", expression, ret)
}

func (r *Role) JoinRoom(roomID string) {

}

func (r *Role) QuitRoom(roomID string) {

}

func (r *Role) ChatRoom(roomID, content string) {

}
