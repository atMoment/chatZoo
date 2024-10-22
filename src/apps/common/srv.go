package common

type _Srv struct {
	entityMgr *_EntityMgr // 可以换成interface
	rpcMgr    *_RpcMgr
}
