package main

import (
	"ChatZoo/common"
	"errors"
	"fmt"
	"net"
	"sync"
)

var entityMgr = NewEntityMgr()

type _EntityMgr struct {
	entityList sync.Map // key: userID
}

func NewEntityMgr() *_EntityMgr {
	return &_EntityMgr{}
}

type _UserInfo struct {
	user *_User
	conn net.Conn
}

func (mgr *_EntityMgr) AddEntity(userID string, conn net.Conn) {
	userInfo := &_UserInfo{
		conn: conn,
		user: &_User{
			userID: userID,
		},
	}
	mgr.entityList.Store(userID, userInfo)
}

func (mgr *_EntityMgr) AddOrGetEntity(userID string, conn net.Conn) (*_UserInfo, error) {
	userInfo := &_UserInfo{
		conn: conn,
		user: &_User{
			userID: userID,
		},
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.entityList.LoadOrStore(userID, userInfo)
	if !ok {
		return userInfo, nil
	}

	ret, transOk := entityInfo.(*_UserInfo)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) GetEntity(userID string) (*_UserInfo, error) {
	entityInfo, ok := mgr.entityList.Load(userID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(*_UserInfo)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) DeleteEntity(userID string) {
	mgr.entityList.Delete(userID)
}

func (mgr *_EntityMgr) TravelMgr(f func(key, value any) bool) {
	mgr.entityList.Range(f)
}

func RpcToEntity(userid string, arg string) error {
	entity, err := entityMgr.GetEntity(userid)
	if err != nil {
		return err
	}
	msg := &common.MsgCmdRsp{
		Arg: arg,
	}
	err = common.WriteToConn(entity.conn, msg)
	return err
}

func RpcToEntityList(userIds map[string]struct{}, arg string) {
	f := func(key, value any) bool {
		id, keyOk := key.(string)
		if !keyOk {
			fmt.Println("RpcToEntityList key not string")
			return false
		}
		_, find := userIds[id]
		if !find {
			return true
		}
		info, valOk := value.(*_UserInfo)
		if !valOk {
			fmt.Println("RpcToEntityList value not _UserInfo")
			return false
		}

		msg := &common.MsgCmdRsp{
			Arg: arg}
		err := common.WriteToConn(info.conn, msg)
		if err != nil {
			fmt.Println("RpcToEntityList wrtie conn err ", err)
		}

		return true
	}
	entityMgr.TravelMgr(f)
}
