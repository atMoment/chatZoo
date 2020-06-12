// 网络底层函数
package main

type Message struct {
	size int32
	id int32
	data interface{}
}

func NewMessage(id int32, data interface{}) *Message {
	msg := &Message {
		id : id,
		data: data,
	}
	return msg
}

func (msg *Message) GetID() int32 {
	return msg.id
}

func (msg *Message) GetString() interface{} {
	return msg.data
}

func (msg *Message) SetSize(size int32){
	msg.size = size
}

func (msg *Message) GetSize() int32{
	return msg.size
}
