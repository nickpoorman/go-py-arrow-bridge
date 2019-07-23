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

As you can see below, the amount of time to move data across the Python/Go language boundary stays constant as the number of elements increases.

However, as the number of chunks increase, the amount of time also increases. I believe this is due to the large number of CGO calls happening in loops. One solution might be to implement the schema data being gathered in C to reduce the number of CGO calls. A simple solution for this could be to compress your table down to a single chunk before crossing the language boundry.

These results are from my Mid 2012 MacBook Air (1.8GHz i5 / 8 GB 1600 MHz DDR3).

```
(bullseye) ➜  go-py-arrow-bridge git:(master) ✗ PKG_CONFIG_PATH=/Users/nick/anaconda3/envs/bullseye/lib/pkgconfig LD_LIBRARY_PATH=/Users/nick/anaconda3/envs/bullseye/lib/python3.7:/Users/nick/anaconda3/envs/bullseye/lib PYTHONPATH=/Users/nick/anaconda3/envs/bullseye/lib/python3.7/site-packages:/Users/nick/projects/go-py-arrow-bridge/__python__ go test -bench=. -run=- -cpuprofile cpu.prof
goos: darwin
goarch: amd64
pkg: github.com/nickpoorman/go-py-arrow-bridge
BenchmarkAll/BenchmarkZeroCopyChunks_5-4                    3000            402994 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_7-4                    3000            463610 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_9-4                    2000            553570 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_1000-4                   30          67841714 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_1500-4                   20          72966442 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_2000-4                   20          91238313 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_2500-4                   10         120581479 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_3000-4                   10         149069387 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_3500-4                   10         168897623 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_4000-4                   10         187915637 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_4500-4                    5         209468675 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_5000-4                    5         232361156 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_5500-4                    5         249023617 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_6000-4                    5         274385207 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_6500-4                    5         305522949 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_7000-4                    5         324781757 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_7500-4                    3         349889266 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_8000-4                    3         372640132 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_8500-4                    3         394905472 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_9000-4                    3         413965959 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_9500-4                    3         440292768 ns/op
BenchmarkAll/BenchmarkZeroCopyChunks_10000-4                   3         461282623 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_5-4                 10000            166585 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_7-4                 10000            165732 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_9-4                 10000            177193 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_1000-4              10000            166462 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_1500-4              10000            166774 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_2000-4              10000            169948 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_2500-4              10000            171018 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_3000-4              10000            168100 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_3500-4              10000            171136 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_4000-4              10000            166941 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_4500-4              10000            171599 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_5000-4              10000            169485 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_5500-4              10000            169657 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_6000-4              10000            168274 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_6500-4              10000            171372 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_7000-4              10000            168484 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_7500-4              10000            169056 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_8000-4              10000            166486 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_8500-4              10000            167760 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_9000-4              10000            173118 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_9500-4              10000            166797 ns/op
BenchmarkAll/BenchmarkZeroCopyElements_10000-4             10000            169281 ns/op
PASS
ok      github.com/nickpoorman/go-py-arrow-bridge       86.975s
```

## License

(c) 2019 Nick Poorman. Licensed under the Apache License, Version 2.0.
