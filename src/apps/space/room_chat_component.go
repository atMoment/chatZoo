package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	ComponentChat = "ComponentChat"
)

type IChatComponent interface {
	Chat(member, content string) error
	GetHistoryMsg() []*_ChatComponentMsg
}

type _ChatComponent struct {
	IRoomBase
	msgCache []*_ChatComponentMsg // 消息缓存
	cfg      *_ChatComponentCfg
}

func NewChatComponent(room IRoomBase) *_ChatComponent {
	ChatComponent := &_ChatComponent{
		IRoomBase: room,
		msgCache:  make([]*_ChatComponentMsg, 0),
	}
	return ChatComponent
}

type _ChatComponentCfg struct {
	isSave              bool // 是否存档
	cacheMsgLimit       int8 // 缓存(内存)消息上限
	persistenceMsgLimit int8 // 持久化消息上限
}

type _ChatComponentMsg struct {
	fromID   string
	fromName string
	content  string
	sendTime int64
}

func (r *_ChatComponent) Chat(member, content string) error {
	if !r.MemberIsExist(member) {
		return errors.New("member not exist in room")
	}

	if len(r.msgCache) >= int(r.cfg.cacheMsgLimit) {
		r.msgCache = r.msgCache[r.cfg.cacheMsgLimit/2:] // 前面的直接丢掉
	}
	r.msgCache = append(r.msgCache, &_ChatComponentMsg{
		fromID:   member,
		content:  content,
		sendTime: time.Now().UnixNano(),
	})
	r.NotifyAllMember("Notify_SToCMessage", fmt.Sprintf("%v say: %v", member, content))
	return nil
}

func (r *_ChatComponent) GetHistoryMsg() []*_ChatComponentMsg {
	return r.msgCache
}
