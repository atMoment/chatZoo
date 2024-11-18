package main

import "testing"

func Test_RandomGuessList(t *testing.T) {
	players := []string{"tomato", "cucumber", "pumpkin"}
	ret, err := getPlayerTurn(players)
	t.Logf("ret %+v err:%v", ret, err)
}
