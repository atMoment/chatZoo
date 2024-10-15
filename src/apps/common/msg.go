package common

import (
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrReadFull  = errors.New("io.ReadFull err")
	ErrConnWrite = errors.New("conn write err")
	ErrDecode    = errors.New("decode err")
	ErrEncode    = errors.New("encode err")
)

// ReadFromConn 当前只能处理string
func ReadFromConn(conn net.Conn) (string, error) {
	sizeBuf := make([]byte, 4) // int32 占4个字节,2 ^ 32, 可容纳4G的数据量, 非常够用了
	var err error
	_, err = io.ReadFull(conn, sizeBuf)
	if err != nil {
		fmt.Println("io.ReadFull read size err ", err)
		return "", ErrReadFull
	}
	var size int32 // 把总长度读出来
	err = Decode(sizeBuf, &size)
	if err != nil {
		fmt.Println("decode []byte to type err ", err)
		return "", ErrDecode
	}
	infoBuf := make([]byte, size)
	_, err = io.ReadFull(conn, infoBuf) // 能接着读
	if err != nil {
		fmt.Println("io.ReadFull read info err ", err)
		return "", ErrReadFull
	}
	// todo  这里可以换成任何结构
	var words string
	err = Decode(infoBuf, &words)
	if err != nil {
		fmt.Println("io.ReadFull read info err ", err)
		return "", ErrReadFull
	}

	return words, nil
}

func WriteToConn(conn net.Conn, words string) error {
	// todo 这里可以换成任何结构
	dataSize, data, err := Encode(words)
	if err != nil {
		fmt.Printf("words:%v encode err:%v \n", words, err)
		return ErrEncode
	}

	_, size, err := Encode(int32(dataSize))
	if err != nil {
		fmt.Printf("words size:%v encode err:%v \n", dataSize, err)
		return ErrEncode
	}
	allData := make([]byte, 0)
	allData = append(size, data...)
	_, err = conn.Write(allData)
	if err != nil {
		return ErrConnWrite
	}
	return nil
}
