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
	Push(index int32, method reflect.Value, args []reflect.Value) bool
	Pop() ([]reflect.Value, int32, bool)
	Close()
}

// 目前满了会阻塞, 有优化空间
// 可以用任意容器做队列, 我自己用的链表, 先进先出

type _RpcQueue struct {
	firstNode *_RpcCell
	length    int
	close     bool
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

// Push 返回值为true, 表示已经关闭了, 不再push
func (q *_RpcQueue) Push(index int32, method reflect.Value, args []reflect.Value) bool {
	if q.close { // 如果关闭了就不push
		return true
	}
	// 插入的时候如果满了, 采取的方案是不再接收, 但这肯定不是最好的方案
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()
		q.cond.Signal()
	}()
	for q.length == QueueMacLength && !q.close {
		fmt.Println("q.length is max")
		q.cond.Wait()
	}
	if q.length < QueueMacLength {
		q.queuePush(index, method, args)
		return q.close
	}
	return true
}

// Pop 返回true, 表示是pop完队列中最后一个
func (q *_RpcQueue) Pop() ([]reflect.Value, int32, bool) {
	q.cond.L.Lock()
	defer func() {
		q.cond.L.Unlock()
		q.cond.Signal()
	}()
	for q.length == 0 && !q.close {
		q.cond.Wait()
	}
	if q.length > 0 {
		return q.queuePop()
	}
	return nil, 0, true
}

func (q *_RpcQueue) Close() {
	q.cond.L.Lock()
	q.close = true
	q.cond.L.Unlock()
	q.cond.Signal()
}

func (q *_RpcQueue) queuePush(index int32, method reflect.Value, args []reflect.Value) {
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
func (q *_RpcQueue) queuePop() ([]reflect.Value, int32, bool) {
	method := q.firstNode.method
	args := q.firstNode.args
	msgIndex := q.firstNode.Index

	if msgIndex == HeartBeatIndex {
		return nil, msgIndex, false
	}

	q.firstNode = q.firstNode.next
	q.length--
	before := time.Now()
	rets := method.Call(args)
	after := time.Now()
	if after.Sub(before).Milliseconds() > time.Duration(2000*time.Second).Milliseconds() {
		fmt.Println("exec func too slow ", after.Sub(before).Milliseconds())
	}
	return rets, msgIndex, false
}
