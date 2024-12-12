package common

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

const (
	QueueMacLength = 100
)

type IEntityRpcQueue interface {
	Push(index int32, method reflect.Value, args []reflect.Value)
	Pop() ([]reflect.Value, int32)
}

// 很简陋的单协程队列, 保证每个函数执行时只有它自己, 目前满了直接丢掉
// 可以用任意容器做队列, 我自己用的链表, 先进先出

type _RpcQueue struct {
	firstNode *_RpcCell
	length    int
	lock      sync.Mutex
	cond      *sync.Cond
}

type _RpcCell struct {
	method reflect.Value
	args   []reflect.Value
	next   *_RpcCell
	Index  int32 // 消息号
}

func NewRpcQueue() IEntityRpcQueue {
	ret := &_RpcQueue{
		lock: sync.Mutex{},
	}
	ret.cond = sync.NewCond(&ret.lock)
	return ret
}

func (q *_RpcQueue) Push(index int32, method reflect.Value, args []reflect.Value) {
	// 插入的时候如果满了, 采取的方案是不再接收, 但这肯定不是最好的方案
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()
		q.cond.Signal()
	}()
	if q.length == QueueMacLength {
		fmt.Println("q.length is max")
		q.cond.Wait()
	}
	newNode := &_RpcCell{
		method: method,
		args:   args,
		Index:  index,
	}
	q.length++
	if q.length == 1 {
		q.firstNode = newNode
		return
	}

	n := q.firstNode
	for n.next != nil {
		n = n.next
	}
	n.next = newNode
}
func (q *_RpcQueue) Pop() ([]reflect.Value, int32) {
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()
		q.cond.Signal()
	}()
	if q.length == 0 {
		fmt.Println("q.length is 0")
		q.cond.Wait()
	}
	method := q.firstNode.method
	args := q.firstNode.args
	msgIndex := q.firstNode.Index

	q.firstNode = q.firstNode.next
	q.length--
	before := time.Now()
	rets := method.Call(args)
	after := time.Now()
	if after.Sub(before).Milliseconds() > time.Duration(1*time.Second).Milliseconds() {
		fmt.Println("exec func too slow ", after.Sub(before).Milliseconds())
	}
	return rets, msgIndex
}
