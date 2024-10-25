package common

type ISrv interface {
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
