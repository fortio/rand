[![GoDoc](https://godoc.org/fortio.org/rand?status.svg)](https://pkg.go.dev/fortio.org/rand)
[![Go Report Card](https://goreportcard.com/badge/fortio.org/rand)](https://goreportcard.com/report/fortio.org/rand)
[![GitHub Release](https://img.shields.io/github/release/fortio/rand.svg?style=flat)](https://github.com/fortio/rand/releases/)
[![CI Checks](https://github.com/fortio/rand/actions/workflows/include.yml/badge.svg)](https://github.com/fortio/rand/actions/workflows/include.yml)
[![codecov](https://codecov.io/github/fortio/rand/graph/badge.svg?token=Yx6QaeQr1b)](https://codecov.io/github/fortio/rand)

# rand

Random number wrapper whose goal is to be used one instance per goroutine

## Install
You can get the binary from [releases](https://github.com/fortio/rand/releases)

Or just run
```
CGO_ENABLED=0 go install fortio.org/rand@latest  # to install (in ~/go/bin typically) or just
CGO_ENABLED=0 go run fortio.org/rand@latest  # to run without install
```

or
```
brew install fortio/tap/rand
```

or
```
docker run -ti fortio/rand
```


## Usage

```
rand help

flags:
```
