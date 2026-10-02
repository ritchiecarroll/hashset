// SPDX-License-Identifier: MIT
// Copyright (c) 2021-2026 The go2cs Authors

// Package hashset provides HashSet, a generic set of comparable values built
// on a Go map, with a method surface that mirrors .NET's HashSet<T>.
package hashset

// HashSet represents a distinct collection of T values, i.e., a set.
// A HashSet is not sorted and will not contain duplicate elements.
//
// The zero value is a nil HashSet. It reads as an empty set: len, range and
// every method that only reads the set work on it. Adding an element to it
// panics, so sets are created with NewHashSet.
//
// A floating-point NaN never equals itself, so each Add of a NaN adds a new
// element, Contains and Remove never find it, and only emptying the set
// removes it: Clear, or IntersectWith or IntersectWithSet with an empty
// slice or set.
//
// A HashSet is not safe for concurrent use. To use one from more than one
// goroutine, hold a lock around every use of it, len and range included.
type HashSet[T comparable] map[T]struct{}

// NewHashSet creates a new set containing all elements in the specified slice.
// Returns a new HashSet that contains the specified slice items.
func NewHashSet[T comparable](items []T) HashSet[T] {
	hs := make(HashSet[T], len(items))
	hs.UnionWith(items)
	return hs
}

// Add adds the specified element to a set.
// Returns true if item was added to the set; otherwise, false.
func (hs HashSet[T]) Add(item T) bool {
	if hs.Contains(item) {
		return false
	}

	hs[item] = struct{}{}
	return true
}

// Remove removes the specified element from a set.
// Returns true if item was removed from the set; otherwise, false.
func (hs HashSet[T]) Remove(item T) bool {
	if hs.Contains(item) {
		delete(hs, item)
		return true
	}

	return false
}

// RemoveWhere removes all elements that match the conditions defined by the
// specified predicate from a set.
// Returns the number of elements that were removed from the set.
func (hs HashSet[T]) RemoveWhere(predicate func(T) bool) int {
	var removedCount int

	for k := range hs {
		if predicate(k) && hs.Remove(k) {
			removedCount++
		}
	}

	return removedCount
}

// IsEmpty determines if a set contains no elements.
// Returns true if the set is empty, i.e., has a length of zero; otherwise,
// false.
func (hs HashSet[T]) IsEmpty() bool {
	return len(hs) == 0
}

// Clear removes all elements from a set.
func (hs HashSet[T]) Clear() {
	clear(hs)
}

// Contains determines whether a set contains the specified element.
// Returns true if the set contains the specified element; otherwise, false.
func (hs HashSet[T]) Contains(item T) bool {
	_, ok := hs[item]
	return ok
}

// Keys copies the elements of a set to a slice.
// Returns a slice containing the elements of the set.
func (hs HashSet[T]) Keys() []T {
	keys := make([]T, len(hs))
	i := 0

	for k := range hs {
		keys[i] = k
		i++
	}

	return keys
}

// ExceptWith removes all elements in the specified slice from the current set.
func (hs HashSet[T]) ExceptWith(other []T) {
	for _, v := range other {
		delete(hs, v)
	}
}

// ExceptWithSet removes all elements in the specified set from the current
// set.
func (hs HashSet[T]) ExceptWithSet(other HashSet[T]) {
	hs.ExceptWith(other.Keys())
}

// SymmetricExceptWith modifies the current set to contain only elements that
// are present either in the set or in the specified slice, but not both.
func (hs HashSet[T]) SymmetricExceptWith(other []T) {
	hs.SymmetricExceptWithSet(NewHashSet(other))
}

// SymmetricExceptWithSet modifies the current set to contain only elements
// that are present either in the set or in the specified set, but not both.
func (hs HashSet[T]) SymmetricExceptWithSet(other HashSet[T]) {
	if hs.IsEmpty() {
		// If set is empty, then symmetric difference is other
		hs.UnionWithSet(other)
		return
	}

	for _, v := range other.Keys() {
		if !hs.Remove(v) {
			hs.Add(v)
		}
	}
}

// IntersectWith modifies the current set to contain only elements that are
// present in the set and in the specified slice.
func (hs HashSet[T]) IntersectWith(other []T) {
	hs.IntersectWithSet(NewHashSet(other))
}

// IntersectWithSet modifies the current set to contain only elements that are
// present in the set and in the specified set.
func (hs HashSet[T]) IntersectWithSet(other HashSet[T]) {
	// Intersection of anything with empty set is empty set, so return if
	// count is 0
	if len(hs) == 0 {
		return
	}

	if len(other) == 0 {
		hs.Clear()
		return
	}

	for k := range hs {
		if !other.Contains(k) {
			hs.Remove(k)
		}
	}
}

// UnionWith modifies the current set to contain all elements that are present
// in the set, the specified slice, or both.
func (hs HashSet[T]) UnionWith(other []T) {
	for _, v := range other {
		hs.Add(v)
	}
}

// UnionWithSet modifies the current set to contain all elements that are
// present in the set, the specified set, or both.
func (hs HashSet[T]) UnionWithSet(other HashSet[T]) {
	hs.UnionWith(other.Keys())
}

// SetEquals determines whether a set and the specified slice contain the same
// elements.
// Returns true if the set is equal to other slice items; otherwise, false.
func (hs HashSet[T]) SetEquals(other []T) bool {
	return hs.SetEqualsSet(NewHashSet(other))
}

// SetEqualsSet determines whether a set and the specified set contain the same
// elements.
// Returns true if the set is equal to other set; otherwise, false.
func (hs HashSet[T]) SetEqualsSet(other HashSet[T]) bool {
	if len(hs) != len(other) {
		return false
	}

	for k := range other {
		if !hs.Contains(k) {
			return false
		}
	}

	return true
}

// Overlaps determines whether the current set and a specified slice share
// common elements.
// Returns true if the current set and other slice items share at least one
// common element; otherwise, false.
func (hs HashSet[T]) Overlaps(other []T) bool {
	if hs.IsEmpty() {
		return false
	}

	for _, v := range other {
		if hs.Contains(v) {
			return true
		}
	}

	return false
}

// OverlapsSet determines whether the current set and a specified set share
// common elements.
// Returns true if the current set and other set share at least one common
// element; otherwise, false.
func (hs HashSet[T]) OverlapsSet(other HashSet[T]) bool {
	return hs.Overlaps(other.Keys())
}

// IsSubsetOf determines whether a set is a subset of the specified slice.
// Returns true if the set is a subset of other slice items; otherwise, false.
func (hs HashSet[T]) IsSubsetOf(other []T) bool {
	return hs.IsSubsetOfSet(NewHashSet(other))
}

// IsSubsetOfSet determines whether a set is a subset of the specified set.
// Returns true if the set is a subset of other set; otherwise, false.
func (hs HashSet[T]) IsSubsetOfSet(other HashSet[T]) bool {
	// The empty set is a subset of any set
	if hs.IsEmpty() {
		return true
	}

	// If set has more elements than the other set, then it can't be a subset
	if len(hs) > len(other) {
		return false
	}

	for k := range hs {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

// IsProperSubsetOf determines whether a set is a proper subset of the
// specified slice.
// Returns true if the set is a proper subset of other slice items; otherwise,
// false.
func (hs HashSet[T]) IsProperSubsetOf(other []T) bool {
	return hs.IsProperSubsetOfSet(NewHashSet(other))
}

// IsProperSubsetOfSet determines whether a set is a proper subset of the
// specified set.
// Returns true if the set is a proper subset of other set; otherwise, false.
func (hs HashSet[T]) IsProperSubsetOfSet(other HashSet[T]) bool {
	// The empty set is a proper subset of anything but the empty set
	if hs.IsEmpty() {
		return len(other) > 0
	}

	// If set has more or equal elements than the other set, then it can't be
	// a proper subset
	if len(hs) >= len(other) {
		return false
	}

	for k := range hs {
		if !other.Contains(k) {
			return false
		}
	}

	return true
}

// IsSupersetOf determines whether a set is a superset of the specified slice.
// Returns true if the set is a superset of other slice items; otherwise,
// false.
func (hs HashSet[T]) IsSupersetOf(other []T) bool {
	return hs.IsSupersetOfSet(NewHashSet(other))
}

// IsSupersetOfSet determines whether a set is a superset of the specified set.
// Returns true if the set is a superset of other set; otherwise, false.
func (hs HashSet[T]) IsSupersetOfSet(other HashSet[T]) bool {
	// If other is the empty set then this is a superset
	if len(other) == 0 {
		return true
	}

	// If other set has more elements than set, then it can't be a superset
	if len(other) > len(hs) {
		return false
	}

	for k := range other {
		if !hs.Contains(k) {
			return false
		}
	}

	return true
}

// IsProperSupersetOf determines whether a set is a proper superset of the
// specified slice.
// Returns true if the set is a proper superset of other slice items;
// otherwise, false.
func (hs HashSet[T]) IsProperSupersetOf(other []T) bool {
	return hs.IsProperSupersetOfSet(NewHashSet(other))
}

// IsProperSupersetOfSet determines whether a set is a proper superset of the
// specified set.
// Returns true if the set is a proper superset of other set; otherwise, false.
func (hs HashSet[T]) IsProperSupersetOfSet(other HashSet[T]) bool {
	// The empty set is not a proper superset of any set
	if hs.IsEmpty() {
		return false
	}

	// If other is the empty set then this is a proper superset
	if len(other) == 0 {
		return true // Set has at least one element, based on prior check
	}

	// If other set has more or equal elements than set, then it can't be a
	// proper superset
	if len(other) >= len(hs) {
		return false
	}

	for k := range other {
		if !hs.Contains(k) {
			return false
		}
	}

	return true
}
