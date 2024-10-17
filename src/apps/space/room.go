package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// 一些特别简单的狗屎代码

var roomMgr = &_RoomMgr{}

type _RoomMgr struct {
	rooms sync.Map // key: roomID val: room
}

func (mgr *_RoomMgr) AddEntity(userID string) {
	room := &_Room{
		createTime: time.Now().UnixNano(),
	}
	mgr.rooms.Store(userID, room)
}

func (mgr *_RoomMgr) AddOrGetEntity(userID string) (*_Room, error) {
	room := &_Room{
		createTime: time.Now().UnixNano(),
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.rooms.LoadOrStore(userID, room)
	if !ok {
		return room, nil
	}

	ret, transOk := entityInfo.(*_Room)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) GetEntity(userID string) (*_Room, error) {
	entityInfo, ok := mgr.rooms.Load(userID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(*_Room)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) DeleteEntity(userID string) {
	mgr.rooms.Delete(userID)
}

type _Room struct {
	memberList map[string]struct{}
	limit      uint8 // 最多255人
	createTime int64
	msgCache   []*_RoomChatMsg // 消息缓存
}

type _RoomChatMsg struct {
	fromID   string
	fromName string
	content  string
	sendTime int64
}

func (r *_Room) joinRoom(member string) {
	r.memberList[member] = struct{}{}
}

func (r *_Room) quitRoom(member string) {
	delete(r.memberList, member)
}

func (r *_Room) chat(member, memberName, content string) {
	m := fmt.Sprintf("%v chat: %v", memberName, content)
	RpcToEntityList(r.memberList, m)
}
