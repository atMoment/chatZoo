package main

import (
	"ChatZoo/common"
	"errors"
	"time"
)

type IRoomBase interface {
	GetType() int
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
	typ        int
	memberList map[string]struct{}
	limit      int
	createTime int64
}

func NewRoom(limit int) *_Room {
	return &_Room{
		typ:        RoomType_None,
		limit:      limit,
		createTime: time.Now().Unix(),
		memberList: make(map[string]struct{}),
	}
}

func (r *_Room) GetType() int {
	return r.typ
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
