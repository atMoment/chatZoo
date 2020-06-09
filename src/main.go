// 调用函数入口
package main


func main() {
	host := "127.0.0.1:3344"
	connect_type := "tcp"

	w := NewWorld()
	s := NewService(host, connect_type)
	s.ReqMessageHandle(w.DealMessage)
	s.AcceptConn()
}


// 需要改的地方
// service 的功能很复杂，需要明确分摊到每个步骤层面 比如sendALL
// 需要新增一个类型 player,由player管理会话
// session 应该有handel函数，对应每一个会话处理
// 使用 sync.map 就直接使用他的原子操作。 加的if条件反而会导致协程并发的问题
// 聊天记录放room身上，玩家本地保存已读的序号
// 底层上层使用包分隔开
