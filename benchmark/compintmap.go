// DO NOT EDIT this code was generated using go-ordmap code generation
// go run github.com/edofic/go-ordmap/cmd/gen --pkg ordmap --name=CompIntMap --key=comp --value=int --target=compintmap_test.go --less=.Less
//lint:file-ignore ST1000 ignore problems in generated code
package ordmap

type CompIntMapEntry struct {
	K comp
	V int
}

type CompIntMap struct {
	CompIntMapEntry CompIntMapEntry
	h               int
	len             int
	children        [2]*CompIntMap
}

func (node *CompIntMap) Height() int {
	if node == nil {
		return 0
	}
	return node.h
}

// suffix CompIntMap is needed because this will get specialised in codegen
func combinedDepthCompIntMap(n1, n2 *CompIntMap) int {
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

// suffix CompIntMap is needed because this will get specialised in codegen
func mkCompIntMap(entry CompIntMapEntry, left *CompIntMap, right *CompIntMap) *CompIntMap {
	len := 1
	if left != nil {
		len += left.len
	}
	if right != nil {
		len += right.len
	}
	return &CompIntMap{
		CompIntMapEntry: entry,
		h:               combinedDepthCompIntMap(left, right),
		len:             len,
		children:        [2]*CompIntMap{left, right},
	}
}

func (node *CompIntMap) Get(key comp) (value int, ok bool) {
	finger := node
	for {
		if finger == nil {
			ok = false
			return // using named returns so we keep the zero value for `value`
		}
		if key.Less(finger.CompIntMapEntry.K) {
			finger = finger.children[0]
		} else if finger.CompIntMapEntry.K.Less(key) {
			finger = finger.children[1]
		} else {
			// equal
			return finger.CompIntMapEntry.V, true
		}
	}
}

func (node *CompIntMap) Insert(key comp, value int) *CompIntMap {
	if node == nil {
		return mkCompIntMap(CompIntMapEntry{key, value}, nil, nil)
	}
	entry, left, right := node.CompIntMapEntry, node.children[0], node.children[1]
	if node.CompIntMapEntry.K.Less(key) {
		right = right.Insert(key, value)
	} else if key.Less(node.CompIntMapEntry.K) {
		left = left.Insert(key, value)
	} else { // equals
		entry = CompIntMapEntry{key, value}
	}
	return rotateCompIntMap(entry, left, right)
}

func (node *CompIntMap) Remove(key comp) *CompIntMap {
	if node == nil {
		return nil
	}
	entry, left, right := node.CompIntMapEntry, node.children[0], node.children[1]
	if node.CompIntMapEntry.K.Less(key) {
		right = right.Remove(key)
	} else if key.Less(node.CompIntMapEntry.K) {
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
	return rotateCompIntMap(entry, left, right)
}

// suffix CompIntMap is needed because this will get specialised in codegen
func rotateCompIntMap(entry CompIntMapEntry, left *CompIntMap, right *CompIntMap) *CompIntMap {
	if right.Height()-left.Height() > 1 { // implies right != nil
		// single left
		rl := right.children[0]
		rr := right.children[1]
		if combinedDepthCompIntMap(left, rl)-rr.Height() > 1 {
			// double rotation
			return mkCompIntMap(
				rl.CompIntMapEntry,
				mkCompIntMap(entry, left, rl.children[0]),
				mkCompIntMap(right.CompIntMapEntry, rl.children[1], rr),
			)
		}
		return mkCompIntMap(right.CompIntMapEntry, mkCompIntMap(entry, left, rl), rr)
	}
	if left.Height()-right.Height() > 1 { // implies left != nil
		// single right
		ll := left.children[0]
		lr := left.children[1]
		if combinedDepthCompIntMap(right, lr)-ll.Height() > 1 {
			// double rotation
			return mkCompIntMap(
				lr.CompIntMapEntry,
				mkCompIntMap(left.CompIntMapEntry, ll, lr.children[0]),
				mkCompIntMap(entry, lr.children[1], right),
			)
		}
		return mkCompIntMap(left.CompIntMapEntry, ll, mkCompIntMap(entry, lr, right))
	}
	return mkCompIntMap(entry, left, right)
}

func (node *CompIntMap) Len() int {
	if node == nil {
		return 0
	}
	return node.len
}

func (node *CompIntMap) Entries() []CompIntMapEntry {
	elems := make([]CompIntMapEntry, 0, node.Len())
	if node == nil {
		return elems
	}
	type frame struct {
		node     *CompIntMap
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
			elems = append(elems, top.node.CompIntMapEntry)
			if top.node.children[1] != nil {
				stack = append(stack, frame{top.node.children[1], false})
			}
		}
	}
	return elems
}

func (node *CompIntMap) extreme(dir int) *CompIntMapEntry {
	if node == nil {
		return nil
	}
	finger := node
	for finger.children[dir] != nil {
		finger = finger.children[dir]
	}
	return &finger.CompIntMapEntry
}

func (node *CompIntMap) Min() *CompIntMapEntry {
	return node.extreme(0)
}

func (node *CompIntMap) Max() *CompIntMapEntry {
	return node.extreme(1)
}

func (node *CompIntMap) Iterate() CompIntMapIterator {
	return newIteratorCompIntMap(node, 0, nil)
}

func (node *CompIntMap) IterateFrom(k comp) CompIntMapIterator {
	return newIteratorCompIntMap(node, 0, &k)
}

func (node *CompIntMap) IterateReverse() CompIntMapIterator {
	return newIteratorCompIntMap(node, 1, nil)
}

func (node *CompIntMap) IterateReverseFrom(k comp) CompIntMapIterator {
	return newIteratorCompIntMap(node, 1, &k)
}

type CompIntMapIteratorStackFrame struct {
	node  *CompIntMap
	state int8
}

type CompIntMapIterator struct {
	direction    int
	stack        []CompIntMapIteratorStackFrame
	currentEntry CompIntMapEntry
}

// suffix CompIntMap is needed because this will get specialised in codegen
func newIteratorCompIntMap(node *CompIntMap, direction int, startFrom *comp) CompIntMapIterator {
	if node == nil {
		return CompIntMapIterator{}
	}
	stack := make([]CompIntMapIteratorStackFrame, 1, node.Height())
	stack[0] = CompIntMapIteratorStackFrame{node: node, state: 0}
	iter := CompIntMapIterator{direction: direction, stack: stack}
	if startFrom != nil {
		stack[0].state = 2
		iter.seek(*startFrom)
	} else {
		iter.Next()
	}
	return iter
}

func (i *CompIntMapIterator) Done() bool {
	return len(i.stack) == 0
}

func (i *CompIntMapIterator) GetKey() comp {
	return i.currentEntry.K
}

func (i *CompIntMapIterator) GetValue() int {
	return i.currentEntry.V
}

func (i *CompIntMapIterator) Next() {
	for len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		switch frame.state {
		case 0:
			if frame.node == nil {
				last := len(i.stack) - 1
				i.stack[last] = CompIntMapIteratorStackFrame{} // zero out
				i.stack = i.stack[:last]                       // pop
			} else {
				frame.state = 1
			}
		case 1:
			i.stack = append(i.stack, CompIntMapIteratorStackFrame{node: frame.node.children[i.direction], state: 0})
			frame.state = 2
		case 2:
			i.currentEntry = frame.node.CompIntMapEntry
			frame.state = 3
			return
		case 3:
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = CompIntMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 0}
		}

	}
}

func (i *CompIntMapIterator) seek(k comp) {
LOOP:
	for {
		frame := &i.stack[len(i.stack)-1]
		if frame.node == nil {
			last := len(i.stack) - 1
			i.stack[last] = CompIntMapIteratorStackFrame{} // zero out
			i.stack = i.stack[:last]                       // pop
			break LOOP
		}
		if (i.direction == 0 && !(frame.node.CompIntMapEntry.K.Less(k))) || (i.direction == 1 && !(k.Less(frame.node.CompIntMapEntry.K))) {
			i.stack = append(i.stack, CompIntMapIteratorStackFrame{node: frame.node.children[i.direction], state: 2})
		} else {
			// override frame - tail call optimisation
			i.stack[len(i.stack)-1] = CompIntMapIteratorStackFrame{node: frame.node.children[1-i.direction], state: 2}
		}
	}
	if len(i.stack) > 0 {
		frame := &i.stack[len(i.stack)-1]
		i.currentEntry = frame.node.CompIntMapEntry
		frame.state = 3
	}
}
