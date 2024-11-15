package encrypt

import (
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

/*
dh 算法
使用Pair 生成公钥、密钥
客户端和服务器双方交换公钥，Key(自己的密钥+对方的公钥) 得到 同一个值
*/

var (
	rnd        = rand.New(rand.NewSource(time.Now().UnixNano())) // 随机数种子
	dhBase     = big.NewInt(3)                                   //Diffie-Hellman 交换密钥算法的质数
	dhPrime, _ = big.NewInt(0).SetString("0x7FFFFFC3", 0)        // Diffie-Hellman 交换密钥算法的质数
	maxNum     = big.NewInt(math.MaxInt64)

	randLock = sync.Mutex{} // 加锁是因为 rand不安全吗？
)

// Pair 返回一对 私钥和匹配的公钥
func Pair() (*big.Int, *big.Int) {
	randLock.Lock()
	privateKey := big.NewInt(0).Rand(rnd, maxNum)
	randLock.Unlock()
	publicKey := big.NewInt(0).Exp(dhBase, privateKey, dhPrime)
	return privateKey, publicKey
}

// Key 私钥+公钥计算出的值
func Key(privateKey *big.Int, otherKey *big.Int) *big.Int {
	return big.NewInt(0).Exp(otherKey, privateKey, dhPrime)
}
