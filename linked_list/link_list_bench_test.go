package linked_list

import (
	"fmt"
	"strings"
	"testing"
)

const N_ELEM = 128

func Benchmark_fast_LL_iter(b *testing.B) {
	ll := DefaultLinkedList[string]()

	for range N_ELEM {
		ll.PushHead("Hello World")
	}

	for b.Loop() {
		for v := range ll.Iterator() {
			_ = v
			// blackBox(v)
		}
	}
}

func Benchmark_fast_LL_loop(b *testing.B) {
	ll := DefaultLinkedList[string]()

	for range N_ELEM {
		ll.PushHead("Hello World")
	}

	for b.Loop() {
		pointer := ll.chunk

		for pointer != nil {
			for i := 0; i < len(pointer.values); i++ {
				_ = pointer.values[i]
			}
			pointer = pointer.next
		}
	}
}

func Benchmark_fast_LL_first_chunk(b *testing.B) {
	ll := DefaultLinkedList[string]()

	for range N_ELEM {
		ll.PushHead("Hello World")
	}

	for b.Loop() {
		for i := 0; i < len(ll.chunk.values); i++ {
			_ = ll.chunk.values[i]
		}
	}
}

func Benchmark_slow_LL_iter(b *testing.B) {
	ll := NewSlowLL()

	for b.Loop() {
		pointer := ll.next
		for pointer != nil {
			_ = pointer.value
			pointer = pointer.next
		}
	}
	blackBox(ll)
	// PrintLL(0, ll)
}

func PrintLL(rec int, ll *SlowLL[string]) {
	fmt.Println()
	fmt.Printf("%v %v", strings.Repeat(" ", rec), ll)
	pointer := ll.next
	for pointer != nil {
		PrintLL(rec + 1, pointer)
		_ = pointer.value
		pointer = pointer.next
	}
}

type SlowLL[T any] struct {
	next  *SlowLL[T]
	value T
}

func NewSlowLL() *SlowLL[string] {
	last := SlowLL[string]{
		value: "Hello World",
		next:  nil,
	}

	pointer := &last

	for idx := range N_ELEM {
		pointer = &SlowLL[string]{
			value: "Hello World" + fmt.Sprint(idx),
			next:  pointer,
		}
	}
	return pointer
}

var result interface{}

func blackBox(x interface{}) interface{} {
	result = x
	return result
}
