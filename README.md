# Parallel Fibonacci sequence
This repository contains an implementation of an algorithm to calculate the Fibonacci sequence
in parallel. The [algorithm](./algorithm.md) divides the sequence interval based on the number of available cores
to speed up the computation process. The main program compares the traditional approach of computing
the Fibonacci sequence with the parallelized method implemented in this algorithm.

## Building
No dependency is required to build the main program
```sh
  go build
```

## Running
```sh
  go run .
```
