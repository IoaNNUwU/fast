package linked_list

import (
	"iter"
	"slices"
)

// LinkedList optimized for fast iteration and random insertion.
//
// LinkedList consists of ListChunk's connected by pointers in 1 direction.
//
// ListChunk manages slice of values which allows fast iteration over Linked List.
type LinkedList[T any] struct {
	chunk *ListChunk[T]
}

// Creates new LinkedList with specified slice capacity.
// This LinkedList internally uses chunks of slices connected by pointers.
func NewLinkedList[T any](slice_cap int) LinkedList[T] {
	return LinkedList[T]{
		chunk: chunk[T](slice_cap),
	}
}

// Creates new LinkedList with default slice capacity of 128
func DefaultLinkedList[T any]() LinkedList[T] {
	return NewLinkedList[T](128)
}

// Returns number of elements in LinkedList
func (l LinkedList[T]) Length() int {
	if l.chunk.next != nil {
		return l.chunk.next.lenTail() + len(l.chunk.values)
	} else {
		return len(l.chunk.values)
	}
}

func (l LinkedList[T]) PushTail(values ...T) {
	l.chunk.pushTail(values...)
}

func (l LinkedList[T]) PushHead(values ...T) {
	for _, value := range values {
		if !l.chunk.isFull() {
			l.chunk.values = slices.Insert(l.chunk.values, 0, value)
		} else {
			old_values := l.chunk.values
			old_next := l.chunk.next

			new_chunk := chunk[T](cap(l.chunk.values))
			new_chunk.next = old_next
			new_chunk.values = old_values

			l.chunk.next = new_chunk
			l.chunk.values = []T{value}
		}
	}
}

// Adds a value to any empty space inside this list.
//
// Tries to reuse already allocated space up until rec_limit times,
// then falls back to PushHead().
func (l LinkedList[T]) PushUnordered(rec_limit int, value T) {

	next_chunk, stop := iter.Pull(l.Chunks())
	defer stop()

	for range rec_limit {
		chunk, ok := next_chunk()

		if ok && !chunk.isFull() {
			chunk.values = append(chunk.values, value)
			return
		}
	}
	l.PushHead(value)
}

// Returns an iterator over incices and values
func (l LinkedList[T]) Iterator() func(yield func(int, T) bool) {
	pointer := l.chunk
	idx_accum := 0

	return func(yield func(int, T) bool) {
		for pointer != nil {
			for idx, value := range pointer.values {
				if !yield(idx_accum+idx, value) {
					return
				}
			}
			idx_accum += len(pointer.values)
			pointer = pointer.next
		}
	}
}

// Returns an iterator over incices and values
func (l LinkedList[T]) Chunks() func(yield func(*ListChunk[T]) bool) {
	pointer := l.chunk
	return func(yield func(*ListChunk[T]) bool) {
		for pointer != nil {
			if !yield(pointer) {
				return
			}
			pointer = pointer.next
		}
	}
}

// ListChunk is a building block of LinkedList. Chunks are connected together by pointers.
//
// ListChunk manages slice of values which allows fast iteration over Linked List.
type ListChunk[T any] struct {
	next   *ListChunk[T]
	values []T
}

func (l *ListChunk[T]) lenTail() int {
	if l.next != nil {
		return len(l.values) + l.next.lenTail()
	} else {
		return len(l.values)
	}
}

// Pushes a value to a tail of the list
func (c *ListChunk[T]) pushTail(values ...T) {
	for _, value := range values {
		if c.next != nil {
			c.next.pushTail(value)
		} else {
			if !c.isFull() {
				c.values = append(c.values, value)
			} else {
				c.next = chunk[T](cap(c.values))
				c.next.pushTail(value)
			}
		}
	}
}

func (l *ListChunk[T]) isFull() bool {
	return cap(l.values) == len(l.values)
}

func (l *ListChunk[T]) SetLocal(idx int, value T) {
	l.values[idx] = value
}

func (l *ListChunk[T]) DeleteLocal(idx int) {
	_ = slices.Delete(l.values, idx, idx)
}

// slice_cap argument gets saved internally as capacity of LinkedList.Values
// and gets reused for next list chunks.
func chunk[T any](slice_cap int) *ListChunk[T] {
	return &ListChunk[T]{
		values: make([]T, 0, slice_cap),
	}
}
