// 上层函数
package main

import (
	"sync"
)

const MaxRecord  = 5


type Room struct {
	rid int                          // 房间标识符
	name string                      // 房间名字
	players  *sync.Map               // 房间内玩家
	//record []byte                  // 消息记录  暂时没有对消息记录做操作的打算
}

func NewRoom() *Room {
	r := &Room {
		players: &sync.Map{},
		name:"",
	}
	return r
}

// 房间销毁，把房间内玩家全踢下线
func (r *Room) DeepDestroy(k, v interface{}) {
	v.(*Player).Destroy()                          // 把玩家踢下线
	r.players.Delete(k)                            //一边遍历一遍删除会有问题吗？会迭代器失效吗？
}

// 房间销毁，把房间内玩家全踢出去
func (r *Room) Destroy(k, v interface{}){
	v.(*Player).SetRid(0)
	r.players.Delete(k)
}

func (r *Room) TravelPlayers(f func(k, v interface{})) {
	mf := func(k, v interface{}) bool {
		f(k, v)
		return true
	}
	r.players.Range(mf)
}

func (r *Room) AddUser(uid string, player *Player) {
	_, ok := r.players.LoadOrStore(uid, player)
	if ok {
		player.uid = uid
	}
}

func (r *Room) RemoveUser(uid int) {
	r.players.Delete(uid)
}

// 遍历该房间所有玩家，并返回uid集合
func (r *Room) GetPlayers() []string {
	list := make([]string, 0)
	f := func(k, v interface{}) bool {
		list = append(list, k.(string))
		return true
	}
	r.players.Range(f)

	return list
}

// 给房间内所有玩家发消息
func (r *Room) SendRoom(msg *Message) {
	mf := func(k, v interface{}) bool {
		v.(*Player).SendUser(msg)
		return true
	}
	r.players.Range(mf)
}

func (r *Room) OnRecvMessage(from string, content string) {

}

// 玩家上下线问题
// 没有账号，没有登录，每次都是临时玩家，提供临时的操作，自然也没有消息记录。临时玩家不需要消息记录，也没有。
// 1、服务器重启以后，玩家所有数据丢失。 玩家数据存在服务器缓存
// 2、服务器重启以后，玩家数据不会丢失。存在程序外部json文件或者mysql等外部存储设备
// 现在做的是1

