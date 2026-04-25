package linked_list

import (
	"iter"
	"slices"
)

// LinkedList optimized for fast iteration and random insertion.
// Consists of ListChunks connected by a singly linked list.
// Each chunk holds a slice of values, which speeds up iteration.
type LinkedList[T any] struct {
	head *ListChunk[T]
}

const defaultChunkCap = 128

func NewLinkedList[T any](chunkCap int) LinkedList[T] {
	return LinkedList[T]{head: newChunk[T](chunkCap)}
}

func DefaultLinkedList[T any]() LinkedList[T] {
	return NewLinkedList[T](defaultChunkCap)
}

// Length returns the total number of elements in the list.
func (l LinkedList[T]) Length() int {
	total := 0
	for chunk := range l.Chunks() {
		total += len(chunk.values)
	}
	return total
}

// PushTail appends values to the end of the list.
func (l LinkedList[T]) PushTail(values ...T) {
	// Find the last chunk
	last := l.head
	for last.next != nil {
		last = last.next
	}

	for _, value := range values {
		if last.isFull() {
			last.next = newChunk[T](cap(last.values))
			last = last.next
		}
		last.values = append(last.values, value)
	}
}

// PushHead prepends values to the beginning of the list (order is preserved).
func (l LinkedList[T]) PushHead(values ...T) {
	for i := len(values) - 1; i >= 0; i-- {
		value := values[i]
		if !l.head.isFull() {
			l.head.values = slices.Insert(l.head.values, 0, value)
		} else {
			newHead := newChunk[T](cap(l.head.values))
			newHead.next = l.head.next
			newHead.values = l.head.values

			l.head.next = newHead
			l.head.values = []T{value}
		}
	}
}

// PushUnordered inserts a value into the first non-full chunk found.
// Scans at most recLimit chunks before falling back to PushHead.
func (l LinkedList[T]) PushUnordered(recLimit int, value T) {
	nextChunk, stop := iter.Pull(l.Chunks())
	defer stop()

	for range recLimit {
		chunk, ok := nextChunk()
		if !ok {
			break
		}
		if !chunk.isFull() {
			chunk.values = append(chunk.values, value)
			return
		}
	}

	l.PushHead(value)
}

// Set sets a value at the given global index across all chunks.
// Returns false if the index is out of bounds.
func (l LinkedList[T]) Set(globalIdx int, value T) bool {
	if globalIdx < 0 {
		return false
	}
	offset := 0
	for chunk := range l.Chunks() {
		chunkLen := len(chunk.values)
		if globalIdx < offset+chunkLen {
			chunk.values[globalIdx-offset] = value
			return true
		}
		offset += chunkLen
	}
	return false
}

// Get returns a value at the given global index across all chunks.
// Returns the value and false if the index is out of bounds.
func (l LinkedList[T]) Get(globalIdx int) (T, bool) {
	if globalIdx < 0 {
		var zero T
		return zero, false
	}

	offset := 0
	for chunk := range l.Chunks() {
		chunkLen := len(chunk.values)
		if globalIdx < offset+chunkLen {
			return chunk.values[globalIdx-offset], true
		}
		offset += chunkLen
	}

	var zero T
	return zero, false
}
// Delete removes a value at the given global index across all chunks.
// Returns false if the index is out of bounds.
func (l *LinkedList[T]) Delete(globalIdx int) bool {
	if globalIdx < 0 {
		return false
	}

	offset := 0
	for chunk := range l.Chunks() {
		chunkLen := len(chunk.values)
		if globalIdx < offset+chunkLen {
			localIdx := globalIdx - offset

			copy(chunk.values[localIdx:], chunk.values[localIdx+1:])

			// Zeroing last element for GC to collect
			var zero T
			chunk.values[len(chunk.values)-1] = zero
			chunk.values = chunk.values[:len(chunk.values)-1]

			if len(chunk.values) == 0 {
				l.removeChunk(chunk)
			}

			return true
		}
		offset += chunkLen
	}

	return false
}

// removeChunk removes an empty chunk from the linked list.
func (l *LinkedList[T]) removeChunk(target *ListChunk[T]) {
	if l.head == target {
		if l.head.next != nil {
			l.head.values = l.head.next.values
			l.head.next = l.head.next.next
		} else {
			l.head.values = []T{}
		}
		return
	}

	for chunk := range l.Chunks() {
		if chunk.next == target {
			chunk.next = target.next
			return
		}
	}
}

// Iterator returns an iterator over all elements with their global indices.
func (l LinkedList[T]) Iterator() func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		globalIdx := 0
		for chunk := range l.Chunks() {
			for localIdx, value := range chunk.values {
				if !yield(globalIdx+localIdx, value) {
					return
				}
			}
			globalIdx += len(chunk.values)
		}
	}
}

// Chunks returns an iterator over all chunks in the list.
func (l LinkedList[T]) Chunks() func(yield func(*ListChunk[T]) bool) {
	return func(yield func(*ListChunk[T]) bool) {
		for current := l.head; current != nil; current = current.next {
			if !yield(current) {
				return
			}
		}
	}
}

// ListChunk is a building block of LinkedList.
// Holds a slice of values and a pointer to the next chunk.
type ListChunk[T any] struct {
	next   *ListChunk[T]
	values []T
}

func (c *ListChunk[T]) isFull() bool {
	return len(c.values) == cap(c.values)
}

// SetLocal sets a value at the given local index within the chunk.
func (c *ListChunk[T]) SetLocal(idx int, value T) {
	c.values[idx] = value
}

// DeleteLocal removes an element at the given local index within the chunk.
func (c *ListChunk[T]) DeleteLocal(idx int) {
	c.values = slices.Delete(c.values, idx, idx+1)
}

func newChunk[T any](chunkCap int) *ListChunk[T] {
	return &ListChunk[T]{
		values: make([]T, 0, chunkCap),
	}
}