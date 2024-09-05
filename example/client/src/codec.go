package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)


// decode from []byte to Message
// 读
func DDecode (data [] byte) (*Message, error) {
	buf_reader := bytes.NewReader(data)

	data_size := int32(len(data))

	var msg_id int32
	err := binary.Read(buf_reader, binary.LittleEndian, &msg_id)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.id = msg_id
	message.size = data_size + 4

	data_buf := data[4:]
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
// 写
func EEncode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)
	tmp_buffer := new(bytes.Buffer)

	size, err := Encode(tmp_buffer, msg.data)

	if err != nil {
		return nil, err
	}

	msg.SetSize(int32(size + 8))                                           // 加上id4字节 + size4字节 + data长度字节
	err = binary.Write(buffer, binary.LittleEndian, msg.GetSize())
	if err != nil {
		return nil, err
	}

	err =  binary.Write(buffer, binary.LittleEndian, msg.id)
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