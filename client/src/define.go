package main

const Request_join = 1
const Response_join = 2
const Response_chat = 3
const Request_chat = 4


type ReqJoin struct {
	room_id int32
}

type ResJoin struct {
}

type ReqChat struct {
	words string
}

type ResChat struct {
	words string
}