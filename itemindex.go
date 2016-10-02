package cdw

import (
	"errors"
	"sync"
)

type ItemIndex struct {
	sync.Mutex
	cache map[string]int32
}

func NewItemIndex() *ItemIndex {
	return &ItemIndex{
		cache: make(map[string]int32),
	}
}

func (ii *ItemIndex) set(key string, value int) {
	ii.cache[key] = int32(value)
}

func (ii *ItemIndex) get(key string) (result int32, err error) {
	if ii.count() > 0 {
		result, ok := ii.cache[key]
		if !ok {
			return 0, errors.New("No item")
		}
		err = nil
		return result, err
	}
	return 0, errors.New("No Item")
}

func (ii *ItemIndex) count() int {
	return len(ii.cache)
}

func (ii *ItemIndex) Set(key string, value int) {
	ii.Lock()
	defer ii.Unlock()

	ii.set(key, value)
}

func (ii *ItemIndex) Get(key string) (result int32) {
	ii.Lock()
	defer ii.Unlock()

	result, err := ii.get(key)
	if err != nil {
		number := ii.count()
		result = int32(number)
		ii.set(key, number)
	}

	return
}

func (ii *ItemIndex) Count() int {
	ii.Lock()
	defer ii.Unlock()
	return ii.count()
}

func (ii *ItemIndex) GetValues() []string {
	ii.Lock()
	defer ii.Unlock()

	total := ii.count()
	results := make([]string, total)
	for k, v := range ii.cache {
		results[int(v)] = k
	}
	return results
}
