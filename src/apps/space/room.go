package main

import (
	"errors"
	"sync"
)

// 一些特别简单的狗屎代码
const (
	RoomType_None = iota
	RoomType_Chat
	RoomType_Chain
)

type IRoom interface {
	IRoomBase
}

var roomMgr = &_RoomMgr{}

type _RoomMgr struct {
	rooms sync.Map // key: roomID val: room
	typ   int
}

func (mgr *_RoomMgr) AddEntity(userID string, typ, limit int) {
	mgr.rooms.Store(userID, createEntity(typ, limit))
}

func (mgr *_RoomMgr) AddOrGetEntity(userID string, typ, limit int) (IRoom, error) {
	room := createEntity(typ, limit)
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

func (mgr *_RoomMgr) GetEntity(userID string) (IRoom, error) {
	entityInfo, ok := mgr.rooms.Load(userID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(IRoom)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	switch ret.GetType() {
	case RoomType_None:
		entity, roomOk := entityInfo.(*_Room)
		if !roomOk {
			return entity, errors.New("can't trans room")
		}
		return entity, nil
	case RoomType_Chat:
		entity, roomOk := entityInfo.(*_ChatRoom)
		if !roomOk {
			return entity, errors.New("can't trans chat room")
		}
		return entity, nil
	case RoomType_Chain:
		entity, roomOk := entityInfo.(*_ChainRoom)
		if !roomOk {
			return entity, errors.New("can't trans chain room")
		}
		return entity, nil
	default:
		return nil, errors.New("room typ illegal")
	}
}

func (mgr *_RoomMgr) DeleteEntity(userID string) {
	mgr.rooms.Delete(userID)
}

func createEntity(typ int, limit int) IRoom {
	var room IRoom
	switch typ {
	case RoomType_Chain:
		room = NewChainRoom(limit)
	case RoomType_Chat:
		room = NewChatRoom(limit)
	default:
		room = NewRoom(limit)
	}
	return room
}
