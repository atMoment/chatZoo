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
	DefaultSrvEntity = &_Srv{
		&_EntityMgr{},
		&_RpcMgr{},
	}
}
