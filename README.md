# freelist

[![Go Reference](https://pkg.go.dev/badge/github.com/Snawoot/freelist.svg)](https://pkg.go.dev/github.com/Snawoot/freelist)

Pure Go generic implementation of freelist allocator.

It allows faster allocation and reuse of small objects, making some implementations of datastructures (linked lists, trees, ...) significantly more efficient.

Key features:

* Unlike many other implementations, it doesn't ruin garbage collection for pointers inside allocated objects and plays nicely with garbage collector.
* Unlike many other implementations, it doesn't not depend on CGO and doesn't encumber builds.

## Usage

See [godoc examples](https://pkg.go.dev/github.com/Snawoot/freelist#pkg-examples).

## Benchmarks

Here is benchmarks of original `container/list.List` (`BenchmarkContainerList.*`) versus List augmented by this freelist package (`BenchmarkFreelistList.*`).

```
goos: linux
goarch: amd64
pkg: github.com/Snawoot/freelist
cpu: Intel(R) N100
BenchmarkContainerList/PushFront-4                     	 9096877	        141.2 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFront-4                  	 9105859	        160.8 ns/op	      55 B/op	       1 allocs/op
BenchmarkContainerList/PushPopFrontImmediate-4         	12810354	        87.22 ns/op	      56 B/op	       1 allocs/op
BenchmarkFreelistList/PushFront-4                      	13434729	        85.49 ns/op	      67 B/op	       0 allocs/op
BenchmarkFreelistList/PushPopFront-4                   	14114340	        84.38 ns/op	      65 B/op	       0 allocs/op
BenchmarkFreelistList/PushPopFrontImmediate-4          	32279160	        32.19 ns/op	       8 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushFront-4              	41868144	        27.81 ns/op	       7 B/op	       0 allocs/op
BenchmarkFreelistList/WarmedUpPushPopFront-4           	34706287	        34.47 ns/op	       7 B/op	       0 allocs/op
PASS
ok  	github.com/Snawoot/freelist	15.069s
```

As you can see it performs on par with original `container/list.List` in worst case and performs 3-5 times faster once freelist reaches right size.
