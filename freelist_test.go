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
	if m.Len() != 1000 {
		t.Fatalf("length assertion failed: got %d, expected 1000", m.Len())
	}
	if m.Cap() != 1024 {
		t.Fatalf("capacity assertion failed: got %d, expected 1024", m.Cap())
	}
	gc()

	for i, ptr := range allocated {
		if *ptr != i {
			t.Fatalf("index %d has value %d", i, *ptr)
		}
		m.Free(ptr)
	}
	if m.Len() != 0 {
		t.Fatalf("length assertion failed: got %d, expected 0", m.Len())
	}
	if m.Cap() != 1024 {
		t.Fatalf("capacity assertion failed: got %d, expected 1024", m.Cap())
	}

	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("length assertion failed: got %d, expected 0", m.Len())
	}
	if m.Cap() != 0 {
		t.Fatalf("capacity assertion failed: got %d, expected 0", m.Cap())
	}
}

func TestDefaultCapacityFunc(t *testing.T) {
	for i := 0; i < 4096; i++ {
		nextCap := defaultNextCap(i)
		if nextCap-i < 64 || nextCap < (i+i/4) {
			t.Fatalf("unexpectedly small allocation (%d) for old cap %d", nextCap, i)
		}
	}
}
