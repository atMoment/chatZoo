package main

import (
	"errors"
	"fmt"
	"time"
)

type _ChatRoom struct {
	IRoom
	msgCache []*_ChatRoomMsg // 消息缓存
	cfg      *_ChatRoomCfg
}

type _ChatRoomCfg struct {
	isSave              bool // 是否存档
	cacheMsgLimit       int8 // 缓存(内存)消息上限
	persistenceMsgLimit int8 // 持久化消息上限
}

type _ChatRoomMsg struct {
	fromID   string
	fromName string
	content  string
	sendTime int64
}

func (r *_ChatRoom) chat(member, content string) error {
	if !r.MemberIsExist(member) {
		return errors.New("member not exist in room")
	}

	if len(r.msgCache) >= int(r.cfg.cacheMsgLimit) {
		r.msgCache = r.msgCache[r.cfg.cacheMsgLimit/2:] // 前面的直接丢掉
	}
	r.msgCache = append(r.msgCache, &_ChatRoomMsg{
		fromID:   member,
		content:  content,
		sendTime: time.Now().UnixNano(),
	})
	r.NotifyAllMember("Notify_SToCMessage", fmt.Sprintf("%v say: %v", member, content))
	return nil
}

func (r *_ChatRoom) getHistoryMsg() []*_ChatRoomMsg {
	return r.msgCache
}
