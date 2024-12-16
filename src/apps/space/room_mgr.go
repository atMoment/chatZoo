package main

import (
	"ChatZoo/common/encrypt"
	"errors"
	"sync"
)

type IRoom interface {
	IRoomBase
}

var roomMgr = &_RoomMgr{}

type _RoomMgr struct {
	rooms sync.Map // key: roomID val: room
}

func (mgr *_RoomMgr) AddEntity(id string, typ, limit int) (IRoom, error) {
	if len(id) == 0 {
		return nil, errors.New("userid is empty")
	}
	_, ok := mgr.rooms.Load(id)
	if ok {
		return nil, errors.New("userid already exist")
	}
	room := createEntity(typ, id, limit)
	mgr.rooms.Store(id, room)
	return room, nil
}

func (mgr *_RoomMgr) AddOrGetEntity(id string, typ, limit int) (IRoom, error) {
	if len(id) == 0 {
		return nil, errors.New("userid is empty")
	}
	room := createEntity(typ, encrypt.NewGUID(), limit)
	// 找不到 返回 false
	entityInfo, ok := mgr.rooms.LoadOrStore(id, room)
	if !ok {
		return room, nil
	}

	ret, transOk := entityInfo.(IRoom)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) GetEntity(id string) (IRoom, error) {
	if len(id) == 0 {
		return nil, errors.New("userid is empty")
	}
	entityInfo, ok := mgr.rooms.Load(id)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(IRoom)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) TravelRoom() []string {
	ret := make([]string, 0)
	f := func(key, value any) bool {
		id, _ := key.(string)
		ret = append(ret, id)
		return true
	}
	mgr.rooms.Range(f)
	return ret
}

func (mgr *_RoomMgr) DeleteEntity(userID string) {
	mgr.rooms.Delete(userID)
}

func createEntity(typ int, entityID string, limit int) IRoom {
	room := NewRoom(entityID, limit)
	room.AddComponent(ComponentChain, NewChainComponent(room))
	room.AddComponent(ComponentChat, NewChatComponent(room))
	room.AddComponent(ComponentBase, NewBaseComponent(room))
	return room
}
