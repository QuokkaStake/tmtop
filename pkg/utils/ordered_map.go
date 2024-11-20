package utils

import (
	"fmt"
	"strings"
)

// OrderedMap represents a map that maintains insertion order
type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
}

// New creates and returns a new OrderedMap
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   make([]K, 0),
		values: make(map[K]V),
	}
}

// Set adds a key-value pair to the map
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := m.values[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

// Get retrieves a value from the map by key
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.values[key]
	return value, exists
}

// Delete removes a key-value pair from the map
func (m *OrderedMap[K, V]) Delete(key K) {
	if _, exists := m.values[key]; exists {
		delete(m.values, key)
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
				break
			}
		}
	}
}

// Len returns the number of elements in the map
func (m *OrderedMap[K, V]) Len() int {
	return len(m.keys)
}

// Clear removes all elements from the map
func (m *OrderedMap[K, V]) Clear() {
	m.keys = make([]K, 0)
	m.values = make(map[K]V)
}

// Keys returns a slice of all keys in the map, in order
func (m *OrderedMap[K, V]) Keys() []K {
	return append([]K{}, m.keys...)
}

// Values returns a slice of all values in the map, in order
func (m *OrderedMap[K, V]) Values() []V {
	values := make([]V, len(m.keys))
	for i, key := range m.keys {
		values[i] = m.values[key]
	}
	return values
}

// String returns a string representation of the map
func (m *OrderedMap[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, key := range m.keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v: %v", key, m.values[key]))
	}
	sb.WriteString("}")
	return sb.String()
}

// Range iterates over the map in order, calling the given function for each key-value pair
func (m *OrderedMap[K, V]) Range(f func(key K, value V)) {
	for _, key := range m.keys {
		f(key, m.values[key])
	}
}

// GetKeyIndex returns the index of a key in the ordered list of keys
func (m *OrderedMap[K, V]) GetKeyIndex(key K) (int, bool) {
	for i, k := range m.keys {
		if k == key {
			return i, true
		}
	}
	return -1, false
}

// GetByIndex returns the key-value pair at the given index
func (m *OrderedMap[K, V]) GetByIndex(index int) (K, V, bool) {
	if index < 0 || index >= len(m.keys) {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	key := m.keys[index]
	return key, m.values[key], true
}
