package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/edofic/go-ordmap/v2"
)

const (
	mapsize     int = 10_000_000
	readers     int = 7
	concurrency int = readers + 1
)

func run(duration time.Duration) {
	var done uint32 = 0

	m := ordmap.NewBuiltin[int, string]()

	writer := func() {
		// Insert all the keys.
		for i := 0; i < mapsize; i += 1 {
			// FIXME: Should be something like atomic.Store(&m, m.Insert(...))
			m = m.Insert(i, strconv.FormatInt(int64(i), 10))
			if 0 != atomic.LoadUint32(&done) {
				return
			}
		}
		fmt.Println("---------------------")

		// We need a separate random number source for each concurrent
		// thread, because access to the default source is serialized.
		r := rand.New(rand.NewSource(time.Now().UnixMicro()))

		// Randomly update the ordmap.
		n := mapsize
		for 0 == atomic.LoadUint32(&done) {
			i := r.Intn(mapsize)
			// FIXME: Should be something like atomic.Store(&m, m.Insert(...))
			m = m.Insert(i, strconv.FormatInt(int64(n), 10))
			n += 1
		}
	}

	reader := func() {
		r := rand.New(rand.NewSource(time.Now().UnixMicro()))
		for 0 == atomic.LoadUint32(&done) {
			// FIXME: Should be something like atomic.Load(&m).Get(...)
			i := r.Intn(mapsize)
			s, ok := m.Get(i)
			if ok {
				j, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					log.Println("not integer:", s, err)
					panic("bad map")
				}
				if int64(i) > j {
					log.Println("key", i, "> value", j)
					panic("bad map")
				}
			}
		}
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

	// Start the writer.
	go func() {
		defer join.Done()
		init.Done()
		start.Wait()
		writer()
	}()

	// Start all readers.
	for i := 0; i < readers; i += 1 {
		go func() {
			defer join.Done()
			init.Done()
			start.Wait()
			reader()
		}()
	}

	init.Wait()
	start.Done()
	time.Sleep(duration)
	atomic.StoreUint32(&done, 1)
	join.Wait()
}

func main() {
	var freeosmem bool
	var duration time.Duration
	flag.BoolVar(&freeosmem, "freeosmem", false, "Force freeing of unused OS memory")
	flag.DurationVar(&duration, "duration", 2*time.Minute, "Test run duration")
	flag.Parse()

	if freeosmem {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		go func() {
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			for {
				_, ok := <-t.C
				if !ok {
					break
				}
				frees := stats.Frees
				runtime.ReadMemStats(&stats)
				if stats.Frees > frees {
					debug.FreeOSMemory()
				}
			}
		}()
	}

	// Print internal memory usage statistics once every second.
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	go func() {
		t := time.Now()
		var stats runtime.MemStats
		for {
			_, ok := <-timer.C
			if !ok {
				break
			}
			d := uint64(time.Now().Sub(t))
			runtime.ReadMemStats(&stats)
			fmt.Println((d+500_000_000)/1_000_000_000, stats.Sys)
		}
	}()

	run(duration)
}
