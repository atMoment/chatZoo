package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ModuleNameFourOperationCalculate = "FourOperationCalculate"
)

// FourOperationCalculate 根据标准输入向服务器发送 四则运算表达式运算请求
func (u *_User) FourOperationCalculate() (string, string) {
	fmt.Println("已连接计算服务器,请输入你的四则运算公式, 空格分割, \\n 为结束符, 例如 [3 * 3 + 9]")
	//fmt.Scanln(&word) // 从标准控制中输入,以空格分隔
	inputReader := bufio.NewReader(os.Stdin)
	input, inputErr := inputReader.ReadString('\n') // 回车
	if inputErr != nil {
		fmt.Println("os.stdin read err ", inputErr)
		return "", ""
	}

	// 把字符串中的\r\n筛选出来
	words := ""
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		words += v // todo 好像这种写法很消耗,有新的写法
	}

	return "Calculate", words
}

// dealInput  // 把字符串中的\r\n筛选出来
func dealInput(input string) string {
	splitStrings := strings.Split(input, "\r\n")
	words := make([]string, 0, len(splitStrings))
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		words = append(words, v)
	}
	return strings.Join(words, "")
}

// dealInput2 方法2  // 把字符串中的\r\n筛选出来
func dealInput2(input string) string {
	words := ""
	for _, v := range strings.Split(input, "\r\n") {
		if v == "\r\n" {
			break
		}
		words += v // todo 好像这种写法很消耗,有新的写法
	}
	return words
}

// dealInput3 方法3  // 把字符串中的\r\n筛选出来
func dealInput3(input string) string {
	// 因为\r\n是在最末尾, \r\n 是转义字符,应该有码表示, 发现到特别的码, 就删掉
	// 但不知道怎么写
	return ""
}
