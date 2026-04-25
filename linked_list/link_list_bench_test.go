package linked_list

import (
	"testing"
	"fmt"
)

var(
	globalInt int
	globalBool bool

	lengths []int = []int{1000, 10000, 100000, 1000000}
	chunkCaps []int = []int{6400, 64000, 640000}
)

func BenchmarkGet(b *testing.B) {
	for _, length := range lengths {
		for _, chunkCap := range chunkCaps {
			b.Run(fmt.Sprintf("LinkedList/len=%d/cap=%d", length, chunkCap), func(b *testing.B) {
				list := NewLinkedList[int](chunkCap)
				for i := range length {
					list.PushTail(i)
				}
				midIdx := length / 2

				b.ResetTimer()
				for range b.N {
					v, ok := list.Get(midIdx)
					globalInt = v
					globalBool = ok
				}
			})
		}

		b.Run(fmt.Sprintf("Slice/len=%d", length), func(b *testing.B) {
			slice := make([]int, length)
			for i := range length {
				slice[i] = i
			}
			midIdx := length / 2

			b.ResetTimer()
			for range b.N {
				globalInt = slice[midIdx]
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	for _, length := range lengths {
		for _, chunkCap := range chunkCaps {
			b.Run(fmt.Sprintf("LinkedList/len=%d/cap=%d", length, chunkCap), func(b *testing.B) {
				lists := make([]*LinkedList[int], b.N)
				for i := range b.N {
					
					ll := NewLinkedList[int](chunkCap)
					lists[i] = &ll

					for j := range length {
						lists[i].PushTail(j)
					}
				}
				midIdx := length / 2

				b.ResetTimer()
				for i := range b.N {
					globalBool = lists[i].Delete(midIdx)
				}
			})
		}

		b.Run(fmt.Sprintf("Slice/len=%d", length), func(b *testing.B) {
			slices := make([][]int, b.N)
			for i := range b.N {
				slices[i] = make([]int, length)
				for j := range length {
					slices[i][j] = j
				}
			}
			midIdx := length / 2

			b.ResetTimer()
			for i := range b.N {
				slices[i] = append(slices[i][:midIdx], slices[i][midIdx+1:]...)
			}
		})
	}
}