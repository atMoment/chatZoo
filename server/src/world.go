// 世界管理器
// 管理多个房间和一个玩家
package main

import (
	"fmt"
	"sync"
)


type World struct {
	rooms *sync.Map
	player *sync.Map
	//maptable map[int] interface{}    回调函数关系可以用关系表存储
}

func NewWorld() *World {
	w := &World {
		rooms: & sync.Map{},
		player: &sync.Map{},
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
	// 销毁玩家列表稍后再补
	//w.player.Destroy()
}

func (w * World) GetRoom(rid int) *Room {
    r, _ := w.rooms.LoadOrStore(rid, NewRoom())
    return r.(*Room)
}

func (w *World) RemoveRoom(rid int) {

	w.rooms.Delete(rid)
}

func (w *World) GetPlayer(uid string) *Player {
	u, _ := w.player.Load(uid)
	return u.(*Player)
}

func (w *World) AddPlayer (player * Player) {
	w.player.Store(player.GetUid(), player)
}

func (w *World) RemovePlayer(uid string) {
	w.player.Delete(uid)
}

// 收到客户端消息处理
// 整体处理消息，可以考虑分开，执行一个函数， 内部处理由各个系统自己定义 ，怎么写
func (w *World) DealMessage (msg *Message,  sess *_TcpSession) {
	id := msg.GetID()

	switch id {
	case Request_join:
		w.ResponseJoin(msg, sess)
		// 打印聊天记录
	case Request_chat:
		// 广播用户说的话
		w.ResponseChat(msg, sess)
	}
}

func (w *World) ResponseJoin(msg *Message, sess *_TcpSession) {
	str := msg.GetString()
	var rid int
	fmt.Scanln(str, &rid)                          // 这个很耗

	uid := sess.GetSessionUid()
	u := w.GetPlayer(uid)
	fmt.Println("response rid is , str is", rid, str)
	if u == nil {
		fmt.Println("player is not regisiter in world")
		return
	}else {
		u.SetRid(rid)
	}
	r := w.GetRoom(rid)
	fmt.Println("")
	r.AddUser(uid, u)
	_ = r.GetPlayers()

	rstr := "congraulation  join!!!!"
	rmsg := NewMessage(Response_join, rstr)

	r.SendRoom(rmsg)
}

func (w *World) ResponseChat(msg *Message, sess *_TcpSession) {
	reschat := ResChat{msg.GetString().(ReqChat).words}
	rmsg := NewMessage(Response_chat, reschat)
	//fmt.Printf("response chat is rstr  %s\n", msg.GetString())

	uid := sess.GetSessionUid()
	u := w.GetPlayer(uid)

	var rid int
	if u == nil {
		fmt.Println("world player is nil!!")
		return
	}else{
		rid = u.GetRid()
	}
	w.SendRoom(rid, rmsg)
}

// 客户端连接成功就创建一个玩家
func (w *World) ConnectHandle(sess *_TcpSession) {
	player := NewPlayer(sess.GetSessionUid(), sess)
	w.AddPlayer(player)
}

// 客户端下线就销毁该玩家
func (w *World) DisconnetHandle(sess *_TcpSession) {
	uid := sess.GetSessionUid()
	u := w.GetPlayer(uid)
	if (u == nil) {
		fmt.Println("player is nil")
		return
	}
	u.Destroy()
	w.RemovePlayer(uid)
}

func (w *World) SendRoom(rid int, msg *Message) {
	r := w.GetRoom(rid)
	r.SendRoom(msg)
}
