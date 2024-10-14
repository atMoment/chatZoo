package common

import (
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	ErrReadFull = errors.New("io.ReadFull err")
	ErrDecode   = errors.New("decode err")
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
	size := uint32(0)
	err = Decode(sizeBuf, &size)
	if err != nil {
		fmt.Println("decode []byte to type err ", err)
		return "", ErrDecode
	}
	infoBuf := make([]byte, size)
	_, err = io.ReadFull(conn, infoBuf)
	if err != nil {
		fmt.Println("io.ReadFull read info err ", err)
		return "", ErrReadFull
	}
	return string(infoBuf), nil
}
