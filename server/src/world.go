package main

import (
	"fmt"
	"sync"
)
const Request_join = 1
const Response_join = 2
const Response_chat = 3
const Request_chat = 4

type World struct {
	rooms *sync.Map
}


func NewWorld() *World {
	w := &World {
		rooms: & sync.Map{},
	}
	return w
}

func (w * World) AddRoom(rid int) *Room {
	q, ok := w.rooms.Load(rid)
	if !ok {
		r := NewRoom()
		w.rooms.Store(rid, r)
		return r
	}else{
		return q.(*Room)
	}
}

func (w *World) RemoveRoom(rid int) {
	_, ok := w.rooms.Load(rid)
	if ok {
		w.rooms.Delete(rid)
	}
}

func (w *World) DealMessage (msg *Message, s *Service) {
	id := msg.GetID()

	switch id {
	case Request_join:
		w.ResponseJoin(msg, s)
		// 打印聊天记录
		fmt.Println("congraulation you join!!!!")
	case Request_chat:
		// 广播用户说的话
		w.ResponseChat(msg, s)
	}
}

func (w *World) ResponseJoin(msg *Message, s *Service) {
	str := msg.GetString()
	var rid int
	fmt.Scanln(str, &rid)

	rstr := "congraulation  join!"
	rmsg := NewMessage(Response_join, []byte(rstr))

	r := w.AddRoom(rid)
	list := r.TravelUser()
	s.SendAll(rmsg, list)
}

func (w *World) ResponseChat(msg *Message, s *Service) {
	//fmt.Printf("msg is %s", msg.GetString())
	// 广播给大家
	// 这里应该建造一个实体类，保存名字，房间号等信息。就不用在消息里获得了
	//str := msg.GetString()
	//var rid int
	//var rstr string
	//fmt.Scanln(str, &rid, &rstr)

	rmsg := NewMessage(Response_chat, msg.GetString())
	fmt.Printf("response chat is rstr is %s\n", msg.GetString())
	r := w.AddRoom(1)
	list := r.TravelUser()

	for k, _ := range list {
		fmt.Println("range k is ", k)
	}
	s.SendAll(rmsg, list)
}

