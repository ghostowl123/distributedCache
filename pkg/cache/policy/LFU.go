package policy

import "container/list"

type LFU[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
}


