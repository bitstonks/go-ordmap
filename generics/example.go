package main

import "fmt"

type i int

func (v i) Less(v2 i) bool {
	return int(v) < int(v2)
}

type Lesser[E any] interface {
	Less(E) bool
}

type Entry[K Lesser[K], V any] struct {
	K K
	V V
}

type OrdMap[K Lesser[K], V any] struct {
	Entry    Entry[K, V]
	h        int
	len      int
	children [2]*OrdMap[K, V]
}

func (node *OrdMap[K, V]) Height() int {
	if node == nil {
		return 0
	}
	return node.h
}

func combinedDepth[K Lesser[K], V any](n1, n2 *OrdMap[K, V]) int {
	d1 := n1.Height()
	d2 := n2.Height()
	var d int
	if d1 > d2 {
		d = d1
	} else {
		d = d2
	}
	return d + 1
}

func mk_OrdMap[K Lesser[K], V any](entry Entry[K, V], left, right *OrdMap[K, V]) *OrdMap[K, V] {
	len := 1
	if left != nil {
		len += left.len
	}
	if right != nil {
		len += right.len
	}
	return &OrdMap[K, V]{
		Entry:    entry,
		h:        combinedDepth(left, right),
		len:      len,
		children: [2]*OrdMap[K, V]{left, right},
	}
}

func (node *OrdMap[K, V]) Get(key K) (value V, ok bool) {
	finger := node
	for {
		if finger == nil {
			ok = false
			return // using named returns so we keep the zero value for `value`
		}
		if key.Less(finger.Entry.K) {
			finger = finger.children[0]
		} else if finger.Entry.K.Less(key) {
			finger = finger.children[1]
		} else {
			// equal
			return finger.Entry.V, true
		}
	}
}

func (node *OrdMap[K, V]) Insert(key K, value V) *OrdMap[K, V] {
	if node == nil {
		return mk_OrdMap(Entry[K, V]{key, value}, nil, nil)
	}
	entry, left, right := node.Entry, node.children[0], node.children[1]
	if node.Entry.K.Less(key) {
		right = right.Insert(key, value)
	} else if key.Less(node.Entry.K) {
		left = left.Insert(key, value)
	} else { // equals
		entry = Entry[K, V]{key, value}
	}
	return rotate(entry, left, right)
}

func (node *OrdMap[K, V]) Remove(key K) *OrdMap[K, V] {
	if node == nil {
		return nil
	}
	entry, left, right := node.Entry, node.children[0], node.children[1]
	if node.Entry.K.Less(key) {
		right = right.Remove(key)
	} else if key.Less(node.Entry.K) {
		left = left.Remove(key)
	} else { // equals
		max := left.Max()
		if max == nil {
			return right
		} else {
			left = left.Remove(max.K)
			entry = *max
		}
	}
	return rotate(entry, left, right)
}

func rotate[K Lesser[K], V any](entry Entry[K, V], left, right *OrdMap[K, V]) *OrdMap[K, V] {
	if right.Height()-left.Height() > 1 { // implies right != nil
		// single left
		rl := right.children[0]
		rr := right.children[1]
		if combinedDepth(left, rl)-rr.Height() > 1 {
			// double rotation
			return mk_OrdMap(
				rl.Entry,
				mk_OrdMap(entry, left, rl.children[0]),
				mk_OrdMap(right.Entry, rl.children[1], rr),
			)
		}
		return mk_OrdMap(right.Entry, mk_OrdMap(entry, left, rl), rr)
	}
	if left.Height()-right.Height() > 1 { // implies left != nil
		// single right
		ll := left.children[0]
		lr := left.children[1]
		if combinedDepth(right, lr)-ll.Height() > 1 {
			// double rotation
			return mk_OrdMap(
				lr.Entry,
				mk_OrdMap(left.Entry, ll, lr.children[0]),
				mk_OrdMap(entry, lr.children[1], right),
			)
		}
		return mk_OrdMap(left.Entry, ll, mk_OrdMap(entry, lr, right))
	}
	return mk_OrdMap(entry, left, right)
}

func (node *OrdMap[K, V]) Len() int {
	if node == nil {
		return 0
	}
	return node.len
}

type entriesFrame[K Lesser[K], V any] struct {
	node     *OrdMap[K, V]
	leftDone bool
}

func (node *OrdMap[K, V]) Entries() []Entry[K, V] {
	elems := make([]Entry[K, V], 0, node.Len())
	if node == nil {
		return elems
	}
	var preallocated [20]entriesFrame[K, V] // preallocate on stack for common case
	stack := preallocated[:0]
	stack = append(stack, entriesFrame[K, V]{node, false})
	for len(stack) > 0 {
		top := &stack[len(stack)-1]

		if !top.leftDone {
			if top.node.children[0] != nil {
				stack = append(stack, entriesFrame[K, V]{top.node.children[0], false})
			}
			top.leftDone = true
		} else {
			stack = stack[:len(stack)-1] // pop
			elems = append(elems, top.node.Entry)
			if top.node.children[1] != nil {
				stack = append(stack, entriesFrame[K, V]{top.node.children[1], false})
			}
		}
	}
	return elems
}

func (node *OrdMap[K, V]) extreme(dir int) *Entry[K, V] {
	if node == nil {
		return nil
	}
	finger := node
	for finger.children[dir] != nil {
		finger = finger.children[dir]
	}
	return &finger.Entry
}

func (node *OrdMap[K, V]) Min() *Entry[K, V] {
	return node.extreme(0)
}

func (node *OrdMap[K, V]) Max() *Entry[K, V] {
	return node.extreme(1)
}

func main() {
	fmt.Println("hello")
	var m1 *OrdMap[i, string] // zero value is the empty map
	fmt.Println(m1.Entries())
	m1 = m1.Insert(1, "foo") // adding entries
	fmt.Println(m1.Entries())
	m1 = m1.Insert(2, "baz")
	fmt.Println(m1.Entries())
	m1 = m1.Insert(2, "bar") // will override
	fmt.Println(m1.Entries())
	fmt.Println(m1.Get(2)) // access by key
	m1 = m1.Insert(3, "baz")
	fmt.Println(m1.Entries())
	m1 = m1.Remove(1) // can also remove entries
	fmt.Println(m1.Entries())
}
