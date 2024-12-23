package policy

import (
	"container/heap"
)

type Item[V any] struct {
	priority int
	value    V
	index    int
}

// PriorityQueue Implementation
type PriorityQueue[V any] []*Item[V]

func (pq PriorityQueue[V]) Len() int { return len(pq) }

func (pq PriorityQueue[V]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[V]) Push(x any) {
	task := x.(*Item[V])
	task.index = len(*pq)
	*pq = append(*pq, task)
}

func (pq *PriorityQueue[V]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[V]) update(item *Item[V], value V, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// LFU implementation
// PriorityQueue + Hashmap
type LFU[K comparable, V any] struct {
	capacity      int
	priorityqueue PriorityQueue[*LFUEntry[K, V]]
	items         map[K]*LFUEntry[K, V]
}
type LFUEntry[K comparable, V any] struct {
	key      K
	value    V
	priority int
	index    int
}

func NewLFU[K comparable, V any](capacity int) *LFU[K, V] {
	pq := make(PriorityQueue[*LFUEntry[K, V]], 0)
	heap.Init(&pq)
	return &LFU[K, V]{
		capacity:      capacity,
		priorityqueue: pq,
		items:         make(map[K]*LFUEntry[K, V]),
	}
}
func (l *LFU[K, V]) Add(key K, value V) {
	if entry, exist := l.items[key]; exist {
		// Update value and increment priority
		item := &Item[*LFUEntry[K, V]]{
			value:    entry,
			priority: entry.priority + 1,
			index:    entry.index,
		}
		l.priorityqueue.update(item, item.value, item.priority+1)
		return
	}

	// If the key does not exist, handle capacity
	if len(l.items) >= l.capacity {
		l.Evict()
	}

	// Add the new entry
	newItem := &LFUEntry[K, V]{
		key:      key,
		value:    value,
		priority: 1, // New items start with priority 1
	}
	heap.Push(&l.priorityqueue, newItem)
	l.items[key] = newItem
}

func (l *LFU[K, V]) RecordAccess(key K, value V) {
	if entry, exist := l.items[key]; exist {
		// Update value and increment priority
		item := &Item[*LFUEntry[K, V]]{
			value:    entry,
			priority: entry.priority + 1,
			index:    entry.index,
		}
		l.priorityqueue.update(item, item.value, item.priority+1)
	}
}

func (l *LFU[K, V]) Evict() (K, V, bool) {
	if len(l.items) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	item := heap.Pop(&l.priorityqueue).(*LFUEntry[K, V])
	delete(l.items, item.key)
	return item.key, item.value, true
}
func (l *LFU[K, V]) Remove(key K) (V, bool) {
	if entry, exist := l.items[key]; exist {
		heap.Remove(&l.priorityqueue, entry.index)
		delete(l.items, entry.key)
		return entry.value, true
	}
	var zeroV V
	return zeroV, false
}
