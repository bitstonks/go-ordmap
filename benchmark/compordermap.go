// DO NOT EDIT this code was generated using go-ordmap code generation
// go run github.com/edofic/go-ordmap/cmd/gen --pkg ordmap --name=CompOrderMap --key=comp --value=*order --target=compordermap_test.go --less=.Less
//lint:file-ignore ST1000 ignore problems in generated code
package ordmap

type CompOrderMapEntry struct {
	K comp
	V *order
}

type CompOrderMap struct {
	CompOrderMapEntry CompOrderMapEntry
	h                 int
	len               int
	children          [2]*CompOrderMap
}

func (node *CompOrderMap) Height() int {
	if node == nil {
		return 0
	}
	return node.h
}

// suffix CompOrderMap is needed because this will get specialised in codegen
func combinedDepthCompOrderMap(n1, n2 *CompOrderMap) int {
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

// suffix CompOrderMap is needed because this will get specialised in codegen
func mkCompOrderMap(entry CompOrderMapEntry, left *CompOrderMap, right *CompOrderMap) *CompOrderMap {
	len := 1
	if left != nil {
		len += left.len
	}
	if right != nil {
		len += right.len
	}
	return &CompOrderMap{
		CompOrderMapEntry: entry,
		h:                 combinedDepthCompOrderMap(left, right),
		len:               len,
		children:          [2]*CompOrderMap{left, right},
	}
}

func (node *CompOrderMap) Get(key comp) (value *order, ok bool) {
	finger := node
	for {
		if finger == nil {
			ok = false
			return // using named returns so we keep the zero value for `value`
		}
		if key.Less(finger.CompOrderMapEntry.K) {
			finger = finger.children[0]
		} else if finger.CompOrderMapEntry.K.Less(key) {
			finger = finger.children[1]
		} else {
			// equal
			return finger.CompOrderMapEntry.V, true
		}
	}
}

func (node *CompOrderMap) Insert(key comp, value *order) *CompOrderMap {
	if node == nil {
		return mkCompOrderMap(CompOrderMapEntry{key, value}, nil, nil)
	}
	entry, left, right := node.CompOrderMapEntry, node.children[0], node.children[1]
	if node.CompOrderMapEntry.K.Less(key) {
		right = right.Insert(key, value)
	} else if key.Less(node.CompOrderMapEntry.K) {
		left = left.Insert(key, value)
	} else { // equals
		entry = CompOrderMapEntry{key, value}
	}
	return rotateCompOrderMap(entry, left, right)
}

func (node *CompOrderMap) Remove(key comp) *CompOrderMap {
	if node == nil {
		return nil
	}
	entry, left, right := node.CompOrderMapEntry, node.children[0], node.children[1]
	if node.CompOrderMapEntry.K.Less(key) {
		right = right.Remove(key)
	} else if key.Less(node.CompOrderMapEntry.K) {
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
	return rotateCompOrderMap(entry, left, right)
}

// suffix CompOrderMap is needed because this will get specialised in codegen
func rotateCompOrderMap(entry CompOrderMapEntry, left *CompOrderMap, right *CompOrderMap) *CompOrderMap {
	if right.Height()-left.Height() > 1 { // implies right != nil
		// single left
		rl := right.children[0]
		rr := right.children[1]
		if combinedDepthCompOrderMap(left, rl)-rr.Height() > 1 {
			// double rotation
			return mkCompOrderMap(
				rl.CompOrderMapEntry,
				mkCompOrderMap(entry, left, rl.children[0]),
				mkCompOrderMap(right.CompOrderMapEntry, rl.children[1], rr),
			)
		}
		return mkCompOrderMap(right.CompOrderMapEntry, mkCompOrderMap(entry, left, rl), rr)
	}
	if left.Height()-right.Height() > 1 { // implies left != nil
		// single right
		ll := left.children[0]
		lr := left.children[1]
		if combinedDepthCompOrderMap(right, lr)-ll.Height() > 1 {
			// double rotation
			return mkCompOrderMap(
				lr.CompOrderMapEntry,
				mkCompOrderMap(left.CompOrderMapEntry, ll, lr.children[0]),
				mkCompOrderMap(entry, lr.children[1], right),
			)
		}
		return mkCompOrderMap(left.CompOrderMapEntry, ll, mkCompOrderMap(entry, lr, right))
	}
	return mkCompOrderMap(entry, left, right)
}

func (node *CompOrderMap) Len() int {
	if node == nil {
		return 0
	}
	return node.len
}

func (node *CompOrderMap) Entries() []CompOrderMapEntry {
	elems := make([]CompOrderMapEntry, 0, node.Len())
	if node == nil {
		return elems
	}
	type frame struct {
		node     *CompOrderMap
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
			elems = append(elems, top.node.CompOrderMapEntry)
			if top.node.children[1] != nil {
				stack = append(stack, frame{top.node.children[1], false})
			}
		}
	}
	return elems
}

func (node *CompOrderMap) extreme(dir int) *CompOrderMapEntry {
	if node == nil {
		return nil
	}
	finger := node
	for finger.children[dir] != nil {
		finger = finger.children[dir]
	}
	return &finger.CompOrderMapEntry
}

func (node *CompOrderMap) Min() *CompOrderMapEntry {
	return node.extreme(0)
}

func (node *CompOrderMap) Max() *CompOrderMapEntry {
	return node.extreme(1)
}

func (node *CompOrderMap) Iterate() CompOrderMapIterator {
	return newIteratorCompOrderMap(node, 0, nil)
}

func (node *CompOrderMap) IterateFrom(k comp) CompOrderMapIterator {
	return newIteratorCompOrderMap(node, 0, &k)
}

func (node *CompOrderMap) IterateReverse() CompOrderMapIterator {
	return newIteratorCompOrderMap(node, 1, nil)
}

func (node *CompOrderMap) IterateReverseFrom(k comp) CompOrderMapIterator {
	return newIteratorCompOrderMap(node, 1, &k)
}

type CompOrderMapIteratorStackFrame struct {
	node  *CompOrderMap
	state int8
}

type CompOrderMapIterator struct {
	direction    int
	stack        []CompOrderMapIteratorStackFrame
	currentEntry CompOrderMapEntry
}

// suffix CompOrderMap is needed because this will get specialised in codegen
func newIteratorCompOrderMap(node *CompOrderMap, direction int, startFrom *comp) CompOrderMapIterator {
	if node == nil {
		return CompOrderMapIterator{}
	}
	stack := make([]CompOrderMapIteratorStackFrame, 1, node.Height())
	stack[0] = CompOrderMapIteratorStackFrame{node: node, state: 0}
	iter := CompOrderMapIterator{direction: direction, stack: stack}
	if startFrom != nil {
		stack[0].state = 2
		iter.seek(*startFrom)
	} else {
		iter.Next()
	}
	return iter
}

func (i *CompOrderMapIterator) Done() bool {
	return len(i.stack) == 0
}

func (i *CompOrderMapIterator) GetKey() comp {
	return i.currentEntry.K
}

func (i *CompOrderMapIterator) GetValue() *order {
	return i.currentEntry.V
}

func (i *CompOrderMapIterator) Next() {
	for len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		switch frame.state {
		case 0:
			if frame.node == nil {
				last := len(i.stack) - 1
				i.stack[last] = CompOrderMapIteratorStackFrame{} // zero out
				i.stack = i.stack[:last]                         // pop
			} else {
				frame.state = 1
			}
		case 1:
			i.stack = append(i.stack, CompOrderMapIteratorStackFrame{node: frame.node.children[i.direction], state: 0})
			frame.state = 2
		case 2:
			i.currentEntry = frame.node.CompOrderMapEntry
			frame.state = 3
			return
		case 3:
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = CompOrderMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 0}
		}

	}
}

func (i *CompOrderMapIterator) seek(k comp) {
LOOP:
	for {
		frame := &i.stack[len(i.stack)-1]
		if frame.node == nil {
			last := len(i.stack) - 1
			i.stack[last] = CompOrderMapIteratorStackFrame{} // zero out
			i.stack = i.stack[:last]                         // pop
			break LOOP
		}
		if (i.direction == 0 && !(frame.node.CompOrderMapEntry.K.Less(k))) || (i.direction == 1 && !(k.Less(frame.node.CompOrderMapEntry.K))) {
			i.stack = append(i.stack, CompOrderMapIteratorStackFrame{node: frame.node.children[i.direction], state: 2})
		} else {
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = CompOrderMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 2}
		}
	}
	if len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		i.currentEntry = frame.node.CompOrderMapEntry
		frame.state = 3
	}
}
