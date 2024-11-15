package rrand

import (
	"errors"
	"fmt"
	"math/rand"
)

type _Weight interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type _RandomWeight[K comparable, W _Weight] struct {
	keylist    []K
	weightlist []W
}

func NewRandomWeight[K comparable, W _Weight]() *_RandomWeight[K, W] {
	return &_RandomWeight[K, W]{
		keylist:    make([]K, 0, 1),
		weightlist: make([]W, 0, 1),
	}
}

func (r *_RandomWeight[K, W]) Add(key K, weight W) {
	for _, k := range r.keylist {
		if k == key {
			return
		}
	}
	r.keylist = append(r.keylist, key)
	r.weightlist = append(r.weightlist, weight)
}

func (r *_RandomWeight[K, W]) Delete(key K) {
	for i, k := range r.keylist {
		if k == key {
			r.keylist = append(r.keylist[:i], r.keylist[i+1:]...)
			r.weightlist = append(r.weightlist[:i], r.weightlist[i+1:]...)
			return
		}
	}
}

func (r *_RandomWeight[K, W]) Clean() {
	r.keylist = make([]K, 0, 1)
	r.weightlist = make([]W, 0, 1)
}
func (r *_RandomWeight[K, W]) Random() (K, error) {
	var result K
	if len(r.keylist) == 0 {
		return result, errors.New("keylist is empty")
	}
	var total W
	for _, v := range r.weightlist {
		total += v
	}
	if total == 0 {
		return result, errors.New("weight total is 0")
	}

	var num W
	randomNum := W(rand.Float64() * float64(total))
	for i := 0; i < len(r.weightlist); i++ {
		num += r.weightlist[i]
		if randomNum < num {
			return r.keylist[i], nil
		}
	}
	return result, fmt.Errorf("can't find suit weight, randNum:%v", randomNum)
}
func (r *_RandomWeight[K, W]) RandomMultiple(times int, isNoRepeated bool) ([]K, error) {
	if times == 0 {
		return nil, errors.New("times is 0")
	}
	if isNoRepeated && times > len(r.keylist) {
		return nil, errors.New("no repeated Random, times > len(keylist)")
	}
	if isNoRepeated {
		cloneTarget := NewRandomWeight[K, W]()
		cloneTarget.clone(r)
		return cloneTarget.randomMultiple(times, true)
	} else {
		return r.randomMultiple(times, false)
	}
}

func (r *_RandomWeight[K, W]) GetKeyListNum() int {
	return len(r.keylist)
}

func (r *_RandomWeight[K, W]) randomMultiple(times int, isNoRepeated bool) ([]K, error) {
	ret := make([]K, 0)
	var result K
	var err error
	for i := 0; i < times; i++ {
		result, err = r.Random()
		if isNoRepeated {
			r.Delete(result)
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, result)
	}
	return ret, nil
}

func (r *_RandomWeight[K, W]) clone(target *_RandomWeight[K, W]) {
	r.keylist = make([]K, len(target.keylist))
	r.weightlist = make([]W, len(target.weightlist))

	copy(r.keylist, target.keylist)
	copy(r.weightlist, target.weightlist)
}
