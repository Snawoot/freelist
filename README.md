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
BenchmarkContainerList/PushFront-4                     	11164729	        122.9 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFront-4                  	10755362	        123.4 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFrontImmediate-4         	14260502	        90.09 ns/op	      56 B/op	       1 allocs/op
BenchmarkFreelistList/PushFront-4                      	16091882	        63.31 ns/op	      64 B/op	       0 allocs/op
BenchmarkFreelistList/PushPopFront-4                   	17773947	        69.63 ns/op	      58 B/op	       0 allocs/op
BenchmarkFreelistList/PushPopFrontImmediate-4          	37800050	        28.32 ns/op	       8 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushFront-4              	44983108	        23.58 ns/op	       7 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushPopFront-4           	38999410	        29.72 ns/op	       7 B/op	       0 allocs/op
PASS
ok  	github.com/Snawoot/freelist	14.004s
```

As you can see it performs about twice faster than original `container/list.List` in worst case and performs 3-5 times faster once freelist reaches right size.
