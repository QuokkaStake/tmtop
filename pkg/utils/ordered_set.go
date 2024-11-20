package utils

import (
	"fmt"
	"strings"
)

// OrderedSet represents a set that maintains insertion order
type OrderedSet[T comparable] struct {
	elements []T
	set      map[T]struct{}
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		elements: []T{},
		set:      make(map[T]struct{}),
	}
}

// Add adds an element to the set if it doesn't already exist
func (s *OrderedSet[T]) Add(element T) {
	if _, exists := s.set[element]; !exists {
		s.elements = append(s.elements, element)
		s.set[element] = struct{}{}
	}
}

// Remove removes an element from the set
func (s *OrderedSet[T]) Remove(element T) {
	if _, exists := s.set[element]; exists {
		delete(s.set, element)
		for i, v := range s.elements {
			if v == element {
				s.elements = append(s.elements[:i], s.elements[i+1:]...)
				break
			}
		}
	}
}

// Contains checks if an element exists in the set
func (s *OrderedSet[T]) Contains(element T) bool {
	_, exists := s.set[element]
	return exists
}

// Size returns the number of elements in the set
func (s *OrderedSet[T]) Size() int {
	return len(s.elements)
}

// Clear removes all elements from the set
func (s *OrderedSet[T]) Clear() {
	s.elements = []T{}
	s.set = make(map[T]struct{})
}

// Elements returns a slice of all elements in the set, in order
func (s *OrderedSet[T]) Elements() []T {
	return append([]T{}, s.elements...)
}

// String returns a string representation of the set
func (s *OrderedSet[T]) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, element := range s.elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", element))
	}
	sb.WriteString("]")
	return sb.String()
}

// Union returns a new OrderedSet containing all elements from both sets
func (s *OrderedSet[T]) Union(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for _, element := range s.elements {
		result.Add(element)
	}
	for _, element := range other.elements {
		result.Add(element)
	}
	return result
}

// Intersection returns a new OrderedSet containing elements common to both sets
func (s *OrderedSet[T]) Intersection(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for _, element := range s.elements {
		if other.Contains(element) {
			result.Add(element)
		}
	}
	return result
}

// Difference returns a new OrderedSet containing elements in s that are not in other
func (s *OrderedSet[T]) Difference(other *OrderedSet[T]) *OrderedSet[T] {
	result := NewOrderedSet[T]()
	for _, element := range s.elements {
		if !other.Contains(element) {
			result.Add(element)
		}
	}
	return result
}
