package main

const Request_join = 1
const Response_join = 2
const Response_chat = 3
const Request_chat = 4


type ReqJoin struct {
	Room_id int32
}

type ResJoin struct {
	Words string
}

type ReqChat struct {
	Words string
}

type ResChat struct {
	Words string
}