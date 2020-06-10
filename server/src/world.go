// 世界管理器
// 管理多个房间和一个玩家
package main

import (
	"fmt"
	"sync"
)

// 现在是临时行为，没有注册登录等存外部存储的操作，服务器重启时，所有房间和玩家都销毁。
// 可以把数据都存在json里，下次读
// 服务器不重启，玩家上线下线（客户端开闭），这个要做把
// 下线后再上线， 客户端端口变了，uid变了，现在没有账号的，上线下线也在做不了。

type World struct {
	rooms *sync.Map
	player *Player
	//maptable map[int] interface{}    回调函数关系可以用关系表存储
}

func NewWorld() *World {
	w := &World {
		rooms: & sync.Map{},
		player: nil,
		//maptable: make(map[int] interface{}),
	}
	return w
}

// 服务器关闭，世界销毁就销毁所有房间、当前玩家
func (w *World) Destroy() {
	// 每个房间每个玩家sess依次关闭
	f := func(k, v interface{}) bool {
		v.(*Room).TravelPlayers(v.(*Room).DeepDestroy)
		return true
	}
	// 世界里每个房间依次走这个流程
	w.rooms.Range(f)
	w.player.Destroy()
}

func (w * World) GetRoom(rid int) *Room {
    r, _ := w.rooms.LoadOrStore(rid, NewRoom())
    return r.(*Room)
}

func (w *World) RemoveRoom(rid int) {

	w.rooms.Delete(rid)
}

// 收到客户端消息处理
// 整体处理消息，可以考虑分开，执行一个函数， 内部处理由各个系统自己定义 ，怎么写
func (w *World) DealMessage (msg *Message) {
	id := msg.GetID()

	switch id {
	case Request_join:
		w.ResponseJoin(msg)
		// 打印聊天记录
	case Request_chat:
		// 广播用户说的话
		w.ResponseChat(msg)
	}
}

func (w *World) ResponseJoin(msg *Message) {
	str := msg.GetString()
	var rid int
	fmt.Scanln(str, &rid)                          // 这个很耗
	if w.player == nil {
		fmt.Println("world player is nil!!")
		return
	}else {
		w.player.SetRid(rid)
	}
	r := w.GetRoom(rid)
	r.AddUser(w.player.GetUid(), w.player)

	rstr := "congraulation  join!!!!"
	rmsg := NewMessage(Response_join, []byte(rstr))

	r.SendRoom(rmsg)
}

func (w *World) ResponseChat(msg *Message) {
	rmsg := NewMessage(Response_chat, msg.GetString())
	//fmt.Printf("response chat is rstr  %s\n", msg.GetString())
	var rid int
	if w.player == nil {
		fmt.Println("world player is nil!!")
		return
	}else{
		rid = w.player.GetRid()
	}
	w.SendRoom(rid, rmsg)
}

// 客户端连接成功就创建一个玩家
func (w *World) ConnectHandle(sess *_TcpSession) {
	w.player = NewPlayer(sess.GetSessionUid(), sess)
}

// 客户端下线就销毁该玩家
func (w *World) DisconnetHandle() {
	if (w.player != nil) {
		w.player.Destroy()
		w.player = nil
	}
}

func (w *World) SendRoom(rid int, msg *Message) {
	r := w.GetRoom(rid)
	r.SendRoom(msg)
}


