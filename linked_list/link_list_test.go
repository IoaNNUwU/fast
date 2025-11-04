package linked_list

import (
	"slices"
	"testing"
)

func TestPushTail(t *testing.T) {
	list := NewLinkedList[string](3)

	list.PushTail("hello", "world", "hello2", "DOESN'T FIT")

	expected_first := []string{"hello", "world", "hello2"}

	if !slices.Equal(list.chunk.values, expected_first) {
		t.Errorf("Expected `%v` in 1st bucket, got `%v`", expected_first, list.chunk.values)
	}
	if list.chunk.next == nil {
		t.Error("Next node should not be nil")
	}
	expected_sec := []string{"DOESN'T FIT"}
	if !slices.Equal(list.chunk.next.values, expected_sec) {
		t.Errorf("Expected `%v` in 2nd bucket, got `%v`", expected_sec, list.chunk.next.values)
	}
}

func TestPushHead(t *testing.T) {
	list := NewLinkedList[string](3)

	list.PushHead("last", "third", "second", "FIRST")

	expected_first := []string{"FIRST"}

	if !slices.Equal(list.chunk.values, expected_first) {
		t.Errorf("Expected `%v` in 1st bucket, got `%v`", expected_first, list.chunk.values)
	}
	if list.chunk.next == nil {
		t.Error("Next node should not be nil")
	}
	expected_sec := []string{"second", "third", "last"}
	if !slices.Equal(list.chunk.next.values, expected_sec) {
		t.Errorf("Expected `%v` in 2nd bucket, got `%v`", expected_sec, list.chunk.next.values)
	}
}

func TestLinkedIteratorIncides(t *testing.T) {
	list := NewLinkedList[string](3)

	list.PushTail("hello", "world", "hello2", "DOESN'T FIT", "DOESN'T FIT 2")

	elements := make([]string, 0, 4)
	indices := make([]int, 0, 4)

	for idx, elem := range list.Iterator() {
		elements = append(elements, elem)
		indices = append(indices, idx)
	}

	expected := []string{"hello", "world", "hello2", "DOESN'T FIT", "DOESN'T FIT 2"}
	if !slices.Equal(elements, expected) {
		t.Errorf("Expected %v, but got %v", expected, elements)
	}
	
	expected_indices := []int { 0, 1, 2, 3, 4 }
	if !slices.Equal(indices, expected_indices) {
		t.Errorf("Expected %v, but got %v", expected_indices, indices)
	}
}
