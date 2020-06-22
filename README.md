# cache

![gh actions status](https://github.com/topliceanu/cache/workflows/CI%20Workflow/badge.svg)

This repository contains example implementation of various cache replacement strategies.
It accompanies this blog post [Cache replacement strategies](http://alexandrutopliceanu.ro/)

## Notice

Do not use these implementations for production workloads, they are not thread-safe and only support `int` keys and values.
They are meant for educational and experimentation purposes.
However feel free to get inspiration from the source code when designing your own cache system for your use-case.

## Contents

There are current six algorithms implemented in this repository as well as tests and benchmarks for them.

The interface of the package is intentionally left small to allow for more flexibility.
See [godoc](https://godoc.org/github.com/topliceanu/cache).
All algorithms implement the `Cache` interface.
To keep it simple, the only supported type for keys and values is `int`. This may be extended in the future.
To get an instance of a cache implementation, you need to use the `Factory` function.

## Build

```bash
$ ./script/build
```

## Tests

```bash
$ ./script/test
```

## Lint

```bash
$ ./script/lint
```

## Benchmark

```bash
$ ./script/benchmark
```

Calculate hit-rates with random input stream

```bash
$ go run ./cmd/hit-rate/main.go
```
