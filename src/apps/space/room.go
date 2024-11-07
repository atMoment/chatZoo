package main

import (
	"ChatZoo/common"
	"ChatZoo/common/rrand"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 一些特别简单的狗屎代码
// todo 房间中有人离开了,如果正在进行某游戏流程,应该怎么做？
const (
	_ = iota
	chat
	guess
)

var roomMgr = &_RoomMgr{}

type _RoomMgr struct {
	rooms sync.Map // key: roomID val: room
	typ   int
}

func (mgr *_RoomMgr) AddEntity(userID string) {
	room := &_Room{
		createTime: time.Now().UnixNano(),
	}
	mgr.rooms.Store(userID, room)
}

func (mgr *_RoomMgr) AddOrGetEntity(userID string) (*_Room, error) {
	room := &_Room{
		createTime: time.Now().UnixNano(),
		memberList: make(map[string]struct{}),
	}
	// 找不到 返回 false
	entityInfo, ok := mgr.rooms.LoadOrStore(userID, room)
	if !ok {
		return room, nil
	}

	ret, transOk := entityInfo.(*_Room)
	if !transOk {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) GetEntity(userID string) (*_Room, error) {
	entityInfo, ok := mgr.rooms.Load(userID)
	if !ok {
		return nil, errors.New("userid not find")
	}
	ret, ok := entityInfo.(*_Room)
	if !ok {
		return nil, errors.New("trans userinfo err")
	}
	return ret, nil
}

func (mgr *_RoomMgr) DeleteEntity(userID string) {
	mgr.rooms.Delete(userID)
}

type _Room struct {
	memberList map[string]struct{}
	limit      uint8 // 最多255人
	createTime int64
	msgCache   []*_RoomChatMsg // 消息缓存
}

type _RoomChatMsg struct {
	fromID   string
	fromName string
	content  string
	sendTime int64
}

func (r *_Room) joinRoom(member string) {
	r.memberList[member] = struct{}{}
}

func (r *_Room) quitRoom(member string) {
	delete(r.memberList, member)
}

func (r *_Room) chat(member, memberName, content string) {
	m := fmt.Sprintf("%v say: %v", member, content)
	common.DefaultSrvEntity.SendNotifyToEntityList(r.memberList, "Notify_SToCMessage", m)
}

/*
话语接龙房间玩法：
所有玩家 ready 后游戏开始
游戏开始后, 从库里随机出一个字, 分发给房间内玩家
玩家在规定时间内用该字作为头部,组成词语/成语/俗语
后面的玩家将用前面玩家的内容的最后一个字组成头部
最后展示
未在规定时间内作答的玩家轮次将无法进行。有任一玩家下线, 游戏退出
// 多人游戏 玩的是信息差
*/

type _guessGame struct {
	room       *_Room
	memberList []string
	turns      [][]string          // key: 第几道 val: 预备参加的选手名单按次序
	first      []string            // key: 第几道, val: 第一个字
	records    map[int][]*_record  // key: 第几道, val: 选手及选手答案
	readyMap   map[string]struct{} // key: 已经准备好的选手
	limit      int                 // 房间总共多少人
}

type _record struct {
	actor   string
	content string
}

func (r *_guessGame) ready(player string) {
	// 所有人全部ready 游戏自动开始
	r.readyMap[player] = struct{}{}
	if len(r.readyMap) == r.limit {
		r.start()
	}
}

func (r *_guessGame) start() {
	//先从池里随机出关键字

}

func (r *_guessGame) guess() {
	var err error
	r.first, err = r.getFirstKey(r.limit)
	if err != nil {
		panic(fmt.Sprintf("setFirstKey err:%v", err))
	}
	r.turns, err = r.getPlayerTurn(r.memberList)

	// 同时给大家出题
	// 所有人都答完(太复杂先不做)或者倒计时结束不许答题
	// 如果有人没有答,此条线路结束
}

func (r *_guessGame) Update() {

}

func (r *_guessGame) getFirstKey(num int) ([]string, error) {
	allKey := []string{"番茄", "天", "地", "云", "光", "日", "月", "好", "猫", "狗"}
	pool := rrand.NewRandomWeight[string, int]()
	for _, v := range allKey {
		pool.Add(v, 1)
	}
	keylist, err := pool.RandomMultiple(num, true)
	if err != nil {
		panic(fmt.Sprintf("setFirstKey err:%v", err))
	}
	return keylist, err
}

func (r *_guessGame) getPlayerTurn(keyList []string) ([][]string, error) {
	pool := rrand.NewRandomWeight[string, int]()

	all := make([][]string, len(keyList))
	for i := 0; i < len(keyList); i++ {
		switch i {
		case 0:
			for _, key := range keyList {
				pool.Add(key, 1)
			}
			list, err := pool.RandomMultiple(len(keyList), true)
			if err != nil {
				fmt.Println("first err ", err)
				return nil, err
			}
			all[i] = list
		default:
			pool.Clean()
			for _, key := range keyList {
				pool.Add(key, 1)
			}
			for index := 0; index < len(keyList); index++ {
				for i_index := i - 1; i_index >= 0; i_index-- {
					pool.Delete(all[i_index][index])
				}

				element, randomErr := pool.Random()
				if randomErr != nil {
					fmt.Println("random err ", randomErr)
					return nil, randomErr
				}
				all[i] = append(all[i], element)
				for i_index := i - 1; i_index >= 0; i_index-- {
					pool.Add(all[i_index][index], 1)
				}
				pool.Delete(element)
			}
		}
	}
	return all, nil
}
