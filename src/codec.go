// 网络底层函数
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// decode from []byte to Message
// 读
func Decode (data [] byte) (*Message, error) {
	buf_reader := bytes.NewReader(data)

	data_size := int32(len(data))

	var msg_id int32
	err0 := binary.Read(buf_reader, binary.LittleEndian, &msg_id)
	if err0 != nil {
		fmt.Println("read id failed err is", err0)
		return nil, err0
	}
	fmt.Println("read id", msg_id)
	length := data_size - 4
	data_buf := make([]byte, length)
	err1 := binary.Read(buf_reader, binary.LittleEndian, &data_buf)
	if err1 != nil {
		fmt.Println("read buf failed err is", err1)
		return nil, err1
	}


	message := &Message{}
	message.data = data_buf
	message.id = msg_id
	message.size = length

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
