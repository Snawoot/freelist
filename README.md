# freelist

[![Go Reference](https://pkg.go.dev/badge/github.com/Snawoot/freelist.svg)](https://pkg.go.dev/github.com/Snawoot/freelist)

Pure Go generic implementation of freelist allocator.

It allows faster allocation and reuse of small objects, making some implementations of datastructures (linked lists, trees, ...) significantly more efficient.

Key features:

* Unlike many other implementations, it doesn't ruin garbage collection for pointers inside allocated objects and plays nicely with garbage collector.
* Unlike many other implementations, it doesn't depend on CGO and doesn't encumber builds.

## Usage

See [godoc examples](https://pkg.go.dev/github.com/Snawoot/freelist#pkg-examples).

## Benchmarks

Here is benchmarks of original `container/list.List` (`BenchmarkContainerList.*`) versus the same List, but augmented by this freelist package (`BenchmarkFreelistList.*`).

```
goos: linux
goarch: amd64
pkg: github.com/Snawoot/freelist
cpu: Intel(R) N100
BenchmarkContainerList/PushFront-4                     	10498634	        138.6 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFront-4                  	 9171722	        173.3 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFrontImmediate-4         	14392966	        97.04 ns/op	      56 B/op	       1 allocs/op
BenchmarkFreelistList/PushFront-4                      	13464580	        89.79 ns/op	      56 B/op	       1 allocs/op
BenchmarkFreelistList/PushPopFront-4                   	12886014	        82.59 ns/op	      56 B/op	       1 allocs/op
BenchmarkFreelistList/PushPopFrontImmediate-4          	30983012	        38.41 ns/op	       8 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushFront-4              	36359986	        32.00 ns/op	       7 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushPopFront-4           	31587457	        38.25 ns/op	       7 B/op	       0 allocs/op
BenchmarkBuiltinNew-4                                  	48650766	        27.10 ns/op	       8 B/op	       1 allocs/op
BenchmarkFreelistAlloc-4                               	68139662	        16.94 ns/op	      18 B/op	       0 allocs/op
BenchmarkWarmedUpFreelistAlloc-4                       	224581018	        5.377 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/Snawoot/freelist	28.662s
```

As you can see it performs about twice faster than original `container/list.List` in worst case and performs 3-5 times faster once freelist reaches right size.
