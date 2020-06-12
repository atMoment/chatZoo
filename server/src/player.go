package main


type Player struct {
	uid  string                 // 玩家标识
	rid  int                 // 玩家所属房间标识
	name string              // 玩家名字
	sess  *_TcpSession           // 管理自己的会话
	// readno int            // 所读消息号   现在没有对消息有操作的打算，应该有一个消息系统，玩家关联消息系统
}

func NewPlayer(uid string, session *_TcpSession) *Player {
	p := &Player{
		uid:uid,
		rid: 0,
		name: "",
		sess: session,
	}
	return p
}

func (u *Player) Destroy() {
	u.sess.Destroy()
}

func (u *Player) SetRid(rid int) {
	u.rid = rid
}

func (u *Player) GetRid() int {
	return u.rid
}

func (u *Player) GetUid() string {
	return u.uid
}

func (u *Player) SendUser(msg *Message) {
	u.sess.Send(msg)
}

func (u *Player) Loop() {
	
}
