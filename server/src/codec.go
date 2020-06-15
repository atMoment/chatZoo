// 网络底层函数
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// decode from []byte to Message
// 编码， 将字节流转为message
func DDecode (data [] byte) (*Message, error) {
	// 此前已经读了size，这里不用读size了
	buf_reader := bytes.NewReader(data)

	data_size := int32(len(data))

	var msg_id int32
	err := binary.Read(buf_reader, binary.LittleEndian, &msg_id)
	if err != nil {
		fmt.Println("read id failed err is", err)
		return nil, err
	}

	message := &Message{}
	message.id = msg_id
	message.size = data_size + 4                                  // data的字节加上id字节和size字节

	data_buf := data[4:]                                          // 从id之后读取数据
	message.data, err = minfo.GetMessageInfo(msg_id)
	if err != nil {
		return nil, err
	}
	err = Decode(data_buf, message.data)
	if err != nil {
		fmt.Println("read buf failed err is", err)
		return nil, err
	}

	return message, nil
}

// encode from Message to []byte
// 解码，将消息转为字节流
// 写竟然可以不规定容器长度，读可以不规定容易长度吗
func EEncode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	tmp_buffer := new(bytes.Buffer)

	size, err := Encode(tmp_buffer, msg.data)
	if err != nil {
		return nil, err
	}

	msg.SetSize(int32(size + 8))            // 加上id4字节 + size4字节 + data长度字节, 用interface也可以这样算长度
	err = binary.Write(buffer, binary.LittleEndian, msg.GetSize())
	if err != nil {
		return nil, err
	}

	err =  binary.Write(buffer, binary.LittleEndian, msg.id)
	if err != nil {
		return nil, err
	}

	_, err = Encode(tmp_buffer, msg.data)
	if err != nil {
		return nil, err
	}
	tmp_data := tmp_buffer.Bytes()

	_, err = buffer.Write(tmp_data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil

}


