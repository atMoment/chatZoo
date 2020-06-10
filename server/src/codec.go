// 网络底层函数
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// decode from []byte to Message
// 编码， 将字节流转为message
func Decode (data [] byte) (*Message, error) {
	// 此前已经读了size，这里不用读size了
	buf_reader := bytes.NewReader(data)

	data_size := int32(len(data))

	var msg_id int32
	err0 := binary.Read(buf_reader, binary.LittleEndian, &msg_id)
	if err0 != nil {
		fmt.Println("read id failed err is", err0)
		return nil, err0
	}

	length := data_size - 4                                      //减去刚刚读的id字节
	data_buf := make([]byte, length)
	err1 := binary.Read(buf_reader, binary.LittleEndian, &data_buf)
	if err1 != nil {
		fmt.Println("read buf failed err is", err1)
		return nil, err1
	}


	message := &Message{}
	message.data = data_buf
	message.id = msg_id
	message.size = length + 8                                    // data的字节加上id字节和size字节

	return message, nil
}

// encode from Message to []byte
// 解码，将消息转为字节流
// 写竟然可以不规定容器长度，读可以不规定容易长度吗
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
