package utils

import (
	"encoding/json"
	"sync"
	"time"
)

type Statistics struct {
	mtx   sync.Mutex
	items map[string]int

	//itemsLastTimes map[string] time.Time
	itemsCountTimes map[string]int
}

var stat Statistics

func init() {
	stat.items = make(map[string]int)
	stat.itemsCountTimes = make(map[string]int)

	go threadStat()
}

func threadStat() {
	for true {
		stat.mtx.Lock()
		for key, value := range stat.itemsCountTimes {
			stat.items[key] = value
			stat.itemsCountTimes[key] = 0
		}
		stat.mtx.Unlock()

		time.Sleep(1 * time.Second)
	}
}

func Inc(name string) {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()
	if val, ok := stat.items[name]; ok {
		stat.items[name] = val + 1
	} else {
		stat.items[name] = 1
	}
}

func Write(name string, value int) {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()
	stat.items[name] = value
}

func Add(name string, value int) {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()

	if _, ok := stat.items[name]; !ok {
		stat.items[name] = value
	} else {
		stat.items[name] += value
	}
}

func Reg(name string) {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()

	if val, ok := stat.itemsCountTimes[name]; ok {
		stat.itemsCountTimes[name] = val + 1
	} else {
		stat.itemsCountTimes[name] = 1
	}
}

func RegAdd(name string, value int) {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()

	if val, ok := stat.itemsCountTimes[name]; ok {
		stat.itemsCountTimes[name] = val + value
	} else {
		stat.itemsCountTimes[name] = value
	}
}

func Get() string {
	stat.mtx.Lock()
	defer stat.mtx.Unlock()

	b, _ := json.MarshalIndent(stat.items, " ", " ")
	return string(b)
}
