package main

import (
	"ChatZoo/common/rrand"
	"fmt"
	"time"
)

/*
话语接龙房间玩法：
所有玩家 ready 后游戏开始
游戏开始后, 从库里随机出一个字, 分发给房间内玩家
玩家在规定时间内用该字作为头部,组成词语/成语/俗语
后面的玩家将用前面玩家的内容的最后一个字组成头部, 定时客户端做不了现在
最后展示
未在规定时间内作答的玩家轮次将无法进行。有任一玩家下线, 游戏退出

	// 同时给大家出题
	// 所有人都答完(太复杂先不做)或者倒计时结束不许答题
	// 如果有人没有答,此条线路结束, 太复杂先不做

// 多人游戏 玩的是信息差
// todo 房间中有人离开了,如果正在进行某游戏流程,应该怎么做？ 这里简单处理为游戏结束, 还是需要动态调整让游戏继续玩下去？
*/
const (
	gameIntervalDuration = 5 * time.Minute
)

// 接龙房间
type _ChainRoom struct {
	IRoom
	readyMap map[string]int     // key: 已经准备好的选手 val: 当前第几道
	turns    [][]string         // key: 第几道 val: 预备参加的选手名单按次序
	first    []string           // key: 第几道, val: 第一个字
	records  map[int][]*_Record // key: 第几道, val: 选手及选手答案
	curTurn  int                // 当前轮次
	timer    *time.Timer
}

type _Record struct {
	actor   string
	content string
}

func (r *_ChainRoom) ready(player string) {
	// 所有人全部ready 游戏自动开始
	r.readyMap[player] = -1
	if len(r.readyMap) == r.GetRoomMemberLimit() {
		r.start()
	}
}

func (r *_ChainRoom) collect(player, content string) {
	if r.curTurn == len(r.first) {
		fmt.Println("轮次结束")
		return
	}
	d := r.readyMap[player]
	r.records[d] = append(r.records[d], &_Record{actor: player, content: content})
	if len(r.records[r.curTurn]) != r.GetRoomMemberLimit() {
		return
	}

	r.curTurn++
	for _, member := range r.turns[r.curTurn] {
		r.NotifyMember(member, "ChainGameTurnBegin", r.records[d])
		r.readyMap[member] = r.curTurn
	}

	if r.curTurn == r.GetRoomMemberLimit() {
		r.gameOver()
	}
}

func (r *_ChainRoom) start() {
	var err error
	r.first, err = getFirstKey(r.GetRoomMemberLimit())
	if err != nil {
		panic(fmt.Sprintf("setFirstKey err:%v", err))
	}
	r.turns, err = getPlayerTurn(r.GetRoomMemberList())
	if err != nil {
		panic(fmt.Sprintf("getPlayerTurn err:%v", err))
	}

	for i, member := range r.turns[0] {
		r.NotifyMember(member, "ChainGameTurnBegin", r.first[i])
		r.readyMap[member] = 0
	}
	r.timer = time.AfterFunc(gameIntervalDuration, r.gameOver)
	fmt.Printf("game start %v %v %v\n", r.first, r.turns[r.curTurn], r.curTurn)
}

func (r *_ChainRoom) gameOver() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.IRoom.NotifyAllMember("ChainGameOver")
	// 游戏结算,展示所有记录, 并发送给所有玩家
	// 现在服务器和客户端暂不支持map, 就先服务器自己看吧
}

// getFirstKey 先从池里随机出关键字
func getFirstKey(num int) ([]string, error) {
	allKey := []string{"春", "天", "地", "云", "光", "日", "月", "鸟", "花", "水"}
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

// getPlayerTurn 得到接龙的玩家的顺序
func getPlayerTurn(keyList []string) ([][]string, error) {
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
