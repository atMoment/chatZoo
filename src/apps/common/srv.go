package common

const (
	KickUserReason_Online      = "online"
	KickUserReason_NoHeartBeat = "no heart beat"
)

type ISrv interface {
	AddEntityToMgr(entityID string, entity IEntityInfo)
	RemEntityFromMgr(entityID, reason string)
	Destroy()

	IEntityMgr
	IRpcMgr
}

type _Srv struct {
	*_EntityMgr // 可以换成interface
	*_RpcMgr
}

var DefaultSrvEntity ISrv

func init() {
	srv := &_Srv{
		_EntityMgr: NewEntityMgr(),
	}
	srv._RpcMgr = NewRpcMgr(srv._EntityMgr)
	DefaultSrvEntity = srv
}

func (s *_Srv) Destroy() {
	s.IEntityMgr.Destroy()
}

func (s *_Srv) AddEntityToMgr(entityID string, entity IEntityInfo) {
	find := s.LoadAndDeleteEntity(entityID)
	if find {
		s.SendLogout(entityID, KickUserReason_Online)
	}
	s.AddEntity(entityID, entity)
}

func (s *_Srv) RemEntityFromMgr(entityID, reason string) {
	s.DeleteEntity(entityID)
	s.SendLogout(entityID, reason)
}
