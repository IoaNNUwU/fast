### Fast structures in Go

**Linked List**

Linked List implementation where nodes are array chunks

```go
LinkedList
    |
    *-> ListChunk -> ListChunk -> ListChunk -> nil
        |            |            |
        [1, 2, 3]    [4, 5, 6]    [7, 8, 9]
```

Benchmark comparing this with simple LinkedList with 1 element per list node:

- Windows (amd64)
- cpu: AMD Ryzen 5 5500

|Name|ops|throughput|bytes per op|allocations|
|---|---|---|---|---|
|Fast LL Iterator|2280390|528.3 ns/op|72 B/op|5 allocs/op|
|Fast LL for-loop|16261704|67.17 ns/op|0 B/op|0 allocs/op|
|Slow LL for-loop|9322726|126.0 ns/op|0 B/op|0 allocs/op|