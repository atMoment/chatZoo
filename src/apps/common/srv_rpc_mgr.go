package common

import (
	"fmt"
)

type IRpcMgr interface {
	SendNotifyToEntityList(userIds map[string]struct{}, arg ...interface{})
}

type _RpcMgr struct {
	IEntityMgr
}

func (s *_RpcMgr) SendNotifyToEntityList(userIds map[string]struct{}, arg ...interface{}) {
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
		entity.GetRpc().SendNotify(arg...)
		return true
	}
	s.TravelMgr(f)
}
