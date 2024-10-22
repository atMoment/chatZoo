package common

import (
	"errors"
	"net"
	"sync"
)

type _EntityMgr struct {
	entityList sync.Map // key: entityID  val: EntityInfo
}

type IEntityMgr interface {
	AddEntity(entityID string, conn net.Conn)
	AddOrGetEntity(entityID string, conn net.Conn) (*EntityInfo, error)
	GetEntity(entityID string) (*EntityInfo, error)
	DeleteEntity(entityID string)
	TravelMgr(f func(key, value any) bool)
}

func (mgr *_EntityMgr) AddEntity(entityID string, conn net.Conn) {
	userInfo := &EntityInfo{
		conn:     conn,
		entityID: entityID,
	}
	mgr.entityList.Store(entityID, userInfo)
}

func (mgr *_EntityMgr) AddOrGetEntity(entityID string, conn net.Conn) (*EntityInfo, error) {
	userInfo := &EntityInfo{
		conn:     conn,
		entityID: entityID,
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.entityList.LoadOrStore(entityID, userInfo)
	if !ok {
		return userInfo, nil
	}

	ret, transOk := entityInfo.(*EntityInfo)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) GetEntity(entityID string) (*EntityInfo, error) {
	entityInfo, ok := mgr.entityList.Load(entityID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(*EntityInfo)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) DeleteEntity(entityID string) {
	mgr.entityList.Delete(entityID)
}

func (mgr *_EntityMgr) TravelMgr(f func(key, value any) bool) {
	mgr.entityList.Range(f)
}
