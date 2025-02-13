package msg

import (
	"fmt"
	"testing"
)

func Test_DecodeAndEncode(t *testing.T) {
	var temp, temp_copy []bool
	temp = []bool{true, false}
	size, data, err := Encode(temp)
	if err != nil {
		fmt.Printf("encode err:%v\n", err)
		return
	}
	if err = Decode(data, &temp_copy); err != nil {
		fmt.Printf("decode err:%v\n", err)
		return
	}
	fmt.Printf("size:%v final result:%v \n", size, temp_copy)
}
