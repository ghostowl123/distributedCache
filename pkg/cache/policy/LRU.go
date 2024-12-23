package policy

import "container/list"

//doubly linked list + hashmap (key: key, value: node reference)
type LRU[K comparable, V any] struct {
	capacity int
	list     *list.List
	items    map[K]*list.Element
}

type LRUEntry[K comparable, V any] struct {
	key   K
	value V
}

//constructor
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		capacity: capacity,
		list:     list.New(),
		items:    make(map[K]*list.Element),
	}
}

// add new key-value pair into list
func (l *LRU[K, V]) Add(key K, value V) {
	if elem, exists := l.items[key]; exists {
		entry := elem.Value.(*LRUEntry[K, V])
		entry.value = value
		l.list.MoveToFront(elem)
		return
	}

	entry := &LRUEntry[K, V]{key: key, value: value}
	elem := l.list.PushFront(entry)
	l.items[key] = elem

	if l.list.Len() > l.capacity {
		l.Evict()
	}
}

//update the sequence
func (l *LRU[K, V]) RecordAccess(key K, value V) {
	if elem, exist := l.items[key]; exist {
		entry := elem.Value.(*LRUEntry[K, V])
		entry.value = value
		l.list.MoveToFront(elem)
	}
}

// invalid the cache
func (l *LRU[K, V]) Evict() (K, V, bool) {
	elem := l.list.Back()
	if elem == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	entry := elem.Value.(*LRUEntry[K, V])
	l.list.Remove(elem)
	delete(l.items, entry.key)
	return entry.key, entry.value, true
}

func (l *LRU[K, V]) Remove(key K) (V, bool) {
	if elem, exist := l.items[key]; exist {
		entry := elem.Value.(*LRUEntry[K, V])
		l.list.Remove(elem)
		delete(l.items, key)
		return entry.value, true
	}
	var zero V
	return zero, false
}
