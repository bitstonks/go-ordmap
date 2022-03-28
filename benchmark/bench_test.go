package ordmap

import (
	"testing"
)
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
	id int
	prices [1000]byte
}

func BenchmarkNodeBuiltin_Insert(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		m := newIntInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, i)
		}
	})
	b.Run("big value", func(b *testing.B){
		m := newIntOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, &order{id: i})
		}
	})
}

func BenchmarkNode_Insert(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		m := newCompInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, i)
		}
	})
	b.Run("big value", func(b *testing.B){
		m := newCompOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, &order{id: i})
		}
	})
}

func BenchmarkNodeBuiltin_Iterate(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		small := newIntInt()
		for i := 0; i < b.N; i++ {
			small = small.Insert(i, i)
		}
		b.StartTimer()
		for iter := small.Iterate(); !iter.Done(); iter.Next() {
			iter.GetValue()
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		big := newIntOrder()
		for i := 0; i < b.N; i++ {
			big = big.Insert(i, &order{id: i})
		}
		b.StartTimer()
		for iter := big.Iterate(); !iter.Done(); iter.Next() {
			iter.GetValue()
		}
	})
}

func BenchmarkNode_Iterate(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		small := newCompInt()
		for i := 0; i < b.N; i++ {
			small = small.Insert(comp{i, i%2==0}, i)
		}
		b.StartTimer()
		for iter := small.Iterate(); !iter.Done(); iter.Next() {
			iter.GetValue()
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		big := newCompOrder()
		for i := 0; i < b.N; i++ {
			big = big.Insert(comp{i, i%2==0}, &order{id: i})
		}
		b.StartTimer()
		for iter := big.Iterate(); !iter.Done(); iter.Next() {
			iter.GetValue()
		}
	})
}

func BenchmarkNodeBuiltin_Get(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		m := newIntInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, i)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m.Get(i)
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		m := newIntOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, &order{id: i})
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m.Get(i)
		}
	})
}

func BenchmarkNode_Get(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		m := newCompInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, i)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m.Get(comp{i, i%2==0})
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		m := newCompOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, &order{id: i})
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m.Get(comp{i, i%2==0})
		}
	})
}

func BenchmarkNodeBuiltin_Remove(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		m := newIntInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, i)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m = m.Remove(i)
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		m := newIntOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(i, &order{id: i})
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m = m.Remove(i)
		}
	})
}

func BenchmarkNode_Remove(b *testing.B) {
	b.Run("small value", func(b *testing.B){
		b.StopTimer()
		m := newCompInt()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, i)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m = m.Remove(comp{i, i%2==0})
		}
	})
	b.Run("big value", func(b *testing.B){
		b.StopTimer()
		m := newCompOrder()
		for i := 0; i < b.N; i++ {
			m = m.Insert(comp{i, i%2==0}, &order{id: i})
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			m = m.Remove(comp{i, i%2==0})
		}
	})
}