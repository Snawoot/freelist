// Package freelist provides generic implementation of freelist allocator
// in pure Go.
//
// It is useful for implementation of algorithms and data structures using a
// lot of small objects. To some extent it is similar to [sync.Pool].
// But unlike sync.Pool, this package provides more predictable retention,
// type safety and control over lifecycle of allocated objects.
// On the other hand, this package requires allocated objects to be explicitly
// freed to avoid memory leaks.
package freelist

import (
	"sync/atomic"
	"unsafe"
)

const growSize = 1024

// elt is an element of allocation slices, used to contain actual value or
// pointer to the next free element.
type elt[T any] struct {
	value    T       // must be the first field to avoid offset calculations
	nextFree *elt[T] // pointer to the next available element
}

// A Freelist is an instance of freelist allocator of objects of type T.
// The zero value for Freelist is an empty freelist ready to use.
//
// A Freelist should not be copied after first use.
//
// Methods of Freelist are safe for concurrent use by multiple goroutines.
type Freelist[T any] struct {
	// free is the head of freelist
	free atomic.Pointer[elt[T]]
}

// Free deallocates object previously allocated by [Freelist.Alloc].
// Free immediately overwrites freed memory with zero value of corresponding
// type T and marks memory as available for reuse.
//
// Pointer to deallocated object should not be used after call to Free.
func (fl *Freelist[T]) Free(x *T) {
	found := (*elt[T])(unsafe.Pointer(x))

	var zeroT T
	found.value = zeroT

	fl.freelistPush(found)
}

// freelistPop borrows element from freelist for allocation.
func (fl *Freelist[T]) freelistPop() *elt[T] {
	for {
		free := fl.free.Load()
		if free == nil {
			fl.autogrow()
			continue
		}
		if fl.free.CompareAndSwap(free, free.nextFree) {
			free.nextFree = nil
			return free
		}
	}
}

// freelistPush marks element as available for reuse.
func (fl *Freelist[T]) freelistPush(e *elt[T]) {
	for {
		free := fl.free.Load()
		e.nextFree = free
		if fl.free.CompareAndSwap(free, e) {
			return
		}
	}
}

// Grow grows the freelist's capacity to guarantee space for another n objects.
// After Grow(n), at least n objects can be allocated from freelist without
// another allocation from runtime.
// If n is negative, Grow will panic.
func (fl *Freelist[T]) Grow(n int) {
	if n < 0 {
		panic("freelist.Freelist.Grow: negative count")
	}
	if n == 0 {
		return
	}
	newChunk := make([]elt[T], n)
	for i := range newChunk {
		fl.freelistPush(&newChunk[i])
	}
}

// autogrow expands memory allocated from runtime to ensure space
// for new allocations from freelist.
func (fl *Freelist[T]) autogrow() {
	fl.Grow(growSize)
}

// Alloc allocates new object. Allocated pointers should be eventually disposed
// with either:
//   - Passing pointer to [Freelist.Free].
//   - Clearing entire freelist with [Freelist.Clear].
//   - Dropping reference to entire Freelist and all objects allocated from it.
func (fl *Freelist[T]) Alloc() *T {
	found := fl.freelistPop()
	return &found.value
}

// Clear resets freelist to initial empty state.
func (fl *Freelist[T]) Clear() {
	fl.free.Store(nil)
}
