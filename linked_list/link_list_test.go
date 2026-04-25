package linked_list

import (
	"testing"
)

func TestSet(t *testing.T) {

	tests := []struct {
		name        string
		chunkCap    int
		initial     []int
		globalIdx   int
		newValue    int
		wantSuccess bool
		wantValues  []int
	}{
		{
			name:        "set value within single chunk",
			chunkCap:    4,
			initial:     []int{1, 2, 3, 4},
			globalIdx:   2,
			newValue:    99,
			wantSuccess: true,
			wantValues:  []int{1, 2, 99, 4},
		},
		{
			name:        "set value in second chunk",
			chunkCap:    4,
			initial:     []int{1, 2, 3, 4, 5, 6, 7, 8},
			globalIdx:   5,
			newValue:    42,
			wantSuccess: true,
			wantValues:  []int{1, 2, 3, 4, 5, 42, 7, 8},
		},
		{
			name:        "set first element",
			chunkCap:    4,
			initial:     []int{1, 2, 3},
			globalIdx:   0,
			newValue:    10,
			wantSuccess: true,
			wantValues:  []int{10, 2, 3},
		},
		{
			name:        "set last element",
			chunkCap:    4,
			initial:     []int{1, 2, 3, 4, 5},
			globalIdx:   4,
			newValue:    77,
			wantSuccess: true,
			wantValues:  []int{1, 2, 3, 4, 77},
		},
		{
			name:        "index out of bounds",
			chunkCap:    4,
			initial:     []int{1, 2, 3},
			globalIdx:   10,
			newValue:    99,
			wantSuccess: false,
			wantValues:  []int{1, 2, 3},
		},
		{
			name:        "negative index",
			chunkCap:    4,
			initial:     []int{1, 2, 3},
			globalIdx:   -1,
			newValue:    99,
			wantSuccess: false,
			wantValues:  []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushTail(tt.initial...)

			got := list.Set(tt.globalIdx, tt.newValue)

			if got != tt.wantSuccess {
				t.Errorf("Set() returned %v, want %v", got, tt.wantSuccess)
			}

			var actualValues []int
			for _, value := range list.Iterator() {
				actualValues = append(actualValues, value)
			}

			if len(actualValues) != len(tt.wantValues) {
				t.Fatalf("list length = %d, want %d", len(actualValues), len(tt.wantValues))
			}

			for i, value := range actualValues {
				if value != tt.wantValues[i] {
					t.Errorf("values[%d] = %d, want %d", i, value, tt.wantValues[i])
				}
			}
		})
	}
}

func TestPushTail(t *testing.T) {
	tests := []struct {
		name       string
		chunkCap   int
		pushValues []int
		wantValues []int
	}{
		{
			name:       "push single value",
			chunkCap:   4,
			pushValues: []int{1},
			wantValues: []int{1},
		},
		{
			name:       "push multiple values within single chunk",
			chunkCap:   4,
			pushValues: []int{1, 2, 3, 4},
			wantValues: []int{1, 2, 3, 4},
		},
		{
			name:       "push values across multiple chunks",
			chunkCap:   4,
			pushValues: []int{1, 2, 3, 4, 5, 6, 7, 8},
			wantValues: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:       "push values one by one across chunk boundary",
			chunkCap:   2,
			pushValues: []int{10, 20, 30},
			wantValues: []int{10, 20, 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushTail(tt.pushValues...)

			assertValues(t, list, tt.wantValues)
		})
	}
}

func TestPushHead(t *testing.T) {
	tests := []struct {
		name       string
		chunkCap   int
		pushValues []int
		wantValues []int
	}{
		{
			name:       "push single value to empty list",
			chunkCap:   4,
			pushValues: []int{1},
			wantValues: []int{1},
		},
		{
			name:       "push multiple values preserves order",
			chunkCap:   4,
			pushValues: []int{1, 2, 3},
			wantValues: []int{1, 2, 3},
		},
		{
			name:       "push values triggers new chunk allocation",
			chunkCap:   2,
			pushValues: []int{1, 2, 3, 4},
			wantValues: []int{1, 2, 3, 4},
		},
		{
			name:       "push values across multiple chunk boundaries",
			chunkCap:   2,
			pushValues: []int{1, 2, 3, 4, 5},
			wantValues: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushHead(tt.pushValues...)

			assertValues(t, list, tt.wantValues)
		})
	}
}

func TestPushHeadAfterTail(t *testing.T) {
	tests := []struct {
		name       string
		chunkCap   int
		tailValues []int
		headValues []int
		wantValues []int
	}{
		{
			name:       "head value prepended before tail values",
			chunkCap:   4,
			tailValues: []int{2, 3, 4},
			headValues: []int{1},
			wantValues: []int{1, 2, 3, 4},
		},
		{
			name:       "multiple head values prepended in order",
			chunkCap:   4,
			tailValues: []int{3, 4, 5},
			headValues: []int{1, 2},
			wantValues: []int{1, 2, 3, 4, 5},
		},
		{
			name:       "head push on full chunk allocates new chunk",
			chunkCap:   2,
			tailValues: []int{3, 4},
			headValues: []int{1, 2},
			wantValues: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushTail(tt.tailValues...)
			list.PushHead(tt.headValues...)

			assertValues(t, list, tt.wantValues)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		chunkCap   int
		pushValues []int
		globalIdx  int
		wantValue  int
		wantOk     bool
	}{
		{
			name:       "get first element",
			chunkCap:   4,
			pushValues: []int{10, 20, 30},
			globalIdx:  0,
			wantValue:  10,
			wantOk:     true,
		},
		{
			name:       "get last element",
			chunkCap:   4,
			pushValues: []int{10, 20, 30},
			globalIdx:  2,
			wantValue:  30,
			wantOk:     true,
		},
		{
			name:       "get middle element",
			chunkCap:   4,
			pushValues: []int{10, 20, 30},
			globalIdx:  1,
			wantValue:  20,
			wantOk:     true,
		},
		{
			name:       "get element across chunk boundary",
			chunkCap:   2,
			pushValues: []int{10, 20, 30, 40},
			globalIdx:  2,
			wantValue:  30,
			wantOk:     true,
		},
		{
			name:       "get element in last chunk",
			chunkCap:   2,
			pushValues: []int{10, 20, 30, 40},
			globalIdx:  3,
			wantValue:  40,
			wantOk:     true,
		},
		{
			name:       "get out of bounds index",
			chunkCap:   4,
			pushValues: []int{10, 20, 30},
			globalIdx:  5,
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "get negative index",
			chunkCap:   4,
			pushValues: []int{10, 20, 30},
			globalIdx:  -1,
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "get from empty list",
			chunkCap:   4,
			pushValues: []int{},
			globalIdx:  0,
			wantValue:  0,
			wantOk:     false,
		},
		{
			name:       "get single element",
			chunkCap:   4,
			pushValues: []int{42},
			globalIdx:  0,
			wantValue:  42,
			wantOk:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushTail(tt.pushValues...)

			gotValue, gotOk := list.Get(tt.globalIdx)

			if gotOk != tt.wantOk {
				t.Fatalf("Get(%d) ok = %v, want %v", tt.globalIdx, gotOk, tt.wantOk)
			}
			if gotValue != tt.wantValue {
				t.Errorf("Get(%d) value = %v, want %v", tt.globalIdx, gotValue, tt.wantValue)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		chunkCap   int
		pushValues []int
		deleteIdx  int
		wantOk     bool
		wantValues []int
	}{
		{
			name:       "delete first element",
			chunkCap:   4,
			pushValues: []int{1, 2, 3, 4},
			deleteIdx:  0,
			wantOk:     true,
			wantValues: []int{2, 3, 4},
		},
		{
			name:       "delete last element",
			chunkCap:   4,
			pushValues: []int{1, 2, 3, 4},
			deleteIdx:  3,
			wantOk:     true,
			wantValues: []int{1, 2, 3},
		},
		{
			name:       "delete middle element",
			chunkCap:   4,
			pushValues: []int{1, 2, 3, 4},
			deleteIdx:  1,
			wantOk:     true,
			wantValues: []int{1, 3, 4},
		},
		{
			name:       "delete element in second chunk",
			chunkCap:   2,
			pushValues: []int{1, 2, 3, 4},
			deleteIdx:  2,
			wantOk:     true,
			wantValues: []int{1, 2, 4},
		},
		{
			name:       "delete last element in chunk — chunk removed",
			chunkCap:   2,
			pushValues: []int{1, 2, 3, 4, 5},
			deleteIdx:  2,
			wantOk:     true,
			wantValues: []int{1, 2, 4, 5},
		},
		{
			name:       "delete only element — list becomes empty",
			chunkCap:   4,
			pushValues: []int{42},
			deleteIdx:  0,
			wantOk:     true,
			wantValues: []int{},
		},
		{
			name:       "delete out of bounds index",
			chunkCap:   4,
			pushValues: []int{1, 2, 3},
			deleteIdx:  5,
			wantOk:     false,
			wantValues: []int{1, 2, 3},
		},
		{
			name:       "delete negative index",
			chunkCap:   4,
			pushValues: []int{1, 2, 3},
			deleteIdx:  -1,
			wantOk:     false,
			wantValues: []int{1, 2, 3},
		},
		{
			name:       "delete from empty list",
			chunkCap:   4,
			pushValues: []int{},
			deleteIdx:  0,
			wantOk:     false,
			wantValues: []int{},
		},
		{
			name:       "delete head chunk — next chunk becomes head",
			chunkCap:   2,
			pushValues: []int{1, 2, 3, 4},
			deleteIdx:  1,
			wantOk:     true,
			wantValues: []int{1, 3, 4},
		},
		{
			name:       "delete all elements one by one",
			chunkCap:   2,
			pushValues: []int{1, 2, 3},
			deleteIdx:  0,
			wantOk:     true,
			wantValues: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int](tt.chunkCap)
			list.PushTail(tt.pushValues...)

			gotOk := list.Delete(tt.deleteIdx)

			if gotOk != tt.wantOk {
				t.Fatalf("Delete(%d) ok = %v, want %v", tt.deleteIdx, gotOk, tt.wantOk)
			}
			assertValues(t, list, tt.wantValues)
		})
	}
}

func TestDeleteSequential(t *testing.T) {
	t.Run("delete all elements one by one until empty", func(t *testing.T) {
		list := NewLinkedList[int](2)
		list.PushTail(1, 2, 3, 4)

		wantAfterEachDelete := [][]int{
			{2, 3, 4},
			{3, 4},
			{4},
			{},
		}

		for step, wantValues := range wantAfterEachDelete {
			ok := list.Delete(0)
			if !ok {
				t.Fatalf("step %d: Delete(0) returned false unexpectedly", step)
			}
			assertValues(t, list, wantValues)
		}
	})

	t.Run("delete returns false on already empty list", func(t *testing.T) {
		list := NewLinkedList[int](2)

		ok := list.Delete(0)
		if ok {
			t.Fatal("Delete(0) on empty list should return false")
		}
	})
}

func assertValues(t *testing.T, list LinkedList[int], wantValues []int) {
	t.Helper()

	var actualValues []int
	for _, value := range list.Iterator() {
		actualValues = append(actualValues, value)
	}

	if len(actualValues) != len(wantValues) {
		t.Fatalf("list length = %d, want %d; values: %v", len(actualValues), len(wantValues), actualValues)
	}

	for i, value := range actualValues {
		if value != wantValues[i] {
			t.Errorf("values[%d] = %d, want %d", i, value, wantValues[i])
		}
	}
}