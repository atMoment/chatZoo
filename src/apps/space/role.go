package main

import (
	"ChatZoo/common"
	"fmt"
	"net"
)

// 设计上一个玩家只能加入一个房间
type _User struct {
	*common.EntityInfo
	joinRoomID string
}

func NewUser(entityID string, conn net.Conn) (*_User, error) {
	user := &_User{EntityInfo: common.NewEntityInfo(entityID, conn)}
	user.SetRpc(user)
	return user, nil
}

func (r *_User) CRPC_JoinRoom(roomID string) string {
	room, err := roomMgr.GetEntity(roomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	entity, ok := room.(common.IRoomEntity)
	if !ok {
		fmt.Printf("%v room can't trans entity\n ", r.GetEntityID())
		return "failed"
	}
	err = entity.SingleCall(ComponentBase+".PlayerJoinRoom", r.GetEntityID())
	if err != nil {
		fmt.Printf("%v single call failed, err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	r.joinRoomID = roomID
	fmt.Printf("CRPC_JoinRoom success  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CRPC_CreateRoom(typ int, roomID string, limit int) string {
	room, err := roomMgr.AddEntity(roomID, typ, limit)
	if err != nil {
		fmt.Printf("%v CRPC_CreateRoom add entity err:%v \n", r.GetEntityID(), err)
		return "failed"
	}
	entity, ok := room.(common.IRoomEntity)
	if !ok {
		fmt.Printf("%v room can't trans entity\n ", r.GetEntityID())
		return "failed"
	}
	err = entity.SingleCall(ComponentBase+".PlayerJoinRoom", r.GetEntityID())
	if err != nil {
		fmt.Printf("%v single call failed, err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	r.joinRoomID = roomID
	fmt.Printf("CRPC_CreateRoom success  userid:%v, roomid:%v \n ", r.GetEntityID(), roomID)
	return "success"
}

func (r *_User) CRPC_GetRecommendRoom() (string, []string) {
	return "success", roomMgr.TravelRoom()
}
func (r *_User) CRPC_QuitRoom(roomID string) string {
	// todo complete me
	return "success"
}

func (r *_User) CRPC_ChainRoomReady() string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	entity, ok := room.(common.IRoomEntity)
	if !ok {
		fmt.Printf("%v room can't trans entity\n ", r.GetEntityID())
		return "failed"
	}
	err = entity.SingleCall(ComponentChain+".Ready", r.GetEntityID())
	if err != nil {
		fmt.Printf("%v single call failed, err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	fmt.Printf("CRPC_GuessRoomReady success userid:%v, roomid:%v \n ", r.GetEntityID(), r.joinRoomID)
	return "success"
}

func (r *_User) CRPC_ChainRoomSendMsg(content string) string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_JoinRoom get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	entity, ok := room.(common.IRoomEntity)
	if !ok {
		fmt.Printf("%v room can't trans entity\n ", r.GetEntityID())
		return "failed"
	}
	err = entity.SingleCall(ComponentChain+".Collect", r.GetEntityID(), content)
	if err != nil {
		fmt.Printf("%v single call failed, err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	fmt.Printf("CRPC_GuessRoomSendMsg success  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "success"
}

func (r *_User) CRPC_ChatRoomSendMsg(content string) string {
	room, err := roomMgr.GetEntity(r.joinRoomID)
	if err != nil {
		fmt.Printf("%v CRPC_ChatRoomSendMsg get entity err:%v\n ", r.GetEntityID(), err)
		return "failed"
	}
	entity, ok := room.(common.IRoomEntity)
	if !ok {
		fmt.Printf("%v room can't trans entity\n ", r.GetEntityID())
		return "failed"
	}
	entity.SingleCall(ComponentChat+".Chat", r.GetEntityID(), content)
	fmt.Printf("CRPC_ChatRoomSendMsg success  userid:%v, roomid:%v content:%v\n ", r.GetEntityID(), r.joinRoomID, content)
	return "success"
}
