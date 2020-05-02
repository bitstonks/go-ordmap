package ordmap

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type Model struct {
	t       *testing.T
	tree    *Node
	entries []Entry
	r       *rand.Rand
}

func NewModel(t *testing.T) *Model {
	m := &Model{
		t:       t,
		entries: []Entry{},
		r:       rand.New(rand.NewSource(0)),
	}
	m.checkInvariants()
	return m
}

func (m *Model) checkInvariants() {
	m.checkNodesValidity()
	m.checkBalance()
	m.checkElements()
}

func (m *Model) checkNodesValidity() {
	var step func(*Node)
	step = func(n *Node) {
		if n == nil {
			return
		}
		require.GreaterOrEqual(m.t, n.order, uint8(1))
		require.LessOrEqual(m.t, n.order, uint8(MAX))
		for i := int(n.order); i < len(n.entries); i++ {
			require.Equal(m.t, zeroEntry, n.entries[i], fmt.Sprintf("%s: %d %v", n.visual(), n.order, n.entries))
			require.Nil(m.t, n.subtrees[i+1])
		}
		children := 0
		for i := 0; i <= int(n.order); i++ {
			if n.subtrees[i] != nil {
				children += 1
			}
		}
		if n.leaf {
			require.Equal(m.t, children, 0)
			return
		}
		require.Greater(m.t, children, 0)
		for i := 0; i <= int(n.order); i++ {
			step(n.subtrees[i])
		}
	}
	step(m.tree)
}

func (m *Model) checkBalance() {
	var depth func(*Node) int
	depth = func(n *Node) int {
		if n == nil {
			return 0
		}
		if n.leaf {
			return 1
		}
		d1 := depth(n.subtrees[0])
		for _, n = range n.subtrees[1 : n.order+1] {
			require.Equal(m.t, d1, depth(n))
		}
		return d1
	}
	depth(m.tree)
}

func (m *Model) checkElements() {
	require.Equal(m.t, m.entries, m.tree.Entries())
}

func (m *Model) Insert(key Key, value Value) {
	oldTree := m.tree
	oldEntries := oldTree.Entries()

	m.tree = m.tree.Insert(key, value)
	m.insertEntry(key, value)
	m.checkInvariants()

	require.Equal(m.t, oldEntries, oldTree.Entries(), "old tree changed") // persistence check
}

func (m *Model) insertEntry(key int, value Value) {
	for i, e := range m.entries {
		if e.K == key {
			m.entries[i].V = value
			return
		}
	}
	m.entries = append(m.entries, Entry{key, value})
	sort.Slice(m.entries, func(i, j int) bool {
		return m.entries[i].K < m.entries[j].K
	})
}

func (m *Model) Delete(key int) {
	oldTree := m.tree
	oldEntries := oldTree.Entries()

	m.tree = m.tree.Remove(key)
	m.deleteEntry(key)
	m.checkInvariants()

	require.Equal(m.t, oldEntries, oldTree.Entries()) // persistence check
}

func (m *Model) deleteEntry(key int) {
	for i, e := range m.entries {
		if e.K == key {
			copy(m.entries[i:], m.entries[i+1:])
			m.entries = m.entries[:len(m.entries)-1]
		}
	}
}

func TestModel(t *testing.T) {
	sizes := []int{10, 20, 30, 100} // , 400}
	for _, N := range sizes {
		t.Run(fmt.Sprintf("insert_%03d", N), func(t *testing.T) {
			m := NewModel(t)
			for i := 0; i < N; i++ {
				e := m.r.Intn(N)
				m.Insert(e, struct{}{})
			}
		})
	}
	sizes = []int{1, 3, 4, 5, 7, 8, 9, 11, 12, 13, 20, 30, 100} //, 400}
	for _, N := range sizes {
		t.Run(fmt.Sprintf("delete_%03d", N), func(t *testing.T) {
			m := NewModel(t)
			for i := 0; i < N; i++ {
				m.Insert(i, struct{}{})
			}
			for i := 0; i < N; i++ {
				e := m.r.Intn(N)
				m.Delete(e)
			}
		})
	}
}