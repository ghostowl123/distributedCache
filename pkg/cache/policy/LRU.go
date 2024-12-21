package policy

import "container/list"

type LRU[K comparable, V any] struct {
	capacity int
	list     *list.List
	items    map[K]*list.Element
}

type Entry[K comparable, V any] struct {
	key   K
	value V
}

func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		capacity: capacity,
		list:     list.New(),
		items:    make(map[K]*list.Element),
	}
}

func (l *LRU[K, V]) Add(key K, value V) {
	if elem, exists := l.items[key]; exists {
		entry := elem.Value.(*Entry[K, V])
		entry.value = value
		l.list.MoveToFront(elem)
		return
	}

	entry := &Entry[K, V]{key: key, value: value}
	elem := l.list.PushFront(entry)
	l.items[key] = elem

	if l.list.Len() > l.capacity {
		l.Evict()
	}
}

func (l *LRU[K, V]) RecordAccess(key K, value V) {
	if elem, exist := l.items[key]; exist {
		entry := elem.Value.(*Entry[K, V])
		entry.value = value
		l.list.MoveToFront(elem)
	}
}

func (l *LRU[K, V]) Evict() (K, V, bool) {
	elem := l.list.Back()
	if elem == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	entry := elem.Value.(*Entry[K, V])
	l.list.Remove(elem)
	delete(l.items, entry.key)
	return entry.key, entry.value, true
}

func (l *LRU[K, V]) Remove(key K) (V, bool) {
	if elem, exist := l.items[key]; exist {
		entry := elem.Value.(*Entry[K, V])
		l.list.Remove(elem)
		delete(l.items, key)
		return entry.value, true
	}
	var zero V
	return zero, false
}
