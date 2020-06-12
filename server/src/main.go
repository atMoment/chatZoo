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
// 多多使用interface
// 考虑粘包拆包的问题， 使用io.ReadFull可以解决，一直阻塞直到读到固定的长度。而使用conn.Read，如果计算机内部拆包，把字节为4的包拆成1 和 3，那么读到
// 1就会返回了，err也为空。需要对长度进行判断

// 关于session应该有handel函数，对应每一个会话处理。
// 为了保证消息的有序处理，消息由有先后。虽然信号量到达有先后，但是处理函数是协程，不能保证消息有序处理。
// 关于Cinder的message只用了一个interface, 如果使用了 int id, int name ,难道还要用一个结构体包起来？ 可以支持传多参吗？不用
// 粘包和拆包的问题，是不是用的安全链表处理的。不是，是为了解决阻塞。
// 用通道的话，chanel会有长度。如果通道满了，那么消息就会阻塞等待或者直接丢掉。而用安全链表的话，就可以长度不限。
// cinder的消息序列化拆解是先从结构体入手，因为一个消息就是一个结构体，是各种各样的结构体。从里面去解析，递归解析。
// 重大设计问题。 玩家上线后如果不进入房间，而此时另一个玩家上线进入房间，就会丢失当前玩家信息。
// 而客户端下线就销毁玩家信息，则玩家是当前玩家，如果是以前的玩家下线，那么不久销毁错了吗？  想得不够清楚，这是大问题！
// 再说请求聊天，得到当前玩家的房间号。如果1号玩家发言而得到后连接的2号玩家的房间号，是不是发错言了。
// 程序是一定要运行多测试的。

