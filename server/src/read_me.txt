 提供的功能有：
    1、玩家加入某个房间
    2、玩家在某个房间发言
 未来可能支持的功能有：
    1、消息记录
    2、玩家上下线
    ps:   因为现在是临时行为，没有注册登录等存外部存储的操作，服务器重新启动时，所有房间和玩家都销毁, 再重开时，是新的房间
           下线后再上线，客户端端口变了（现在玩家的唯一标识uid是 客户端远程ip  + 端口号）， uid变了，无法找到玩家。   因为没有账号，上下线做不了。每次都是游客登录
    方向： 可以考虑把所有数据都存在json里
               消息记录，把消息结构放房间上，玩家本地保存一个已读的序号，上线时，服务器主动推送所有消息序号（1-10）.注意服务器房间消息有长度，需要更新

  高级建议和修改
  1、当前设计架构服务器只能处理收到客户端消息如何做， 如果考虑服务器向客户端主动推送消息，如何做？
  2、关于session消息处理，可以考虑另外一种实现方式，开读协程、写协程、消息处理协程。 可以把会话连接、会话断连放在消息处理协程里。和当前处理方式有什么差异？
      ps： 注意，消息的处理是有顺序的。
  3、关于消息发送问题，是否需要将客户端的小包消息合成大包进行处理？

 触手可及的修改
   1、分包。把源文件按照层次放在不同的包里，使程序简单易读
   2、服务器处理客户端断线，以及自己宕机的完善。  && 客户端处理自己断线以及服务器宕机
   3、服务器考虑粘包拆包问题。底层读套接字的时候使用io.ReadFull 代替 conn.Read
   4、channel 收到消息的时候，channel有规定长度，如果消息太多来不及处理，channel满了，是自己阻塞还是手动写阻塞？或者把消息丢掉？怎么写？
        ps: 更好的方法是用一个协程安全链表，可以长度不限，动态增减。（太麻烦了不想写了）
   5、消息序列化 按照现有方式非常局限，可以考虑更通用的修改吗？
   6、请求连接会把session id传进来用来当玩家的uid，去查找是哪个玩家发的消息。流程结构能改下吗？在消息
   7、消息是怎么知道有多大的？
   8、客户端、服务器消息处理（收到消息后对应去做哪个函数）可以优化一波（现在两边处理方式不一样）

bug复盘
1、
   Q：客户端隔房间不同也能看到彼此的消息
   A： 服务器在接受请求加入房间的时候，从消息里读取房间号失败，默认都加在房间号为0的房间。  等消息完全序列化再测试，已完成
2、
  Q : 客户端断线后销毁混乱，客户端发言后消息混乱（隔房间不同也能看到彼此消息）
   A：上次把由网络底层管理所有客户端会话改成上层player管理自己会话。但是world管理器值管理当前玩家，如果当前玩家被替代，则无法找到原来玩家，消息处理错误。本次更新已解决 
3、
   Q：在发送消息的时候，根据结构类型将字节流写入类型会报错 reflect.CanSet  返回false
   A：非导出字段不能写，需要将结构里的变量名命名为大写开头，变成可导出字段