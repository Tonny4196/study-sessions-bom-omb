package impl

import (
	"fmt"
	"math/bits"
	"reflect"
)

// Constants from C
const (
	EmptyMarker   = byte(0x80)
	DeletedMarker = byte(0x7F)
	ControlMask   = byte(0x7F)
	GroupSize     = 16 // Define a constant for group size
)

// SwissTable-based hash map implementation using SIMD instructions
type HashMapImplementation struct {
	controlBytes []byte        // Control bytes array including padding
	entries      []bucketEntry // Entries array
	size         int           // Current size
	capacity     int           // Total capacity
	loadFactor   float64       // Load factor (default: 0.75)
	growthFactor float64       // Growth factor (default: 2.0)
}

// Bucket entry structure
type bucketEntry struct {
	key   interface{}
	value interface{}
	hash  uint64
}

// HashMapImpl はHashMapのインターフェース
type HashMapImpl interface {
	Put(key, value interface{})
	Get(key interface{}) (interface{}, bool)
	Remove(key interface{}) bool
	Size() int
	GetAllEntries() map[string]interface{}
}

// NewHashMap は新しいHashMapを作成する
func NewHashMap(initialCapacity int) *HashMapImplementation {
	if initialCapacity <= 0 {
		initialCapacity = 1024
	}

	// Align to 16-byte boundary (optimize for SIMD operations)
	alignedCapacity := (initialCapacity + 15) & ^15

	hashMap := &HashMapImplementation{
		// Important: Add GROUP_SIZE padding to the control array
		controlBytes: make([]byte, alignedCapacity+GroupSize),
		entries:      make([]bucketEntry, alignedCapacity),
		size:         0,
		capacity:     alignedCapacity,
		loadFactor:   0.75,
		growthFactor: 2.0,
	}

	// Initialize all control bytes to EmptyMarker
	for i := range hashMap.controlBytes {
		hashMap.controlBytes[i] = EmptyMarker
	}

	return hashMap
}

// FNV-1a hash constants
const (
	fnvOffset64 = uint64(14695981039346656037)
	fnvPrime64  = uint64(1099511628211)
)

// fnvHash64a implements the FNV-1a hash algorithm for 64-bit hashes
func fnvHash64a(data []byte) uint64 {
	hash := fnvOffset64
	for _, b := range data {
		hash ^= uint64(b)
		hash *= fnvPrime64
	}
	return hash
}

// Calculate a high-quality hash for any key type
func (h *HashMapImplementation) hashKey(key interface{}) uint64 {
	switch k := key.(type) {
	case int:
		return hashInt(uint64(k))
	case int64:
		return hashInt(uint64(k))
	case int32:
		return hashInt(uint64(k))
	case uint:
		return hashInt(uint64(k))
	case uint64:
		return hashInt(k)
	case uint32:
		return hashInt(uint64(k))
	case string:
		// Use FNV-1a for strings
		return fnvHash64a([]byte(k))
	default:
		// Fall back to reflect for other types
		return fnvHash64a([]byte(fmt.Sprintf("%v", key)))
	}
}

// Better hash function for integers
func hashInt(x uint64) uint64 {
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

// Calculate control byte from hash
func (h *HashMapImplementation) calculateControlByte(hash uint64) byte {
	// Use 7 bits from the hash for the control byte, ensure it's not EmptyMarker or DeletedMarker
	h2 := byte(hash & 0x7F)
	if h2 == EmptyMarker&0x7F || h2 == DeletedMarker {
		h2 = 1 // Use a different value
	}
	return h2
}

// Put はキーと値のペアを格納する
func (h *HashMapImplementation) Put(key, value interface{}) {
	// Calculate hash value
	hash := h.hashKey(key)
	h2 := h.calculateControlByte(hash)

	// Check if load factor is exceeded
	if float64(h.size+1) > float64(h.capacity)*h.loadFactor {
		h.resize(int(float64(h.capacity) * h.growthFactor))
	}

	// Calculate start index
	startIdx := int(hash % uint64(h.capacity))

	// Direct lookup (fast path)
	if h.controlBytes[startIdx] == h2 && h.entries[startIdx].hash == hash &&
		equalKeys(h.entries[startIdx].key, key) {
		h.entries[startIdx].value = value
		return
	}

	// Simple linear probing for more reliability
	firstDeletedIdx := -1
	probeIdx := startIdx

	// Search until we find an empty slot or a match
	for i := 0; i < h.capacity; i++ {
		// Check this position
		if h.controlBytes[probeIdx] == EmptyMarker {
			// Found an empty slot - we can insert here or use a previously found deleted slot
			idx := probeIdx

			// Use previously deleted slot if found
			if firstDeletedIdx != -1 {
				idx = firstDeletedIdx
			}

			h.controlBytes[idx] = h2
			h.entries[idx] = bucketEntry{
				key:   key,
				value: value,
				hash:  hash,
			}
			h.size++

			// Copy to padding area if necessary
			if idx < GroupSize {
				h.controlBytes[h.capacity+idx] = h2
			}

			return
		} else if h.controlBytes[probeIdx] == DeletedMarker && firstDeletedIdx == -1 {
			// Remember first deleted slot
			firstDeletedIdx = probeIdx
		} else if h.controlBytes[probeIdx] == h2 {
			// Found a potential match - check for key equality
			if h.entries[probeIdx].hash == hash && equalKeys(h.entries[probeIdx].key, key) {
				// Update existing entry
				h.entries[probeIdx].value = value
				return
			}
		}

		// Simple linear probing - move to next position
		probeIdx = (probeIdx + 1) % h.capacity

		// Safety check to prevent infinite loops
		if probeIdx == startIdx {
			break
		}
	}

	// Use a deleted slot if found and no empty slot was available
	if firstDeletedIdx != -1 {
		h.controlBytes[firstDeletedIdx] = h2
		h.entries[firstDeletedIdx] = bucketEntry{
			key:   key,
			value: value,
			hash:  hash,
		}
		h.size++

		// Copy to padding area if necessary
		if firstDeletedIdx < GroupSize {
			h.controlBytes[h.capacity+firstDeletedIdx] = h2
		}

		return
	}

	// If we get here, the table is completely full
	// Resize and try again
	h.resize(h.capacity * 2)
	h.Put(key, value)
}

// Get はキーに対応する値を取得する
func (h *HashMapImplementation) Get(key interface{}) (interface{}, bool) {
	if h.size == 0 {
		return nil, false
	}

	hash := h.hashKey(key)
	h2 := h.calculateControlByte(hash)
	startIdx := int(hash % uint64(h.capacity))

	// Direct lookup (fast path)
	if h.controlBytes[startIdx] == h2 && h.entries[startIdx].hash == hash &&
		equalKeys(h.entries[startIdx].key, key) {
		return h.entries[startIdx].value, true
	}

	// Simple linear probing for more reliability
	probeIdx := startIdx

	// Search until we find a match or an empty slot
	for i := 0; i < h.capacity; i++ {
		// Check if we found a match
		if h.controlBytes[probeIdx] == h2 &&
			h.entries[probeIdx].hash == hash &&
			equalKeys(h.entries[probeIdx].key, key) {
			return h.entries[probeIdx].value, true
		}

		// If we find an empty slot, the key is not in the table
		if h.controlBytes[probeIdx] == EmptyMarker {
			return nil, false
		}

		// Move to next slot
		probeIdx = (probeIdx + 1) % h.capacity

		// Safety check to prevent infinite loops
		if probeIdx == startIdx {
			break
		}
	}

	return nil, false
}

// Remove はキーに対応するエントリを削除する
func (h *HashMapImplementation) Remove(key interface{}) bool {
	if h.size == 0 {
		return false
	}

	hash := h.hashKey(key)
	h2 := h.calculateControlByte(hash)
	startIdx := int(hash % uint64(h.capacity))

	// Direct lookup (fast path)
	if h.controlBytes[startIdx] == h2 && h.entries[startIdx].hash == hash &&
		equalKeys(h.entries[startIdx].key, key) {
		h.controlBytes[startIdx] = DeletedMarker

		// Update padding area if necessary
		if startIdx < GroupSize {
			h.controlBytes[h.capacity+startIdx] = DeletedMarker
		}

		h.entries[startIdx] = bucketEntry{}
		h.size--
		return true
	}

	// Simple linear probing for reliability
	probeIdx := startIdx

	// Search until we find a match or an empty slot
	for i := 0; i < h.capacity; i++ {
		// Check if we found a match
		if h.controlBytes[probeIdx] == h2 &&
			h.entries[probeIdx].hash == hash &&
			equalKeys(h.entries[probeIdx].key, key) {
			h.controlBytes[probeIdx] = DeletedMarker

			// Update padding area if necessary
			if probeIdx < GroupSize {
				h.controlBytes[h.capacity+probeIdx] = DeletedMarker
			}

			h.entries[probeIdx] = bucketEntry{}
			h.size--
			return true
		}

		// If we find an empty slot, the key is not in the table
		if h.controlBytes[probeIdx] == EmptyMarker {
			return false
		}

		// Move to next slot
		probeIdx = (probeIdx + 1) % h.capacity

		// Safety check to prevent infinite loops
		if probeIdx == startIdx {
			break
		}
	}

	return false
}

// Size は現在の要素数を取得する
func (h *HashMapImplementation) Size() int {
	return h.size
}

// resize はハッシュマップをリサイズする
func (h *HashMapImplementation) resize(newCapacity int) {
	// Ensure we're not downsizing
	if newCapacity <= h.capacity {
		newCapacity = h.capacity * 2
	}

	oldControlBytes := h.controlBytes
	oldEntries := h.entries
	oldCapacity := h.capacity

	// Align to 16-byte boundary (optimize for SIMD operations)
	alignedCapacity := (newCapacity + 15) & ^15

	// Create new arrays - add padding to control bytes
	h.controlBytes = make([]byte, alignedCapacity+GroupSize)
	h.entries = make([]bucketEntry, alignedCapacity)
	h.capacity = alignedCapacity
	h.size = 0

	// Initialize all control bytes to EmptyMarker
	for i := range h.controlBytes {
		h.controlBytes[i] = EmptyMarker
	}

	// Move existing entries to new hash map
	for i := 0; i < oldCapacity; i++ {
		if oldControlBytes[i] != EmptyMarker && oldControlBytes[i] != DeletedMarker {
			entry := oldEntries[i]
			if entry.key != nil {
				// Don't reuse old hash values, recalculate
				h.Put(entry.key, entry.value)
			}
		}
	}
}

// GetAllEntries は全てのエントリを取得する（テスト用）
func (h *HashMapImplementation) GetAllEntries() map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < h.capacity; i++ {
		if h.controlBytes[i] != EmptyMarker && h.controlBytes[i] != DeletedMarker {
			entry := h.entries[i]
			result[fmt.Sprintf("%v", entry.key)] = entry.value
		}
	}
	return result
}

// Check if two keys are equal
func equalKeys(k1, k2 interface{}) bool {
	if k1 == nil || k2 == nil {
		return k1 == k2
	}
	return reflect.DeepEqual(k1, k2)
}

// Calculate log base 2
func log2(x uint64) int {
	return bits.Len64(x) - 1
}
