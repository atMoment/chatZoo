// 上层函数
package main

import (
	"sync"
)

const MaxRecord  = 5


type Room struct {
	players  *sync.Map
}

type Player struct {
	name string
	record []byte
}


func NewRoom() *Room {
	r := &Room {
		players: &sync.Map{},
	}
	return r
}

func NewPlayer() *Player {
	p := &Player{}
	return p
}
func (r *Room) AddUser(uid int) {
	_, ok := r.players.Load(uid)
	if !ok {
		p := NewPlayer()
		r.players.Store(uid, p)
	}
}

func (r *Room) RemoveUser(uid int) {
	_, ok := r.players.Load(uid)
	if ok {
		r.players.Delete(uid)
	}
}

// 遍历该房间所有玩家，并返回uid集合
func (r *Room) TravelUser() []string {
	list := make([]string, 0)
	f := func(k, v interface{}) bool {
		list = append(list, k.(string))
		return true
	}
	r.players.Range(f)

	return list
}

func (r *Room) OnRecvMessage(from string, content string) {

}
