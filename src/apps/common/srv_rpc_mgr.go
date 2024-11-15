package common

import (
	"fmt"
)

type IRpcMgr interface {
	SendNotifyToEntityList(userIds map[string]struct{}, methodName string, arg ...interface{})
	SendNotifyToEntity(userid string, methodName string, arg ...interface{})
}

type _RpcMgr struct {
	IEntityMgr
}

func NewRpcMgr(entityMgr IEntityMgr) *_RpcMgr {
	return &_RpcMgr{entityMgr}
}

func (s *_RpcMgr) SendNotifyToEntityList(userIds map[string]struct{}, methodName string, arg ...interface{}) {
	f := func(key, value any) bool {
		id, keyOk := key.(string)
		if !keyOk {
			fmt.Println("RpcToEntityList key not string")
			return false
		}
		_, find := userIds[id]
		if !find {
			fmt.Println("RpcToEntityList not find this guy ", id)
			return true
		}
		entity, valOk := value.(IEntityInfo)
		if !valOk {
			fmt.Println("RpcToEntityList value not _EntityInfo")
			return false
		}
		entity.GetRpc().SendNotify(methodName, arg...)
		return true
	}
	s.TravelMgr(f)
}

func (s *_RpcMgr) SendNotifyToEntity(userid string, methodName string, arg ...interface{}) {
	f := func(key, value any) bool {
		id, keyOk := key.(string)
		if !keyOk {
			fmt.Println("RpcToEntityList key not string")
			return false
		}
		if id != userid {
			return true
		}
		entity, valOk := value.(IEntityInfo)
		if !valOk {
			fmt.Println("RpcToEntityList value not _EntityInfo")
			return false
		}
		entity.GetRpc().SendNotify(methodName, arg...)
		return false // 退出循环
	}
	s.TravelMgr(f)
}
