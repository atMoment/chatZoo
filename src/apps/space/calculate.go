package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func calculate(s string) (int, error) {
	factorList := make([]int, 0)
	operatorMap := map[string]struct{}{"+": {}, "-": {}, "*": {}, "/": {}}
	operatorList := make([]string, 0)
	for _, w := range strings.Split(s, " ") {
		if w == "\n" || w == "\r" {
			continue
		}
		if _, ok := operatorMap[w]; ok {
			operatorList = append(operatorList, w)
			continue
		}
		n, err := strconv.Atoi(w)
		if err != nil {
			return 0, errors.New("不能识别的符号")
		}
		factorList = append(factorList, n)
	}
	final := loopCalculate(operatorList, factorList)
	fmt.Println(final)
	return final, nil
}

func loopCalculate(operatorList []string, factorList []int) int {
	for {
		if !dealAdvanceOperator(operatorList, factorList) {
			break
		}
	}
	for {
		if !dealNormalOperator(operatorList, factorList) {
			break
		}
	}
	return factorList[0]
}

func dealAdvanceOperator(operatorList []string, factorList []int) bool {
	for i, v := range operatorList {
		switch v {
		case "*":
			newNum := factorList[i] * factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return true
		case "/":
			newNum := factorList[i] / factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return true
		}
	}
	return false
}

func dealNormalOperator(operatorList []string, factorList []int) bool {
	for i, v := range operatorList {
		switch v {
		case "+":
			newNum := factorList[i] * factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return true
		case "-":
			newNum := factorList[i] / factorList[i+1]
			operatorList = append(operatorList[:i], operatorList[i+1:]...)
			factorList[i+1] = newNum
			factorList = append(factorList[:i], factorList[i+1:]...)
			return true
		}
	}
	return false
}
