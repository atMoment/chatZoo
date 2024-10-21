package main

import (
	"ChatZoo/common"
	"ChatZoo/testclient/logic"
)

type ModuleFunc func(string) common.IMessage

var moduleFuncList map[string]ModuleFunc

func init() {
	moduleFuncList = make(map[string]ModuleFunc)
	moduleFuncList[logic.ModuleNameChat] = logic.Chat
	moduleFuncList[logic.ModuleNameFourOperationCalculate] = logic.FourOperationCalculate
}
