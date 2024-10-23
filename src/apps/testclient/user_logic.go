package main

import (
	mmsg "ChatZoo/common/msg"
	"ChatZoo/testclient/logic"
)

type ModuleFunc func(string) mmsg.IMessage

var moduleFuncList map[string]ModuleFunc

func init() {
	moduleFuncList = make(map[string]ModuleFunc)
	moduleFuncList[logic.ModuleNameChat] = logic.Chat
	moduleFuncList[logic.ModuleNameFourOperationCalculate] = logic.FourOperationCalculate
}
