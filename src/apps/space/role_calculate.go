package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (r *_User) Calculate(expression string) string {
	var ret string
	result, err := calculate(expression)
	if err != nil {
		ret = fmt.Sprintf("session server calculate, err:%v", err)
	} else {
		ret = fmt.Sprintf("%s = %d", expression, result)
	}
	fmt.Printf("Calculate success expression:%v, ret:%v\n ", expression, ret)
	return ret
}

func calculate(s string) (int, error) {
	operatorList, factorList, err := analysisInput(s)
	if err != nil {
		return 0, fmt.Errorf("分析输入失败, err: %v", err)
	}
	final, err := loopCalculate2(operatorList, factorList)
	if err != nil {
		return 0, fmt.Errorf("计算失败, err: %v", err)
	}
	return final, nil
}

// analysisInput 从字符串中分析出数值和运算符
func analysisInput(input string) ([]string, []int, error) {
	operatorMap := map[string]struct{}{"+": {}, "-": {}, "*": {}, "/": {}}

	factorList := make([]int, 0)
	operatorList := make([]string, 0)
	for _, word := range strings.Split(input, " ") {
		w := strings.Split(word, "\r")[0]
		// 为了兼容 输完后正常公式后再输入一个空格。 例如 [3 * 3 ] => ["3", "*", "3", "\r\n"]]
		// 不输入空格 [3 * 3] => ["3", "*", "3\r\n"]
		if len(w) == 0 {
			continue
		}
		if _, ok := operatorMap[w]; ok {
			operatorList = append(operatorList, w)
			continue
		}
		n, err := strconv.Atoi(w)
		if err != nil {
			return nil, nil, fmt.Errorf("不能识别的数字或者运算符:%v, err:%v \n", w, err)
		}
		factorList = append(factorList, n)
	}
	if len(operatorList) != len(factorList)-1 {
		return nil, nil, errors.New("输入公式有误, 请检查")
	}
	return operatorList, factorList, nil
}

// loopCalculate2 为正确写法, loopCalculate 是slice 离谱的浅拷贝写法, 欢迎去回顾
func loopCalculate2(operatorList []string, factorList []int) (int, error) {
	operatorListCopy := operatorList
	factorListCopy := factorList
	var ok bool
	var err error
	for {
		operatorListCopy, factorListCopy, ok, err = dealAdvanceOperator2(operatorListCopy, factorListCopy)
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
	}
	for {
		operatorListCopy, factorListCopy, ok, err = dealNormalOperator2(operatorListCopy, factorListCopy)
		if err != nil {
			return 0, err
		}
		if !ok {
			break
		}
	}
	return factorListCopy[0], nil
}

func dealAdvanceOperator2(operatorListCopy []string, factorListCopy []int) ([]string, []int, bool, error) {
	operatorList := make([]string, len(operatorListCopy))
	factorList := make([]int, len(factorListCopy))
	copy(operatorList, operatorListCopy)
	copy(factorList, factorListCopy)

	for i, v := range operatorList {
		switch v {
		case "*":
			newNum := factorList[i] * factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i] = newNum
			factorList = append(factorList[:i+1], factorList[i+2:]...)
			return operatorList, factorList, true, nil
		case "/":
			if factorList[i+1] == 0 {
				return nil, nil, false, errors.New("不可除0")
			}
			newNum := factorList[i] / factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i] = newNum
			factorList = append(factorList[:i+1], factorList[i+2:]...)
			return operatorList, factorList, true, nil
		}
	}
	return operatorList, factorList, false, nil
}

func dealNormalOperator2(operatorListCopy []string, factorListCopy []int) ([]string, []int, bool, error) {
	operatorList := make([]string, len(operatorListCopy))
	factorList := make([]int, len(factorListCopy))
	copy(operatorList, operatorListCopy)
	copy(factorList, factorListCopy)

	for i, v := range operatorList {
		switch v {
		case "+":
			newNum := factorList[i] + factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return operatorList, factorList, true, nil
		case "-":
			newNum := factorList[i] - factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return operatorList, factorList, true, nil
		}
	}
	return operatorList, factorList, false, nil
}
