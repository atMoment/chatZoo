package common

// 先读放在最前面的类型, 再根据类型占多少字节读后面的数据

// UnpackArgs 将[]byte 转化具体类型再转化为为[]interface{}
func UnpackArgs(data []byte) []interface{} {
	if len(data) == 0 {
		return nil
	}
}
