[![GoDoc](https://godoc.org/fortio.org/rand?status.svg)](https://pkg.go.dev/fortio.org/rand)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/rand)](https://goreportcard.com/report/fortio.org/rand)
[![CI Checks](https://github.com/fortio/rand/actions/workflows/include.yml/badge.svg)](https://github.com/fortio/rand/actions/workflows/include.yml)
[![codecov](https://codecov.io/github/fortio/rand/graph/badge.svg?token=Yx6QaeQr1b)](https://codecov.io/github/fortio/rand)

# rand

Random number wrapper whose goal is to be used one instance per goroutine.

It's a wrapper over stdlib math/rand/v2 PCG - ie the fastest available in stdlib.

It provides convenience method to either have different rng each time (seed 0) or a specific repeatable sequence.

## Use

```
go get fortio.org/rand@latest
```

## Benchmarks

```
$ go test -bench . .
goos: darwin
goarch: arm64
pkg: fortio.org/rand
cpu: Apple M3 Pro
BenchmarkSharedRand-11                  10285254               119.0 ns/op
BenchmarkPerGoRoutineRand-11            1000000000               0.7205 ns/op
BenchmarkGlobalRand-11                  1000000000               1.115 ns/op
BenchmarkPerGoRoutineChaCha8-11         1000000000               0.9193 ns/op
BenchmarkPCGUint64-11                   1000000000               0.3845 ns/op
BenchmarkChaCha8Uint64-11               1000000000               0.6418 ns/op
PASS
ok      fortio.org/rand 5.849s
```
