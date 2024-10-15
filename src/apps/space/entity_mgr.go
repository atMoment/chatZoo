package main

import (
	"ChatZoo/common"
	"errors"
	"net"
	"sync"
)

var entityMgr *_EntityMgr

type _EntityMgr struct {
	entityList sync.Map // key: roleID
}

func init() {
	entityMgr = NewEntityMgr()
}

func NewEntityMgr() *_EntityMgr {
	return &_EntityMgr{}
}

type _RoleInfo struct {
	role      *Role
	sessionID string
	conn      net.Conn
}

func (mgr *_EntityMgr) AddEntity(roleID string, sessionID string, conn net.Conn) {
	roleInfo := &_RoleInfo{
		sessionID: sessionID,
		conn:      conn,
		role: &Role{
			roleID: roleID,
		},
	}
	mgr.entityList.Store(roleID, roleInfo)
}

func (mgr *_EntityMgr) AddOrGetEntity(roleID string, sessionID string, conn net.Conn) (*_RoleInfo, error) {
	roleInfo := &_RoleInfo{
		sessionID: sessionID,
		conn:      conn,
		role: &Role{
			roleID: roleID,
		},
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.entityList.LoadOrStore(roleID, roleInfo)
	if !ok {
		return roleInfo, nil
	}

	ret, transOk := entityInfo.(*_RoleInfo)
	if !transOk {
		return nil, errors.New("trans roleinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) GetEntity(roleID string) (*_RoleInfo, error) {
	entityInfo, ok := mgr.entityList.Load(roleID)
	if !ok {
		return nil, errors.New("roleid not find")
	}
	ret, ok := entityInfo.(*_RoleInfo)
	if !ok {
		return nil, errors.New("trans roleinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) DeleteEntity(roleID string) {
	mgr.entityList.Delete(roleID)
}

func Rpc(roleid string, arg string) error {
	entity, err := entityMgr.GetEntity(roleid)
	if err != nil {
		return err
	}
	msg := &common.MsgCmdRsp{
		Arg: arg,
	}
	err = common.WriteToConn(entity.conn, msg)
	return err
}
