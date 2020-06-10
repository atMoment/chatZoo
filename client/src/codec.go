package main

import (
	"bytes"
	"encoding/binary"
)


// decode from []byte to Message
// 读
func Decode (data [] byte) (*Message, error) {
	buf_reader := bytes.NewReader(data)

	data_size := int32(len(data))

	var msg_id int32
	err := binary.Read(buf_reader, binary.LittleEndian, &msg_id)
	if err != nil {
		return nil, err
	}

	length := data_size - 4   // 减去刚刚读的id字节
	data_buf := make([]byte, length)
	err = binary.Read(buf_reader, binary.LittleEndian, &data_buf)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.data = data_buf
	message.id = msg_id
	message.size = length + 8

	return message, nil
}

// encode from Message to []byte
// 写
func Encode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.size)
	if err != nil {
		return nil, err
	}

	err =  binary.Write(buffer, binary.LittleEndian, msg.id)
	if err != nil {
		return nil, err
	}

	err =  binary.Write(buffer, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}