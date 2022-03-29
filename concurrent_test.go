package ordmap

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConcurrentAccess(t *testing.T) {
	const iterations int = 1_000_000
	const readers int = 7
	const concurrency int = readers + 1

	var nums [iterations]int
	for i := 0; i < iterations; i += 1 {
		nums[i] = i
	}
	rand.Seed(time.Now().UnixMicro())
	rand.Shuffle(iterations, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	m := NewBuiltin[int, int]()

	writer := func() {
		for _, i := range nums {
			// FIXME: Should be something like atomic.Store(&m, m.Insert(...))
			m = m.Insert(i, i)
		}
	}

	reader := func(keymap func(i int) int) int {
		found := 0
		repeat := 0
		for found == 0 {
			repeat += 1
			for i := 0; i < iterations; i += 1 {
				j := keymap(i)
				// FIXME: Should be something like atomic.Load(&m).Get(...)
				n, exists := m.Get(j)
				if exists {
					require.Equal(t, n, j)
					found += 1
				}
			}
		}

		// Flag any reader that had to repeat its search.
		if repeat > 1 {
			t.Log("repeat:", repeat-1)
		}
		require.Greater(t, found, 0)
		require.LessOrEqual(t, found, iterations)
		return found
	}

	// All concurrent routines are initialized.
	var init sync.WaitGroup
	init.Add(concurrency)

	// All concurrent routines have completed.
	var join sync.WaitGroup
	join.Add(concurrency)

	// Synchronization point for concurrent writer and readers.
	var start sync.WaitGroup
	start.Add(1)

	go func() {
		defer join.Done()
		init.Done()
		start.Wait()
		writer()
	}()

	for i := 0; i < readers; i += 1 {
		go func() {
			// We need a separate random number source for each concurrent
			// reader, because access to the default source is serialized.
			r := rand.New(rand.NewSource(time.Now().UnixMicro()))

			defer join.Done()
			init.Done()
			start.Wait()
			reader(func(i int) int { return r.Intn(iterations) })
		}()
	}

	init.Wait()
	then := time.Now()
	start.Done()
	join.Wait()
	now := time.Now()
	require.Equal(t, m.Len(), iterations)
	require.Equal(t, iterations, reader(func(i int) int { return i }))

	// This is the actual elapsed time in the concurrent writer and readers,
	// without initialization overhead.
	t.Log("elapsed:", int64(now.Sub(then))/1_000_000, "ms")
}
