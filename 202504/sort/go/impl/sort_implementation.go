package impl

// SortImplementation はソートアルゴリズムの基本実装を提供する
type SortImplementation struct{}

// smallSortThreshold 以下は挿入ソートに切り替える
const smallSortThreshold = 16

func (s *SortImplementation) Sort(data []interface{}) []interface{} {
	n := len(data)
	if n <= 1 {
		return data
	}
	// 元データをコピー
	newArr := make([]interface{}, n)
	copy(newArr, data)

	// 型ごとに in-place ソート
	switch newArr[0].(type) {
	case int:
		quickSortInt(newArr, 0, n-1)
	case string:
		quickSortString(newArr, 0, n-1)
	case float64:
		quickSortFloat(newArr, 0, n-1)
	}
	return newArr
}

// ─── int 用 ───────────────────────────

// 挿入ソート（low～high inclusive）
func insertionIntInterface(a []interface{}, low, high int) {
	for i := low + 1; i <= high; i++ {
		key := a[i].(int)
		j := i - 1
		for j >= low && a[j].(int) > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

// メディアン・オブ・スリー
func median3IntInterface(a []interface{}, low, high int) {
	mid := low + (high-low)/2
	if a[mid].(int) < a[low].(int) {
		a[mid], a[low] = a[low], a[mid]
	}
	if a[high].(int) < a[low].(int) {
		a[high], a[low] = a[low], a[high]
	}
	if a[high].(int) < a[mid].(int) {
		a[high], a[mid] = a[mid], a[high]
	}
	// pivot を末尾に
	a[mid], a[high] = a[high], a[mid]
}

// パーティション
func partitionIntInterface(a []interface{}, low, high int) int {
	pivot := a[high].(int)
	i := low
	for j := low; j < high; j++ {
		if a[j].(int) < pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[high] = a[high], a[i]
	return i
}

// クイックソート本体（末尾再帰最適化＋挿入ソート切り替え）
func quickSortInt(a []interface{}, low, high int) {
	for low < high {
		if high-low <= smallSortThreshold {
			insertionIntInterface(a, low, high)
			return
		}
		median3IntInterface(a, low, high)
		p := partitionIntInterface(a, low, high)
		// 小さいほうを再帰、大きいほうをループで
		if p-low < high-p {
			quickSortInt(a, low, p-1)
			low = p + 1
		} else {
			quickSortInt(a, p+1, high)
			high = p - 1
		}
	}
}

// ─── string 用 ───────────────────────────

func insertionString(a []interface{}, low, high int) {
	for i := low + 1; i <= high; i++ {
		key := a[i].(string)
		j := i - 1
		for j >= low && a[j].(string) > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

func median3String(a []interface{}, low, high int) {
	mid := low + (high-low)/2
	if a[mid].(string) < a[low].(string) {
		a[mid], a[low] = a[low], a[mid]
	}
	if a[high].(string) < a[low].(string) {
		a[high], a[low] = a[low], a[high]
	}
	if a[high].(string) < a[mid].(string) {
		a[high], a[mid] = a[mid], a[high]
	}
	a[mid], a[high] = a[high], a[mid]
}

func partitionString(a []interface{}, low, high int) int {
	pivot := a[high].(string)
	i := low
	for j := low; j < high; j++ {
		if a[j].(string) < pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[high] = a[high], a[i]
	return i
}

func quickSortString(a []interface{}, low, high int) {
	for low < high {
		if high-low <= smallSortThreshold {
			insertionString(a, low, high)
			return
		}
		median3String(a, low, high)
		p := partitionString(a, low, high)
		if p-low < high-p {
			quickSortString(a, low, p-1)
			low = p + 1
		} else {
			quickSortString(a, p+1, high)
			high = p - 1
		}
	}
}

// ─── float64 用 ───────────────────────────

func insertionFloat(a []interface{}, low, high int) {
	for i := low + 1; i <= high; i++ {
		key := a[i].(float64)
		j := i - 1
		for j >= low && a[j].(float64) > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

func median3Float(a []interface{}, low, high int) {
	mid := low + (high-low)/2
	if a[mid].(float64) < a[low].(float64) {
		a[mid], a[low] = a[low], a[mid]
	}
	if a[high].(float64) < a[low].(float64) {
		a[high], a[low] = a[low], a[high]
	}
	if a[high].(float64) < a[mid].(float64) {
		a[high], a[mid] = a[mid], a[high]
	}
	a[mid], a[high] = a[high], a[mid]
}

func partitionFloat(a []interface{}, low, high int) int {
	pivot := a[high].(float64)
	i := low
	for j := low; j < high; j++ {
		if a[j].(float64) < pivot {
			a[i], a[j] = a[j], a[i]
			i++
		}
	}
	a[i], a[high] = a[high], a[i]
	return i
}

func quickSortFloat(a []interface{}, low, high int) {
	for low < high {
		if high-low <= smallSortThreshold {
			insertionFloat(a, low, high)
			return
		}
		median3Float(a, low, high)
		p := partitionFloat(a, low, high)
		if p-low < high-p {
			quickSortFloat(a, low, p-1)
			low = p + 1
		} else {
			quickSortFloat(a, p+1, high)
			high = p - 1
		}
	}
}
