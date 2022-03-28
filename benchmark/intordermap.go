// DO NOT EDIT this code was generated using go-ordmap code generation
// go run github.com/edofic/go-ordmap/cmd/gen --pkg ordmap --name=IntOrderMap --key=int --value=*order --target=intordermap_test.go --less=<
//lint:file-ignore ST1000 ignore problems in generated code
package ordmap

type IntOrderMapEntry struct {
	K int
	V *order
}

type IntOrderMap struct {
	IntOrderMapEntry IntOrderMapEntry
	h                int
	len              int
	children         [2]*IntOrderMap
}

func (node *IntOrderMap) Height() int {
	if node == nil {
		return 0
	}
	return node.h
}

// suffix IntOrderMap is needed because this will get specialised in codegen
func combinedDepthIntOrderMap(n1, n2 *IntOrderMap) int {
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

// suffix IntOrderMap is needed because this will get specialised in codegen
func mkIntOrderMap(entry IntOrderMapEntry, left *IntOrderMap, right *IntOrderMap) *IntOrderMap {
	len := 1
	if left != nil {
		len += left.len
	}
	if right != nil {
		len += right.len
	}
	return &IntOrderMap{
		IntOrderMapEntry: entry,
		h:                combinedDepthIntOrderMap(left, right),
		len:              len,
		children:         [2]*IntOrderMap{left, right},
	}
}

func (node *IntOrderMap) Get(key int) (value *order, ok bool) {
	finger := node
	for {
		if finger == nil {
			ok = false
			return // using named returns so we keep the zero value for `value`
		}
		if key < (finger.IntOrderMapEntry.K) {
			finger = finger.children[0]
		} else if finger.IntOrderMapEntry.K < (key) {
			finger = finger.children[1]
		} else {
			// equal
			return finger.IntOrderMapEntry.V, true
		}
	}
}

func (node *IntOrderMap) Insert(key int, value *order) *IntOrderMap {
	if node == nil {
		return mkIntOrderMap(IntOrderMapEntry{key, value}, nil, nil)
	}
	entry, left, right := node.IntOrderMapEntry, node.children[0], node.children[1]
	if node.IntOrderMapEntry.K < (key) {
		right = right.Insert(key, value)
	} else if key < (node.IntOrderMapEntry.K) {
		left = left.Insert(key, value)
	} else { // equals
		entry = IntOrderMapEntry{key, value}
	}
	return rotateIntOrderMap(entry, left, right)
}

func (node *IntOrderMap) Remove(key int) *IntOrderMap {
	if node == nil {
		return nil
	}
	entry, left, right := node.IntOrderMapEntry, node.children[0], node.children[1]
	if node.IntOrderMapEntry.K < (key) {
		right = right.Remove(key)
	} else if key < (node.IntOrderMapEntry.K) {
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
	return rotateIntOrderMap(entry, left, right)
}

// suffix IntOrderMap is needed because this will get specialised in codegen
func rotateIntOrderMap(entry IntOrderMapEntry, left *IntOrderMap, right *IntOrderMap) *IntOrderMap {
	if right.Height()-left.Height() > 1 { // implies right != nil
		// single left
		rl := right.children[0]
		rr := right.children[1]
		if combinedDepthIntOrderMap(left, rl)-rr.Height() > 1 {
			// double rotation
			return mkIntOrderMap(
				rl.IntOrderMapEntry,
				mkIntOrderMap(entry, left, rl.children[0]),
				mkIntOrderMap(right.IntOrderMapEntry, rl.children[1], rr),
			)
		}
		return mkIntOrderMap(right.IntOrderMapEntry, mkIntOrderMap(entry, left, rl), rr)
	}
	if left.Height()-right.Height() > 1 { // implies left != nil
		// single right
		ll := left.children[0]
		lr := left.children[1]
		if combinedDepthIntOrderMap(right, lr)-ll.Height() > 1 {
			// double rotation
			return mkIntOrderMap(
				lr.IntOrderMapEntry,
				mkIntOrderMap(left.IntOrderMapEntry, ll, lr.children[0]),
				mkIntOrderMap(entry, lr.children[1], right),
			)
		}
		return mkIntOrderMap(left.IntOrderMapEntry, ll, mkIntOrderMap(entry, lr, right))
	}
	return mkIntOrderMap(entry, left, right)
}

func (node *IntOrderMap) Len() int {
	if node == nil {
		return 0
	}
	return node.len
}

func (node *IntOrderMap) Entries() []IntOrderMapEntry {
	elems := make([]IntOrderMapEntry, 0, node.Len())
	if node == nil {
		return elems
	}
	type frame struct {
		node     *IntOrderMap
		leftDone bool
	}
	var preallocated [20]frame // preallocate on stack for common case
	stack := preallocated[:0]
	stack = append(stack, frame{node, false})
	for len(stack) > 0 {
		top := &stack[len(stack)-1]

		if !top.leftDone {
			if top.node.children[0] != nil {
				stack = append(stack, frame{top.node.children[0], false})
			}
			top.leftDone = true
		} else {
			stack = stack[:len(stack)-1] // pop
			elems = append(elems, top.node.IntOrderMapEntry)
			if top.node.children[1] != nil {
				stack = append(stack, frame{top.node.children[1], false})
			}
		}
	}
	return elems
}

func (node *IntOrderMap) extreme(dir int) *IntOrderMapEntry {
	if node == nil {
		return nil
	}
	finger := node
	for finger.children[dir] != nil {
		finger = finger.children[dir]
	}
	return &finger.IntOrderMapEntry
}

func (node *IntOrderMap) Min() *IntOrderMapEntry {
	return node.extreme(0)
}

func (node *IntOrderMap) Max() *IntOrderMapEntry {
	return node.extreme(1)
}

func (node *IntOrderMap) Iterate() IntOrderMapIterator {
	return newIteratorIntOrderMap(node, 0, nil)
}

func (node *IntOrderMap) IterateFrom(k int) IntOrderMapIterator {
	return newIteratorIntOrderMap(node, 0, &k)
}

func (node *IntOrderMap) IterateReverse() IntOrderMapIterator {
	return newIteratorIntOrderMap(node, 1, nil)
}

func (node *IntOrderMap) IterateReverseFrom(k int) IntOrderMapIterator {
	return newIteratorIntOrderMap(node, 1, &k)
}

type IntOrderMapIteratorStackFrame struct {
	node  *IntOrderMap
	state int8
}

type IntOrderMapIterator struct {
	direction    int
	stack        []IntOrderMapIteratorStackFrame
	currentEntry IntOrderMapEntry
}

// suffix IntOrderMap is needed because this will get specialised in codegen
func newIteratorIntOrderMap(node *IntOrderMap, direction int, startFrom *int) IntOrderMapIterator {
	if node == nil {
		return IntOrderMapIterator{}
	}
	stack := make([]IntOrderMapIteratorStackFrame, 1, node.Height())
	stack[0] = IntOrderMapIteratorStackFrame{node: node, state: 0}
	iter := IntOrderMapIterator{direction: direction, stack: stack}
	if startFrom != nil {
		stack[0].state = 2
		iter.seek(*startFrom)
	} else {
		iter.Next()
	}
	return iter
}

func (i *IntOrderMapIterator) Done() bool {
	return len(i.stack) == 0
}

func (i *IntOrderMapIterator) GetKey() int {
	return i.currentEntry.K
}

func (i *IntOrderMapIterator) GetValue() *order {
	return i.currentEntry.V
}

func (i *IntOrderMapIterator) Next() {
	for len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		switch frame.state {
		case 0:
			if frame.node == nil {
				last := len(i.stack) - 1
				i.stack[last] = IntOrderMapIteratorStackFrame{} // zero out
				i.stack = i.stack[:last]                        // pop
			} else {
				frame.state = 1
			}
		case 1:
			i.stack = append(i.stack, IntOrderMapIteratorStackFrame{node: frame.node.children[i.direction], state: 0})
			frame.state = 2
		case 2:
			i.currentEntry = frame.node.IntOrderMapEntry
			frame.state = 3
			return
		case 3:
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = IntOrderMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 0}
		}

	}
}

func (i *IntOrderMapIterator) seek(k int) {
LOOP:
	for {
		frame := &i.stack[len(i.stack)-1]
		if frame.node == nil {
			last := len(i.stack) - 1
			i.stack[last] = IntOrderMapIteratorStackFrame{} // zero out
			i.stack = i.stack[:last]                        // pop
			break LOOP
		}
		if (i.direction == 0 && !(frame.node.IntOrderMapEntry.K < (k))) || (i.direction == 1 && !(k < (frame.node.IntOrderMapEntry.K))) {
			i.stack = append(i.stack, IntOrderMapIteratorStackFrame{node: frame.node.children[i.direction], state: 2})
		} else {
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = IntOrderMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 2}
		}
	}
	if len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		i.currentEntry = frame.node.IntOrderMapEntry
		frame.state = 3
	}
}
