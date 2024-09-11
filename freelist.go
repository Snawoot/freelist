package freelist

import "unsafe"

func defaultNextCap(currentCap int) int {
	if currentCap < 64 {
		return 64
	}
	return currentCap * 2
}

type elt[T any] struct {
	value    T // must be the first field to avoid offset calculations
	nextFree *elt[T]
}

type Freelist[T any] struct {
	NextCapFn func(int) int
	free      *elt[T]
	mem       [][]elt[T]
	cap       int
	len       int
}

func New[T any]() *Freelist[T] {
	return new(Freelist[T])
}

func (fl *Freelist[T]) Free(x *T) {
	found := (*elt[T])(unsafe.Pointer(x))

	var zeroElt elt[T]
	*found = zeroElt

	found.nextFree = fl.free
	fl.free = found
	fl.len--
}

func (fl *Freelist[T]) grow() {
	nextCapFn := defaultNextCap
	if fl.NextCapFn != nil {
		nextCapFn = fl.NextCapFn
	}

	nextCap := nextCapFn(fl.cap)
	if nextCap <= fl.cap {
		panic("NextCapFn returned capacity not larger than current one")
	}
	newChunk := make([]elt[T], nextCap-fl.cap)
	fl.mem = append(fl.mem, newChunk)
	fl.cap = nextCap
	fl.len += len(newChunk)
	for i := range newChunk {
		fl.Free(&newChunk[i].value)
	}

}

func (fl *Freelist[T]) Alloc() *T {
	if fl.free == nil {
		fl.grow()
	}
	found := (*elt[T])(unsafe.Pointer(fl.free))

	fl.free = found.nextFree
	fl.len++

	var zeroElt elt[T]
	*found = zeroElt
	return &found.value
}

func (fl *Freelist[T]) Len() int {
	return fl.len
}

func (fl *Freelist[T]) Cap() int {
	return fl.cap
}

func (fl *Freelist[T]) Clear() {
	fl.len = 0
	fl.cap = 0
	fl.mem = nil
	fl.free = nil
}
