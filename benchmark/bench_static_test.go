package ordmap

import (
	"fmt"
	"testing"
)

const mapSize = 1000000

type comp struct {
	i int
	b bool
}

func (c comp) Less(o comp) bool {
	if c.b != o.b {
		return c.b // i.e. c.b < o.b
	}
	if c.b {
		return c.i < o.i
	}
	return c.i > o.i
}

type order struct {
	id     int
	prices [1000]byte
}

func BenchmarkNodeBuiltin_Insert(b *testing.B) {
	small := newIntInt()
	big := newIntOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(i, i)
		big = big.Insert(i, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Insert(i+mapSize, i)
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			big.Insert(i+mapSize, &order{id: i})
		}
	})
}

func BenchmarkNode_Insert(b *testing.B) {
	small := newCompInt()
	big := newCompOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(comp{i, i%2 == 0}, i)
		big = big.Insert(comp{i, i%2 == 0}, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Insert(comp{i + mapSize, i%2 == 0}, i)
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			big.Insert(comp{i + mapSize, i%2 == 0}, &order{id: i})
		}
	})
}

func BenchmarkNodeBuiltin_Iterate(b *testing.B) {
	small := newIntInt()
	big := newIntOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(i, i)
		big = big.Insert(i, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for iter := small.Iterate(); !iter.Done(); iter.Next() {
				iter.GetValue()
			}
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for iter := big.Iterate(); !iter.Done(); iter.Next() {
				iter.GetValue()
			}
		}
	})
}

func BenchmarkNode_Iterate(b *testing.B) {
	small := newCompInt()
	big := newCompOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(comp{i, i%2 == 0}, i)
		big = big.Insert(comp{i, i%2 == 0}, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for iter := small.Iterate(); !iter.Done(); iter.Next() {
				iter.GetValue()
			}
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for iter := big.Iterate(); !iter.Done(); iter.Next() {
				iter.GetValue()
			}
		}
	})
}

func BenchmarkNodeBuiltin_Get(b *testing.B) {
	small := newIntInt()
	big := newIntOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(i, i)
		big = big.Insert(i, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Get(i % mapSize)
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			big.Get(i % mapSize)
		}
	})
}

func BenchmarkNode_Get(b *testing.B) {
	small := newCompInt()
	big := newCompOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(comp{i, i%2 == 0}, i)
		big = big.Insert(comp{i, i%2 == 0}, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Get(comp{i % mapSize, i%2 == 0})
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			big.Get(comp{i % mapSize, i%2 == 0})
		}
	})
}

func BenchmarkNodeBuiltin_Remove(b *testing.B) {
	small := newIntInt()
	big := newIntOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(i, i)
		big = big.Insert(i, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Remove(i % mapSize)
		}
	})
	b.Run("big value", func(b *testing.B) {
		x := newIntOrder()
		for i := 0; i < b.N; i++ {
			x = big.Remove(i % mapSize)
		}
		if x.Len() == 0 {
			fmt.Println("nil")
		}
	})
}

func BenchmarkNode_Remove(b *testing.B) {
	small := newCompInt()
	big := newCompOrder()
	for i := 0; i < mapSize; i++ {
		small = small.Insert(comp{i, i%2 == 0}, i)
		big = big.Insert(comp{i, i%2 == 0}, &order{id: i})
	}
	b.Run("small value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			small.Remove(comp{i % mapSize, i%2 == 0})
		}
	})
	b.Run("big value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			big.Remove(comp{i % mapSize, i%2 == 0})
		}
	})
}
