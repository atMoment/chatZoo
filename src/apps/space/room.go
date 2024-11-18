package main

import (
	"ChatZoo/common"
	"errors"
	"sync"
	"time"
)

// 一些特别简单的狗屎代码
const (
	_ = iota
	RoomType_Chat
	RoomType_Guess
)

var roomMgr = &_RoomMgr{}

type _RoomMgr struct {
	rooms sync.Map // key: roomID val: room
	typ   int
}

func (mgr *_RoomMgr) AddEntity(userID string) {
	room := &_Room{
		createTime: time.Now().UnixNano(),
	}
	mgr.rooms.Store(userID, room)
}

func (mgr *_RoomMgr) AddOrGetEntity(userID string, typ int) (IRoom, error) {
	switch typ {
	case RoomType_Guess:
		room := NewChainRoom(4)

	}
	room := &_Room{
		createTime: time.Now().UnixNano(),
		memberList: make(map[string]struct{}),
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.rooms.LoadOrStore(userID, room)
	if !ok {
		return room, nil
	}

	ret, transOk := entityInfo.(IRoom)
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

type IRoom interface {
	JoinRoom(member string) error
	QuitRoom(member string)
	MemberIsExist(member string) bool
	GetRoomMemberList() []string // 只能读
	GetRoomMemberLimit() int     // 只能读
	NotifyAllMember(methodName string, args ...interface{})
	NotifyMember(member string, methodName string, args ...interface{}) error
}

// todo 还有创建销毁, load/unload 逻辑

type _Room struct {
	memberList map[string]struct{}
	limit      int
	createTime int64
}

func NewRoom(limit int) IRoom {
	return &_Room{
		limit:      limit,
		createTime: time.Now().Unix(),
		memberList: make(map[string]struct{}),
	}
}

func (r *_Room) JoinRoom(member string) error {
	if len(r.memberList) == int(r.limit) {
		return errors.New("room full")
	}
	r.memberList[member] = struct{}{}
	return nil
}

func (r *_Room) QuitRoom(member string) {
	delete(r.memberList, member)
}

func (r *_Room) MemberIsExist(member string) bool {
	_, ok := r.memberList[member]
	return ok
}

func (r *_Room) GetRoomMemberLimit() int {
	return r.limit
}

func (r *_Room) GetRoomMemberList() []string {
	ret := make([]string, len(r.memberList))
	var i int
	for m := range r.memberList {
		ret[i] = m
		i++
	}
	return ret
}

func (r *_Room) NotifyAllMember(methodName string, args ...interface{}) {
	common.DefaultSrvEntity.SendNotifyToEntityList(r.memberList, methodName, args...)
}

func (r *_Room) NotifyMember(member string, methodName string, args ...interface{}) error {
	if _, ok := r.memberList[member]; !ok {
		return errors.New("member not in room")
	}
	common.DefaultSrvEntity.SendNotifyToEntity(member, methodName, args...)
	return nil
}
