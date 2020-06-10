// 网络底层函数
package main

type Message struct {
	size int32
	id int32
	data []byte
}

func NewMessage(id int32, data[] byte) *Message {
	msg := &Message {
		id : id,
		data: data,
		size : int32(len(data)) + 8,
	}
	return msg
}

func (msg *Message) GetID() int32 {
	return msg.id
}

func (msg *Message) GetString() []byte {
	return msg.data
}
