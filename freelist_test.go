package freelist

import (
	"runtime"
	"testing"
)

func gc() {
	for i := 0; i < 3; i++ {
		runtime.GC()
	}
}

func TestSmoke(t *testing.T) {
	var m Freelist[int]
	var allocated []*int
	for i := 0; i < 1000; i++ {
		ptr := m.Alloc()
		*ptr = i
		allocated = append(allocated, ptr)
	}
	gc()

	for i, ptr := range allocated {
		if *ptr != i {
			t.Fatalf("index %d has value %d", i, *ptr)
		}
		m.Free(ptr)
	}

	m.Clear()
}
