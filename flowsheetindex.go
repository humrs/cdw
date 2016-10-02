package cdw

import (
	"errors"
	"sync"
)

type FlowsheetIndex struct {
	sync.Mutex
	cache map[string]int32
}

func NewFlowsheetIndex() *FlowsheetIndex {
	return &FlowsheetIndex{
		cache: make(map[string]int32),
	}
}
func (fi *FlowsheetIndex) set(key string, value int) {
	fi.cache[key] = int32(value)
}

func (fi *FlowsheetIndex) get(key string) (result int32, err error) {
	if fi.count() > 0 {
		result, ok := fi.cache[key]
		if !ok {
			return 0, errors.New("No item")
		}
		err = nil
		return result, err
	}
	return 0, errors.New("No Item")
}

func (fi *FlowsheetIndex) count() int {
	return len(fi.cache)
}

func (fi *FlowsheetIndex) Set(key string, value int) {
	fi.Lock()
	defer fi.Unlock()

	fi.set(key, value)
}

func (fi *FlowsheetIndex) Get(key string) (result int32) {
	fi.Lock()
	defer fi.Unlock()

	result, err := fi.get(key)
	if err != nil {
		number := fi.count()
		result = int32(number)
		fi.set(key, number)
	}

	return
}

func (fi *FlowsheetIndex) Count() int {
	fi.Lock()
	defer fi.Unlock()
	return fi.count()
}

func (fi *FlowsheetIndex) GetValues() []string {
	fi.Lock()
	defer fi.Unlock()

	total := fi.count()
	results := make([]string, total)
	for k, v := range fi.cache {
		results[int(v)] = k
	}
	return results
}
