package main

import (
	"errors"
	"fmt"
	"reflect"
)

var  MessageInfo  map[int32]interface{}




func addInfo(id int32, info interface{}) {
	MessageInfo[id] = info
}

func GetMessageInfo(id int32) (interface{}, error) {
	info, ok := MessageInfo[id]
	if !ok {
		fmt.Println("[GetInfo] failed msg_id is ", id)
		return nil, errors.New("get message failed")
	}

	m := reflect.New(reflect.TypeOf(info).Elem())
	q := m.Interface().(interface{})

	return q, nil
}
func init() {
	MessageInfo = make(map[int32]interface{})
	addInfo(Request_join,  &RequestJoin{})
	addInfo(Response_join, &ResponseJoin{})
	addInfo(Request_chat,  &RequestChat{})
	addInfo(Response_chat, &ResponseChat{})
}

