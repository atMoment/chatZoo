package main

const Request_join = 1
const Response_join = 2
const Response_chat = 3
const Request_chat = 4


type RequestJoin struct {
	room_id int32
}

type ResponseJoin struct {
}

type RequestChat struct {
	words string
}

type ResponseChat struct {
	words string
}