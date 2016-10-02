package cdw

import (
	"errors"
	"sync"
)

type LocationIndex struct {
	sync.Mutex
	cache map[string]int32
}

func NewLocationIndex() *LocationIndex {
	return &LocationIndex{
		cache: make(map[string]int32),
	}
}

func (li *LocationIndex) set(key string, value int) {
	li.cache[key] = int32(value)
}

func (li *LocationIndex) get(key string) (result int32, err error) {
	if li.count() > 0 {
		result, ok := li.cache[key]
		if !ok {
			return 0, errors.New("No item")
		}
		err = nil
		return result, err
	}
	return 0, errors.New("No Item")
}

func (li *LocationIndex) count() int {
	return len(li.cache)
}

func (li *LocationIndex) Set(key string, value int) {
	li.Lock()
	defer li.Unlock()

	li.set(key, value)
}

func (li *LocationIndex) Get(key string) (result int32) {
	li.Lock()
	defer li.Unlock()

	result, err := li.get(key)
	if err != nil {
		number := li.count()
		result = int32(number)
		li.set(key, number)
	}

	return
}

func (li *LocationIndex) Count() int {
	li.Lock()
	defer li.Unlock()
	return li.count()
}

func (li *LocationIndex) GetValues() []string {
	li.Lock()
	defer li.Unlock()

	total := li.count()
	results := make([]string, total)
	for k, v := range li.cache {
		results[int(v)] = k
	}
	return results
}
