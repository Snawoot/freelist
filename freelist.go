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

import "unsafe"

// defaultNextCap is NextCapFn used by Freelist by default.
func defaultNextCap(currentCap int) int {
	if currentCap < 64 {
		return 64
	}
	return currentCap * 2
}

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
// Methods of Freelist are not safe for concurrent use by multiple goroutines.
type Freelist[T any] struct {
	// If NextCapFn is not nil, it is called to query next capacity value
	// on freelist auto-grow. The currentCap argument of that function
	// is the number of objects freelist can hold at this moment and
	// the returned value is Returned value must be larger than current
	// capacity, otherwise panic will occur.
	//
	// If NextCapFn is nil, default function is used, which doubles capacity
	// each time and initially starts with at least 64 elements.
	//
	// Note that Freelist can be also expanded explicitly by [Freelist.Grow],
	// which means currentCap passed to NextCapFn may be not one of the
	// values returned by NextCapFn previously.
	NextCapFn func(currentCap int) int

	// free is the head of freelist
	free *elt[T]

	// mem is slice which holds extents of memory with actual objects
	mem [][]elt[T]

	// cap is current capacity, the total size of allocated memory extents
	cap int

	// len is current length, the number of allocated objects
	len int
}

// Free deallocates object previously allocated by [Freelist.Alloc].
// Free immediately overwrites freed memory with zero value of corresponding
// type T and marks memory as available for reuse.
//
// Pointer to deallocated object should not be used after call to Free.
func (fl *Freelist[T]) Free(x *T) {
	found := (*elt[T])(unsafe.Pointer(x))

	var zeroElt elt[T]
	*found = zeroElt

	found.nextFree = fl.free
	fl.free = found
	fl.len--
}

// nextCap invokes NextCapFn or default next capacity function if NextCapFn is
// not set.
func (fl *Freelist[T]) nextCap() int {
	if fl.NextCapFn != nil {
		return fl.NextCapFn(fl.cap)
	}
	return defaultNextCap(fl.cap)
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
	fl.mem = append(fl.mem, newChunk)
	fl.cap += n
	fl.len += n
	for i := range newChunk {
		fl.Free(&newChunk[i].value)
	}
}

// autogrow expands memory allocated from runtime to ensure space
// for new allocations from freelist.
func (fl *Freelist[T]) autogrow() {
	growSize := fl.nextCap() - fl.cap
	if growSize <= 0 {
		panic("freelist.Freelist.autogrow: insufficient new capacity")
	}
	fl.Grow(growSize)
}

// Alloc allocates new object. Allocated pointers should be eventually disposed
// with either
//
// - Passing pointer to [Freelist.Free].
//
// - Clearing entire freelist with [Freelist.Clear].
//
// - Dropping reference to entire Freelist and all objects allocated from it.
func (fl *Freelist[T]) Alloc() *T {
	if fl.free == nil {
		fl.autogrow()
	}
	found := fl.free

	fl.free = found.nextFree
	fl.len++

	var zeroElt elt[T]
	*found = zeroElt
	return &found.value
}

// Len returns the number of objects currently allocated from freelist.
func (fl *Freelist[T]) Len() int {
	return fl.len
}

// Cap returns the number of objects that freelist currently can hold.
func (fl *Freelist[T]) Cap() int {
	return fl.cap
}

// Clear resets freelist to initial empty state.
func (fl *Freelist[T]) Clear() {
	fl.len = 0
	fl.cap = 0
	fl.mem = nil
	fl.free = nil
}
