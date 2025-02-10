package main

import "sync"

type IGateSrv interface {
	Destroy()
	AddGateUser(entityID string, gateUser IGateUser)
	RemGateUser(entityID string)
}
type _GateSrv struct {
	gateUserMap sync.Map
}

var DefaultGateSrvEntity IGateSrv

func init() {
	DefaultGateSrvEntity = &_GateSrv{}
}

func (s *_GateSrv) Destroy() {
	f := func(key, value any) bool {
		value.(IGateUser).Destroy()
		return true
	}
	s.gateUserMap.Range(f)
}

func (s *_GateSrv) AddGateUser(entityID string, gateUser IGateUser) {
	user, ok := s.gateUserMap.LoadOrStore(entityID, gateUser)
	if ok {
		user.(IGateUser).Destroy()
		s.RemGateUser(entityID)
	}
	s.gateUserMap.Store(entityID, gateUser)
}

func (s *_GateSrv) RemGateUser(entityID string) {
	s.gateUserMap.Delete(entityID)
}
