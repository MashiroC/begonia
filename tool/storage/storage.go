package storage

import "sync"

var std map[string]interface{}
var list map[string][]interface{}
var stdL sync.Mutex
var listL sync.Mutex

func init() {
	std = make(map[string]interface{})
	list = make(map[string][]interface{})
}

// Store 给本地存储添加数据，并发安全，但是新的会覆盖旧的
func Store(key string, value interface{}) {
	stdL.Lock()
	defer stdL.Unlock()

	std[key] = value
}

// Load 取走数据，会删除key对应的value
func Load(key string) interface{} {
	stdL.Lock()
	defer stdL.Unlock()

	v, ok := std[key]
	if ok {
		delete(std, key)
	}
	return v
}

// NewList 创建新的List，如果key相同会覆盖旧的
func NewList(key string) {
	listL.Lock()
	defer listL.Unlock()

	list[key] = make([]interface{}, 4, 4)
}

// DelteList 从List中删除一个key
func DeleteList(key string) {
	listL.Lock()
	defer listL.Unlock()

	delete(list, key)
}

// Push 给List尾部添加一个新的数据
func Push(key string, value interface{}) {
	listL.Lock()
	defer listL.Unlock()

	list[key] = append(list[key], value)
}

// Pull 一次性取出key中所有数据
func Pull(key string) []interface{} {
	listL.Lock()
	defer listL.Unlock()

	v, ok := list[key]
	if ok {
		list[key] = list[key][:0]
	}

	return v
}
