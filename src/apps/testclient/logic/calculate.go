package logic

import (
	"ChatZoo/common"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FourOperationCalculate 根据标准输入向服务器发送 四则运算表达式运算请求
func FourOperationCalculate() common.IMessage {
	fmt.Println("已连接计算服务器,请输入你的四则运算公式, 空格分割, \\n 为结束符, 例如 [3 * 3 + 9]")
	//fmt.Scanln(&word) // 从标准控制中输入,以空格分隔
	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println("os.stdin read err ", inputErr)
		return nil
	}

	// 把字符串中的\r\n筛选出来
	words := ""
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		words += v // todo 好像这种写法很消耗,有新的写法
	}

	args, err := common.PackArgs(words)
	if err != nil {
		fmt.Println("pack args ", err)
		return nil
	}
	msg := &common.MsgCmdReq{
		MethodName: "Calculate",
		Args:       args,
	}
	return msg
}
