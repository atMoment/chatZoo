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
		size : int32(len(data)) + 4,
	}
	return msg
}

func (msg *Message) GetString() []byte {
	return msg.data
}
