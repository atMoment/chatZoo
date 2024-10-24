package common

import (
	"errors"
	"fmt"
	"sync"
)

type IEntityMgr interface {
	AddEntity(entityID string, entity IEntityInfo)
	AddOrGetEntity(entityID string, entity IEntityInfo) (IEntityInfo, error)
	GetEntity(entityID string) (IEntityInfo, error)
	DeleteEntity(entityID string)
	TravelMgr(f func(key, value any) bool)
}

type _EntityMgr struct {
	entityList sync.Map // key: entityID  val: EntityInfo
}

func (mgr *_EntityMgr) Start() {
	fmt.Println("entity mgr start")
}

func (mgr *_EntityMgr) AddEntity(entityID string, entity IEntityInfo) {
	mgr.entityList.Store(entityID, entity)
}

func (mgr *_EntityMgr) AddOrGetEntity(entityID string, entity IEntityInfo) (IEntityInfo, error) {
	// 找不到 返回 false
	entityInfo, ok := mgr.entityList.LoadOrStore(entityID, entity)
	if !ok {
		return entity, nil
	}

	ret, transOk := entityInfo.(IEntityInfo)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_EntityMgr) GetEntity(entityID string) (IEntityInfo, error) {
	entityInfo, ok := mgr.entityList.Load(entityID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(IEntityInfo)
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
