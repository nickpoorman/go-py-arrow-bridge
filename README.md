# go-py-arrow-bridge

[![GoDoc](https://godoc.org/github.com/nickpoorman/go-py-arrow-bridge?status.svg)](https://godoc.org/github.com/nickpoorman/go-py-arrow-bridge)
[![CircleCI](https://circleci.com/gh/nickpoorman/go-py-arrow-bridge.svg?style=svg)](https://circleci.com/gh/nickpoorman/go-py-arrow-bridge)

A rudimentary bridge for [Apache Arrow](https://github.com/apache/arrow) between Go and Python to facilitate zero-copy.

This Go module demonstrates in the [tests](table_test.go) how easy it is to create an Arrow Table in Python and use the same Arrow Table in Go without copying the underlying buffers.

<!-- ----------------------------------------------------------------------------------------------- -->

## Installation

Add the package to your `go.mod` file:

    require github.com/nickpoorman/go-py-arrow-bridge master

Or, clone the repository:

    git clone --branch master https://github.com/nickpoorman/go-py-arrow-bridge.git $GOPATH/src/github.com/nickpoorman/go-py-arrow-bridge

<!-- ----------------------------------------------------------------------------------------------- -->

## Example

See the [example](cmd/example/example.go) or clone down the repo an run it via `make run`.

## Benchmarks

As you can see below, the amount of time to move data across the Python/Go language boundary stays constant as the amount of data is increased.

These results are from my Mid 2012 MacBook Air (1.8GHz i5 / 8 GB 1600 MHz DDR3).

```
(bullseye) ➜  go-py-arrow-bridge git:(master) ✗ make bench
PKG_CONFIG_PATH=/Users/nick/anaconda3/envs/bullseye/lib/pkgconfig LD_LIBRARY_PATH=/Users/nick/anaconda3/envs/bullseye/lib/python3.7:/Users/nick/anaconda3/envs/bullseye/lib PYTHONPATH=/Users/nick/anaconda3/envs/bullseye/lib/python3.7/site-packages:/Users/nick/projects/go-py-arrow-bridge/__python__ go test  -bench=. -run=- ./...
goos: darwin
goarch: amd64
pkg: github.com/nickpoorman/go-py-arrow-bridge
BenchmarkAll/BenchmarkZeroCopy_0-4                  3000            362634 ns/op
BenchmarkAll/BenchmarkZeroCopy_2-4                  5000            335742 ns/op
BenchmarkAll/BenchmarkZeroCopy_4-4                  5000            349463 ns/op
BenchmarkAll/BenchmarkZeroCopy_6-4                  5000            337202 ns/op
BenchmarkAll/BenchmarkZeroCopy_8-4                  3000            340323 ns/op
BenchmarkAll/BenchmarkZeroCopy_10-4                 5000            323478 ns/op
BenchmarkAll/BenchmarkZeroCopy_1000-4               5000            339729 ns/op
BenchmarkAll/BenchmarkZeroCopy_1500-4               5000            339731 ns/op
BenchmarkAll/BenchmarkZeroCopy_2000-4               5000            336031 ns/op
BenchmarkAll/BenchmarkZeroCopy_2500-4               5000            333809 ns/op
BenchmarkAll/BenchmarkZeroCopy_3000-4               5000            330085 ns/op
BenchmarkAll/BenchmarkZeroCopy_3500-4               5000            368959 ns/op
BenchmarkAll/BenchmarkZeroCopy_4000-4               5000            327952 ns/op
BenchmarkAll/BenchmarkZeroCopy_4500-4               5000            321121 ns/op
BenchmarkAll/BenchmarkZeroCopy_5000-4               3000            343679 ns/op
BenchmarkAll/BenchmarkZeroCopy_5500-4               5000            332056 ns/op
BenchmarkAll/BenchmarkZeroCopy_6000-4               5000            332736 ns/op
BenchmarkAll/BenchmarkZeroCopy_6500-4               5000            327532 ns/op
BenchmarkAll/BenchmarkZeroCopy_7000-4               5000            320282 ns/op
BenchmarkAll/BenchmarkZeroCopy_7500-4               5000            325349 ns/op
BenchmarkAll/BenchmarkZeroCopy_8000-4               5000            324450 ns/op
BenchmarkAll/BenchmarkZeroCopy_8500-4               5000            319664 ns/op
BenchmarkAll/BenchmarkZeroCopy_9000-4               5000            319457 ns/op
BenchmarkAll/BenchmarkZeroCopy_9500-4               5000            322749 ns/op
BenchmarkAll/BenchmarkZeroCopy_10000-4              5000            321663 ns/op
PASS
ok      github.com/nickpoorman/go-py-arrow-bridge       42.747s
```

## License

(c) 2019 Nick Poorman. Licensed under the Apache License, Version 2.0.
