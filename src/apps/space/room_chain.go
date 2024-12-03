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

const (
	ChainPhase_NotBegin = iota
	ChainPhase_Begin
	ChainPhase_End
)

type IChainRoom interface {
	Ready(player string)
	Collect(player, content string)
}

// 接龙房间
type _ChainRoom struct {
	IRoomBase
	readyMap map[string]struct{} // key: 已经准备好的选手
	turns    [][]string          // i: 某赛道的接棒选手们, j: 要接第几棒的选手
	first    []string            // key: 赛道, val: 第一个字
	records  [][]*_Record        // i 赛道  j:某一棒选手表现
	curTurn  int                 // 当前到第几棒
	timer    *time.Timer
	phase    int // 阶段
}

type _Record struct {
	actor   string
	content string
}

func NewChainRoom(limit int) *_ChainRoom {
	chainRoom := &_ChainRoom{
		IRoomBase: NewRoom(limit),
		readyMap:  make(map[string]struct{}),
	}
	chainRoom.SetType(RoomType_Chain)
	return chainRoom
}

// Ready 有玩家准备好时调用
func (room *_ChainRoom) Ready(player string) {
	if room.phase != ChainPhase_NotBegin {
		fmt.Println("phase not illegal")
		return
	}
	// 所有人全部ready 游戏自动开始
	room.readyMap[player] = struct{}{}
	if len(room.readyMap) == room.GetRoomMemberLimit() {
		room.start()
		room.phase = ChainPhase_Begin
	}
}

// Collect 每轮中玩家发言
func (room *_ChainRoom) Collect(player, content string) {
	if room.curTurn == room.GetRoomMemberLimit() {
		fmt.Println("轮次结束,不应该发言")
		return
	}
	// 找到这个选手是哪个赛道
	track := -1
	for i := 0; i < room.GetRoomMemberLimit(); i++ {
		if room.turns[i][room.curTurn] == player {
			track = i
			break
		}
	}
	if track == -1 {
		fmt.Println("player illegal ", player)
		return
	}

	room.records[track] = append(room.records[track], &_Record{actor: player, content: content})
	if !room.canReachNextLevel() {
		fmt.Printf("当前轮次:%v, 赛道:%v player:%v已交棒 :%v,等待其他队伍交棒\n", room.curTurn, track, player, content)
		return
	}
	fmt.Printf("当前轮次:%v, 赛道:%v player:%v已交棒 :%v,所有赛道已交棒,进入下一轮次\n", room.curTurn, track, player, content)

	room.curTurn++
	room.timer.Reset(gameIntervalDuration)
	if room.curTurn == room.GetRoomMemberLimit() {
		room.gameOver()
		return
	}

	for i := 0; i < room.GetRoomMemberLimit(); i++ {
		member := room.turns[i][room.curTurn]          // 各自赛道的第x棒选手
		key := room.records[i][len(room.records[i])-1] // 各自赛道的棒
		// 通知选手起跑
		room.NotifyMember(member, "SRPC_ChainGameTurnBegin", key)
		fmt.Printf("下一轮次开始,选手各就位！ 轮次:%v player:%v 关键词:%v\n", room.curTurn, member, key)
	}
}

func (room *_ChainRoom) start() {
	var err error
	room.first, err = getFirstKey(room.GetRoomMemberLimit())
	if err != nil {
		panic(fmt.Sprintf("setFirstKey err:%v", err))
	}
	room.turns, err = getPlayerTurn(room.GetRoomMemberList())
	if err != nil {
		panic(fmt.Sprintf("getPlayerTurn err:%v", err))
	}

	room.records = make([][]*_Record, len(room.turns))
	for i := 0; i < room.GetRoomMemberLimit(); i++ {
		member := room.turns[i][0] // 各自赛道的第0棒选手
		firstKey := room.first[i]  // 各自赛道的棒
		// 通知选手起跑
		room.NotifyMember(member, "SRPC_ChainGameTurnBegin", firstKey)
		room.records[i] = append(room.records[i], &_Record{
			actor:   "system",
			content: firstKey,
		})
	}
	room.timer = time.AfterFunc(gameIntervalDuration, room.dealTimeOut)
	fmt.Printf("game start firstKey:%v turn:%+v curTurn:%v\n", room.first, room.turns, room.curTurn)
}

func (room *_ChainRoom) dealTimeOut() {
	fmt.Printf("超时期限到,有玩家未回答,游戏结束\n")
	room.gameOver()
}

func (room *_ChainRoom) gameOver() {
	if room.timer != nil {
		room.timer.Stop()
		room.timer = nil
	}
	// 告诉所有玩家游戏结果
	room.NotifyAllMember("SPRC_ChainGameOver")
	fmt.Printf("game over, all result:%+v\n", room.records)
	// todo 触发房间销毁？
}

// canReachNextLevel 是否可以进入到下一棒
func (room *_ChainRoom) canReachNextLevel() bool {
	// 检查每一个赛道, 是否都已经交棒
	for i := 0; i < room.GetRoomMemberLimit(); i++ {
		if len(room.records[i]) != room.curTurn+2 {
			return false
		}
	}
	return true
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
