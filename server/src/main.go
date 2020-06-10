// 调用函数入口
package main


func main() {
	host := "127.0.0.1:3344"
	connect_type := "tcp"

	w := NewWorld()    // 创建一个世界
	s := NewService(host, connect_type)
	s.RegisterMessageHandle(w.DealMessage)
	s.RegisterConnectHandle(w.ConnectHandle)
	s.RegisterDisConnectHandle(w.DisconnetHandle)
	s.AcceptConn()
}


// 需要改的地方
// service 的功能很复杂，需要明确分摊到每个步骤层面 比如sendALL
// 需要新增一个类型 player,由player管理会话
// session 应该有handel函数，对应每一个会话处理
// 使用 sync.map 就直接使用他的原子操作。 加的if条件反而会导致协程并发的问题
// 聊天记录放room身上，玩家本地保存已读的序号
// 底层上层使用包分隔开
// 系列化、反序列化都要写好。用fmt.Scanln从string拆分int会很耗
// 完善销毁函数，以及客户端短线
// 多多使用interface
// 考虑粘包拆包的问题
