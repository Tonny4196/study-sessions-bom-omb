package impl

import (
	"sync"
)

// グローバル定数
const (
	INSERTION_SORT_THRESHOLD = 10
	MAX_PARALLEL_DEPTH       = 4
)

// SortImplementation はジェネリックなソート実装を提供する
type SortImplementation[T any] struct{}

func (s *SortImplementation[T]) Sort(array []T, less func(a, b T) bool) []T {
	arrCopy := make([]T, len(array))
	copy(arrCopy, array)

	ParallelQuickSort(arrCopy, less)
	return arrCopy
}

// ParallelQuickSort メソッドで並列クイックソートを使用
func ParallelQuickSort[T any](arr []T, less func(a, b T) bool) {
	var wg sync.WaitGroup
	wg.Add(1)
	go parallelQuickSort(arr, 0, len(arr)-1, less, 0, &wg)
	wg.Wait()
}

// クイックソート（並列版）
func parallelQuickSort[T any](arr []T, low, high int, less func(a, b T) bool, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if high-low+1 <= INSERTION_SORT_THRESHOLD {
		insertionSort(arr, low, high, less)
		return
	}

	p := partition(arr, low, high, less)

	if depth < MAX_PARALLEL_DEPTH {
		wg.Add(2)
		go parallelQuickSort(arr, low, p-1, less, depth+1, wg)
		go parallelQuickSort(arr, p+1, high, less, depth+1, wg)
	} else {
		quickSort(arr, low, p-1, less)
		quickSort(arr, p+1, high, less)
	}
}

// 通常クイックソート（再帰で使う）
func quickSort[T any](arr []T, low, high int, less func(a, b T) bool) {
	for low < high {
		if high-low+1 <= INSERTION_SORT_THRESHOLD {
			insertionSort(arr, low, high, less)
			break
		}
		p := partition(arr, low, high, less)
		if p-low < high-p {
			quickSort(arr, low, p-1, less)
			low = p + 1
		} else {
			quickSort(arr, p+1, high, less)
			high = p - 1
		}
	}
}

// 挿入ソート
func insertionSort[T any](arr []T, low, high int, less func(a, b T) bool) {
	for i := low + 1; i <= high; i++ {
		key := arr[i]
		j := i - 1
		for j >= low && less(key, arr[j]) {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

// 三点中央値パーティション
func partition[T any](arr []T, low, high int, less func(a, b T) bool) int {
	mid := (low + high) / 2

	if less(arr[mid], arr[low]) {
		arr[low], arr[mid] = arr[mid], arr[low]
	}
	if less(arr[high], arr[low]) {
		arr[low], arr[high] = arr[high], arr[low]
	}
	if less(arr[high], arr[mid]) {
		arr[mid], arr[high] = arr[high], arr[mid]
	}

	arr[mid], arr[high-1] = arr[high-1], arr[mid]
	pivot := arr[high-1]

	i, j := low, high-1
	for {
		for i++; less(arr[i], pivot); i++ {
		}
		for j--; less(pivot, arr[j]); j-- {
		}
		if i >= j {
			break
		}
		arr[i], arr[j] = arr[j], arr[i]
	}
	arr[i], arr[high-1] = arr[high-1], arr[i]
	return i
}
