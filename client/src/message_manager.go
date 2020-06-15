package main

import (
	"errors"
	"fmt"
	"reflect"
)

type _MessageInfo struct {
	msg_info map[int32] interface{}
}
//var  MessageInfo  map[int32]interface{}

var minfo *_MessageInfo = NewMessageInfo()

func NewMessageInfo() *_MessageInfo {
	return &_MessageInfo{
		msg_info: make(map[int32]interface{}),
	}
}

func (m *_MessageInfo) addInfo(id int32, info interface{}) {
	m.msg_info[id] = info
}

func (m * _MessageInfo)GetMessageInfo(id int32) (interface{}, error) {
	info, ok := m.msg_info[id]
	if !ok {
		fmt.Println("[GetInfo] failed msg_id is ", id)
		return nil, errors.New("get message failed")
	}

	t := reflect.New(reflect.TypeOf(info).Elem())
	q := t.Interface().(interface{})

	return q, nil
}
func init() {

	minfo.addInfo(Request_join,  &ReqJoin{})
	minfo.addInfo(Response_join, &ResJoin{})
	minfo.addInfo(Request_chat,  &ReqChat{})
	minfo.addInfo(Response_chat, &ResChat{})
}

