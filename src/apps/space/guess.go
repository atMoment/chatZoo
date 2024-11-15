package main

import (
	"ChatZoo/common/rrand"
	"fmt"
)

/*
算法简述： 游戏开始, 服务器同时发题给ABC三个玩家, 三个玩家答完后, 都进入到下一轮, 答不同的题
例如有玩家 A、B、C
A -> B  -> C
B -> C  -> A
C -> A  -> B
*/

// drawGuessRandom
func drawGuessRandom(keyList []string) ([][]string, error) {
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
