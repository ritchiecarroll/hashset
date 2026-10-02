// SPDX-License-Identifier: MIT
// Copyright (c) 2021-2026 The go2cs Authors

package hashset

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
)

// sorted returns the elements of hs in ascending order.
func sorted(hs HashSet[int]) []int {
	keys := hs.Keys()
	sort.Ints(keys)
	return keys
}

// equalInts reports whether a and b hold the same values in the same order.
func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func expectSet(t *testing.T, hs HashSet[int], want ...int) {
	t.Helper()

	sort.Ints(want)

	if got := sorted(hs); !equalInts(got, want) {
		t.Fatalf("set = %v, want %v", got, want)
	}
}

func TestNewHashSet(t *testing.T) {
	expectSet(t, NewHashSet[int](nil))
	expectSet(t, NewHashSet([]int{}))
	expectSet(t, NewHashSet([]int{3, 1, 2}), 1, 2, 3)
	expectSet(t, NewHashSet([]int{1, 1, 2, 2, 2}), 1, 2)

	// The new set does not alias the input slice.
	items := []int{1, 2}
	hs := NewHashSet(items)
	items[0] = 9
	expectSet(t, hs, 1, 2)

	// The set is a map, so len works directly.
	if len(hs) != 2 {
		t.Fatalf("len = %d, want 2", len(hs))
	}
}

func TestNewHashSetStrings(t *testing.T) {
	hs := NewHashSet([]string{"b", "a", "b", ""})

	if len(hs) != 3 || !hs.Contains("") || !hs.Contains("a") || !hs.Contains("b") {
		t.Fatalf("unexpected set %v", hs.Keys())
	}
}

func TestZeroValueReads(t *testing.T) {
	var hs HashSet[int] // nil map

	if !hs.IsEmpty() {
		t.Fatal("nil set should be empty")
	}

	if hs.Contains(1) {
		t.Fatal("nil set should contain nothing")
	}

	if keys := hs.Keys(); len(keys) != 0 {
		t.Fatalf("Keys = %v, want empty", keys)
	}

	if hs.Remove(1) {
		t.Fatal("Remove on nil set should return false")
	}

	if n := hs.RemoveWhere(func(int) bool { return true }); n != 0 {
		t.Fatalf("RemoveWhere = %d, want 0", n)
	}

	hs.Clear()
	hs.ExceptWith([]int{1})

	if !hs.IsSubsetOf([]int{1}) || !hs.IsProperSubsetOf([]int{1}) || hs.Overlaps([]int{1}) {
		t.Fatal("nil set relations are wrong")
	}
}

// TestZeroValueIsAnEmptySet checks that the zero value, a nil HashSet, reads
// as an empty set: len, range and every method that only reads the set work
// on it and give the answers an empty set gives.
func TestZeroValueIsAnEmptySet(t *testing.T) {
	var hs HashSet[int] // nil map

	if hs != nil {
		t.Fatal("the zero value should be a nil HashSet")
	}

	if len(hs) != 0 {
		t.Fatalf("len = %d, want 0", len(hs))
	}

	for v := range hs {
		t.Fatalf("range over a nil set yielded %v", v)
	}

	if !hs.IsEmpty() || hs.Contains(0) || len(hs.Keys()) != 0 {
		t.Error("IsEmpty, Contains or Keys is wrong on a nil set")
	}

	one := []int{1}
	oneSet := NewHashSet(one)

	if !hs.SetEquals(nil) || hs.SetEquals(one) || !hs.SetEqualsSet(nil) || hs.SetEqualsSet(oneSet) {
		t.Error("SetEquals or SetEqualsSet is wrong on a nil set")
	}

	if hs.Overlaps(nil) || hs.Overlaps(one) || hs.OverlapsSet(nil) || hs.OverlapsSet(oneSet) {
		t.Error("Overlaps or OverlapsSet is wrong on a nil set")
	}

	if !hs.IsSubsetOf(nil) || !hs.IsSubsetOf(one) || !hs.IsSubsetOfSet(nil) || !hs.IsSubsetOfSet(oneSet) {
		t.Error("IsSubsetOf or IsSubsetOfSet is wrong on a nil set")
	}

	if hs.IsProperSubsetOf(nil) || !hs.IsProperSubsetOf(one) || hs.IsProperSubsetOfSet(nil) || !hs.IsProperSubsetOfSet(oneSet) {
		t.Error("IsProperSubsetOf or IsProperSubsetOfSet is wrong on a nil set")
	}

	if !hs.IsSupersetOf(nil) || hs.IsSupersetOf(one) || !hs.IsSupersetOfSet(nil) || hs.IsSupersetOfSet(oneSet) {
		t.Error("IsSupersetOf or IsSupersetOfSet is wrong on a nil set")
	}

	if hs.IsProperSupersetOf(nil) || hs.IsProperSupersetOf(one) || hs.IsProperSupersetOfSet(nil) || hs.IsProperSupersetOfSet(oneSet) {
		t.Error("IsProperSupersetOf or IsProperSupersetOfSet is wrong on a nil set")
	}
}

// panics reports whether f panics.
func panics(f func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	f()
	return false
}

// TestZeroValueAddPanics checks that adding an element to the zero value, a
// nil HashSet, panics, through Add and through every method that adds.
func TestZeroValueAddPanics(t *testing.T) {
	var hs HashSet[int] // nil map

	adders := []struct {
		name string
		add  func()
	}{
		{"Add", func() { hs.Add(1) }},
		{"UnionWith", func() { hs.UnionWith([]int{1}) }},
		{"UnionWithSet", func() { hs.UnionWithSet(NewHashSet([]int{1})) }},
		{"SymmetricExceptWith", func() { hs.SymmetricExceptWith([]int{1}) }},
		{"SymmetricExceptWithSet", func() { hs.SymmetricExceptWithSet(NewHashSet([]int{1})) }},
	}

	for _, a := range adders {
		if !panics(a.add) {
			t.Errorf("%s on a nil set should panic", a.name)
		}
	}

	if len(hs) != 0 {
		t.Fatalf("len = %d after the failed adds, want 0", len(hs))
	}

	// Adding nothing does not panic.
	if panics(func() { hs.UnionWith(nil) }) || panics(func() { hs.UnionWithSet(nil) }) {
		t.Error("an empty union on a nil set should not panic")
	}
}

func TestAdd(t *testing.T) {
	hs := NewHashSet[int](nil)

	if !hs.Add(1) {
		t.Fatal("first Add should return true")
	}

	if hs.Add(1) {
		t.Fatal("duplicate Add should return false")
	}

	if !hs.Add(0) {
		t.Fatal("Add of the zero value should return true")
	}

	expectSet(t, hs, 0, 1)
}

func TestRemove(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.Remove(1) {
		t.Fatal("Remove of a member should return true")
	}

	if hs.Remove(1) {
		t.Fatal("second Remove should return false")
	}

	if hs.Remove(5) {
		t.Fatal("Remove of a non-member should return false")
	}

	expectSet(t, hs, 2)
}

func TestRemoveWhere(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3, 4, 5, 6})

	if n := hs.RemoveWhere(func(v int) bool { return v%2 == 0 }); n != 3 {
		t.Fatalf("RemoveWhere = %d, want 3", n)
	}

	expectSet(t, hs, 1, 3, 5)

	if n := hs.RemoveWhere(func(int) bool { return false }); n != 0 {
		t.Fatalf("RemoveWhere(none) = %d, want 0", n)
	}

	if n := hs.RemoveWhere(func(int) bool { return true }); n != 3 {
		t.Fatalf("RemoveWhere(all) = %d, want 3", n)
	}

	expectSet(t, hs)
}

func TestIsEmptyAndClear(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	if hs.IsEmpty() {
		t.Fatal("set with members reported empty")
	}

	hs.Clear()

	if !hs.IsEmpty() {
		t.Fatal("set not empty after Clear")
	}

	// The set is still usable after Clear.
	hs.Add(4)
	expectSet(t, hs, 4)
}

// TestClearRemovesNaN checks that a set holding NaN is empty after Clear. A
// NaN never equals itself, so a key lookup never finds it and delete cannot
// remove it.
func TestClearRemovesNaN(t *testing.T) {
	hs := NewHashSet([]float64{1, math.NaN(), math.NaN()})
	hs.Add(math.NaN())

	if len(hs) != 4 {
		t.Fatalf("len = %d, want 4: 1 and three NaN elements", len(hs))
	}

	hs.Clear()

	if len(hs) != 0 || !hs.IsEmpty() {
		t.Fatalf("len = %d after Clear, want 0", len(hs))
	}

	// The set is still usable after Clear.
	hs.Add(2)

	if len(hs) != 1 || !hs.Contains(2) {
		t.Fatalf("set = %v after Clear and Add(2), want [2]", hs.Keys())
	}
}

// TestNaNElements checks the documented behaviour of NaN elements: each Add of
// a NaN adds a new element, Contains and Remove never find it, and only
// emptying the set removes it.
func TestNaNElements(t *testing.T) {
	nan := math.NaN()
	hs := NewHashSet([]float64{1})

	// Each Add of a NaN adds a new element.
	if !hs.Add(nan) || !hs.Add(nan) {
		t.Fatal("each Add of a NaN should return true")
	}

	if len(hs) != 3 {
		t.Fatalf("len = %d after adding two NaN to {1}, want 3", len(hs))
	}

	// Contains and Remove never find it.
	if hs.Contains(nan) {
		t.Error("Contains should not find a NaN")
	}

	if hs.Remove(nan) {
		t.Error("Remove should not find a NaN")
	}

	// No other removal takes a NaN out, because none of them finds it.
	nanCount := func() int {
		n := 0

		for v := range hs {
			if math.IsNaN(v) {
				n++
			}
		}

		return n
	}

	if n := hs.RemoveWhere(math.IsNaN); n != 0 {
		t.Errorf("RemoveWhere(IsNaN) = %d, want 0", n)
	}

	hs.ExceptWith([]float64{nan})
	hs.IntersectWith([]float64{1})
	hs.IntersectWithSet(NewHashSet([]float64{1}))

	// The intersections keep 1, which is in their argument, beside the NaNs.
	if got := nanCount(); got != 2 || !hs.Contains(1) || len(hs) != 3 {
		t.Fatalf("set = %v, want 1 and two NaN elements", hs.Keys())
	}

	hs.ExceptWithSet(NewHashSet(hs.Keys()))

	if got := nanCount(); got != 2 {
		t.Fatalf("%d NaN elements left, want 2: no removal finds a NaN", got)
	}

	// The removals above still work on an element that equals itself.
	if hs.Contains(1) || len(hs) != 2 {
		t.Fatalf("set = %v, want two NaN elements and nothing else", hs.Keys())
	}

	// SymmetricExceptWith and SymmetricExceptWithSet do not find a NaN either:
	// each adds the NaN elements of its argument and removes none.
	sym := NewHashSet(hs.Keys())
	sym.SymmetricExceptWith([]float64{nan})
	sym.SymmetricExceptWithSet(NewHashSet(sym.Keys()))

	if len(sym) != 6 {
		t.Errorf("len = %d after the symmetric differences, want 6", len(sym))
	}

	// Intersecting with an empty slice or set empties the set, so that is the
	// one other call that removes a NaN.
	for _, intersect := range []func(HashSet[float64]){
		func(c HashSet[float64]) { c.IntersectWith(nil) },
		func(c HashSet[float64]) { c.IntersectWithSet(nil) },
	} {
		c := NewHashSet(hs.Keys())

		if len(c) != 2 {
			t.Fatalf("copy has %d elements, want 2", len(c))
		}

		intersect(c)

		if len(c) != 0 {
			t.Errorf("len = %d after intersecting with nothing, want 0", len(c))
		}
	}

	hs.Clear()

	if len(hs) != 0 {
		t.Fatalf("len = %d after Clear, want 0", len(hs))
	}
}

func TestContains(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.Contains(1) || !hs.Contains(2) || hs.Contains(3) {
		t.Fatal("Contains gave a wrong answer")
	}
}

func TestKeys(t *testing.T) {
	hs := NewHashSet([]int{5, 3, 1})
	keys := hs.Keys()

	if len(keys) != 3 {
		t.Fatalf("len(Keys) = %d, want 3", len(keys))
	}

	// Keys returns a copy: changing it does not change the set.
	keys[0] = 99
	expectSet(t, hs, 1, 3, 5)
}

func TestExceptWith(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3, 4})
	hs.ExceptWith([]int{2, 2, 4, 9})
	expectSet(t, hs, 1, 3)

	hs.ExceptWith(nil)
	expectSet(t, hs, 1, 3)

	hs.ExceptWith(hs.Keys())
	expectSet(t, hs)
}

func TestExceptWithSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})
	hs.ExceptWithSet(NewHashSet([]int{3, 4}))
	expectSet(t, hs, 1, 2)

	hs.ExceptWithSet(nil)
	expectSet(t, hs, 1, 2)

	hs.ExceptWithSet(hs)
	expectSet(t, hs)
}

func TestSymmetricExceptWith(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})
	hs.SymmetricExceptWith([]int{3, 4, 5})
	expectSet(t, hs, 1, 2, 4, 5)

	hs.SymmetricExceptWith(nil)
	expectSet(t, hs, 1, 2, 4, 5)

	// An empty receiver becomes the other collection.
	empty := NewHashSet[int](nil)
	empty.SymmetricExceptWith([]int{7, 7, 8})
	expectSet(t, empty, 7, 8)
}

func TestSymmetricExceptWithSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2})
	hs.SymmetricExceptWithSet(NewHashSet([]int{2, 3}))
	expectSet(t, hs, 1, 3)

	hs.SymmetricExceptWithSet(nil)
	expectSet(t, hs, 1, 3)

	// The symmetric difference of a set with itself is empty.
	hs.SymmetricExceptWithSet(hs)
	expectSet(t, hs)

	empty := NewHashSet[int](nil)
	empty.SymmetricExceptWithSet(NewHashSet([]int{4}))
	expectSet(t, empty, 4)
}

func TestIntersectWith(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3, 4})
	hs.IntersectWith([]int{2, 4, 4, 6})
	expectSet(t, hs, 2, 4)

	hs.IntersectWith(nil)
	expectSet(t, hs)

	// An empty receiver stays empty.
	hs.IntersectWith([]int{1, 2})
	expectSet(t, hs)
}

func TestIntersectWithSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})
	hs.IntersectWithSet(NewHashSet([]int{3, 1, 9}))
	expectSet(t, hs, 1, 3)

	hs.IntersectWithSet(hs)
	expectSet(t, hs, 1, 3)

	hs.IntersectWithSet(NewHashSet([]int{7}))
	expectSet(t, hs)

	full := NewHashSet([]int{1})
	full.IntersectWithSet(nil)
	expectSet(t, full)

	empty := NewHashSet[int](nil)
	empty.IntersectWithSet(NewHashSet([]int{1}))
	expectSet(t, empty)
}

func TestUnionWith(t *testing.T) {
	hs := NewHashSet([]int{1})
	hs.UnionWith([]int{2, 2, 1, 3})
	expectSet(t, hs, 1, 2, 3)

	hs.UnionWith(nil)
	expectSet(t, hs, 1, 2, 3)
}

func TestUnionWithSet(t *testing.T) {
	hs := NewHashSet([]int{1})
	hs.UnionWithSet(NewHashSet([]int{2, 3}))
	expectSet(t, hs, 1, 2, 3)

	hs.UnionWithSet(hs)
	expectSet(t, hs, 1, 2, 3)

	hs.UnionWithSet(nil)
	expectSet(t, hs, 1, 2, 3)
}

func TestSetEquals(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	cases := []struct {
		other []int
		want  bool
	}{
		{[]int{3, 2, 1}, true},
		{[]int{1, 2}, false},
		{[]int{1, 2, 4}, false},
		{[]int{1, 2, 3, 4}, false},
		{nil, false},
	}

	for _, c := range cases {
		if got := hs.SetEquals(c.other); got != c.want {
			t.Errorf("SetEquals(%v) = %v, want %v", c.other, got, c.want)
		}
	}

	if !NewHashSet[int](nil).SetEquals(nil) {
		t.Error("empty set should equal an empty slice")
	}
}

func TestSetEqualsSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.SetEqualsSet(hs) {
		t.Error("a set should equal itself")
	}

	if !hs.SetEqualsSet(NewHashSet([]int{2, 1})) {
		t.Error("equal sets reported unequal")
	}

	if hs.SetEqualsSet(NewHashSet([]int{1, 3})) || hs.SetEqualsSet(nil) {
		t.Error("unequal sets reported equal")
	}

	if !NewHashSet[int](nil).SetEqualsSet(nil) {
		t.Error("two empty sets should be equal")
	}
}

func TestOverlaps(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.Overlaps([]int{5, 2}) {
		t.Error("sets sharing 2 reported no overlap")
	}

	if hs.Overlaps([]int{5, 6}) || hs.Overlaps(nil) {
		t.Error("disjoint collections reported overlap")
	}

	if NewHashSet[int](nil).Overlaps([]int{1}) {
		t.Error("an empty set overlaps nothing")
	}
}

func TestOverlapsSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.OverlapsSet(hs) || !hs.OverlapsSet(NewHashSet([]int{2, 3})) {
		t.Error("overlapping sets reported disjoint")
	}

	if hs.OverlapsSet(NewHashSet([]int{3})) || hs.OverlapsSet(nil) {
		t.Error("disjoint sets reported overlap")
	}
}

func TestIsSubsetOf(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	cases := []struct {
		other []int
		want  bool
	}{
		{[]int{1, 2}, true},
		{[]int{2, 1, 1, 2}, true},
		{[]int{1, 2, 3}, true},
		{[]int{1, 3, 4}, false},
		{[]int{1}, false},
		{nil, false},
	}

	for _, c := range cases {
		if got := hs.IsSubsetOf(c.other); got != c.want {
			t.Errorf("IsSubsetOf(%v) = %v, want %v", c.other, got, c.want)
		}
	}

	empty := NewHashSet[int](nil)

	if !empty.IsSubsetOf(nil) || !empty.IsSubsetOf([]int{1}) {
		t.Error("the empty set is a subset of every set")
	}
}

func TestIsSubsetOfSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if !hs.IsSubsetOfSet(hs) {
		t.Error("a set is a subset of itself")
	}

	if !hs.IsSubsetOfSet(NewHashSet([]int{1, 2, 3})) {
		t.Error("subset reported as not a subset")
	}

	if hs.IsSubsetOfSet(NewHashSet([]int{1, 3, 4})) || hs.IsSubsetOfSet(nil) {
		t.Error("non-subset reported as a subset")
	}
}

func TestIsProperSubsetOf(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	cases := []struct {
		other []int
		want  bool
	}{
		{[]int{1, 2, 3}, true},
		{[]int{3, 2, 1, 3}, true},
		{[]int{1, 2}, false},
		{[]int{1, 2, 2, 1}, false},
		{[]int{1, 3, 4}, false},
		{[]int{1}, false},
		{nil, false},
	}

	for _, c := range cases {
		if got := hs.IsProperSubsetOf(c.other); got != c.want {
			t.Errorf("IsProperSubsetOf(%v) = %v, want %v", c.other, got, c.want)
		}
	}

	empty := NewHashSet[int](nil)

	if empty.IsProperSubsetOf(nil) {
		t.Error("the empty set is not a proper subset of the empty set")
	}

	if !empty.IsProperSubsetOf([]int{1}) {
		t.Error("the empty set is a proper subset of any non-empty set")
	}
}

func TestIsProperSubsetOfSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2})

	if hs.IsProperSubsetOfSet(hs) {
		t.Error("a set is not a proper subset of itself")
	}

	if !hs.IsProperSubsetOfSet(NewHashSet([]int{1, 2, 3})) {
		t.Error("proper subset reported as not one")
	}

	if hs.IsProperSubsetOfSet(NewHashSet([]int{1, 3, 4})) || hs.IsProperSubsetOfSet(nil) {
		t.Error("non-subset reported as a proper subset")
	}
}

func TestIsSupersetOf(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	cases := []struct {
		other []int
		want  bool
	}{
		{nil, true},
		{[]int{}, true},
		{[]int{1, 3}, true},
		{[]int{1, 2, 3}, true},
		{[]int{1, 4}, false},
		{[]int{1, 2, 3, 4}, false},
	}

	for _, c := range cases {
		if got := hs.IsSupersetOf(c.other); got != c.want {
			t.Errorf("IsSupersetOf(%v) = %v, want %v", c.other, got, c.want)
		}
	}

	empty := NewHashSet[int](nil)

	if !empty.IsSupersetOf(nil) || empty.IsSupersetOf([]int{1}) {
		t.Error("empty set superset relations are wrong")
	}
}

func TestIsSupersetOfSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	if !hs.IsSupersetOfSet(hs) || !hs.IsSupersetOfSet(nil) {
		t.Error("a set is a superset of itself and of the empty set")
	}

	if !hs.IsSupersetOfSet(NewHashSet([]int{2})) {
		t.Error("superset reported as not one")
	}

	if hs.IsSupersetOfSet(NewHashSet([]int{2, 4})) {
		t.Error("non-superset reported as a superset")
	}
}

func TestIsProperSupersetOf(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	cases := []struct {
		other []int
		want  bool
	}{
		{nil, true},
		{[]int{1, 3}, true},
		{[]int{1, 2, 3}, false},
		{[]int{1, 2, 3, 4}, false},
		{[]int{1, 4}, false},
	}

	for _, c := range cases {
		if got := hs.IsProperSupersetOf(c.other); got != c.want {
			t.Errorf("IsProperSupersetOf(%v) = %v, want %v", c.other, got, c.want)
		}
	}
}

func TestIsProperSupersetOfSet(t *testing.T) {
	hs := NewHashSet([]int{1, 2, 3})

	if hs.IsProperSupersetOfSet(hs) {
		t.Error("a set is not a proper superset of itself")
	}

	if !hs.IsProperSupersetOfSet(NewHashSet([]int{3})) || !hs.IsProperSupersetOfSet(nil) {
		t.Error("proper superset reported as not one")
	}

	if hs.IsProperSupersetOfSet(NewHashSet([]int{3, 4})) {
		t.Error("non-superset reported as a proper superset")
	}
}

// TestSetOperationsAgainstModel checks the set-argument operations against a
// plain map model over many random inputs.
func TestSetOperationsAgainstModel(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	randomSet := func() HashSet[int] {
		hs := NewHashSet[int](nil)

		for i, n := 0, rng.Intn(6); i < n; i++ {
			hs.Add(rng.Intn(8))
		}

		return hs
	}

	for i := 0; i < 2000; i++ {
		a, b := randomSet(), randomSet()

		inA, inB := func(v int) bool { return a.Contains(v) }, func(v int) bool { return b.Contains(v) }

		var union, inter, diff, sym []int
		subset := true

		for v := 0; v < 8; v++ {
			if inA(v) || inB(v) {
				union = append(union, v)
			}

			if inA(v) && inB(v) {
				inter = append(inter, v)
			}

			if inA(v) && !inB(v) {
				diff = append(diff, v)
				subset = false
			}

			if inA(v) != inB(v) {
				sym = append(sym, v)
			}
		}

		equal := len(sym) == 0
		superset := len(b) == len(inter)
		overlaps := len(inter) > 0

		check := func(name string, op func(HashSet[int]), want []int) {
			t.Helper()

			c := NewHashSet(a.Keys())
			op(c)

			if got := sorted(c); !equalInts(got, want) {
				t.Fatalf("%s of %v and %v = %v, want %v", name, sorted(a), sorted(b), got, want)
			}
		}

		check("UnionWithSet", func(c HashSet[int]) { c.UnionWithSet(b) }, union)
		check("IntersectWithSet", func(c HashSet[int]) { c.IntersectWithSet(b) }, inter)
		check("ExceptWithSet", func(c HashSet[int]) { c.ExceptWithSet(b) }, diff)
		check("SymmetricExceptWithSet", func(c HashSet[int]) { c.SymmetricExceptWithSet(b) }, sym)

		if a.SetEqualsSet(b) != equal {
			t.Fatalf("SetEqualsSet(%v, %v) != %v", sorted(a), sorted(b), equal)
		}

		if a.OverlapsSet(b) != overlaps {
			t.Fatalf("OverlapsSet(%v, %v) != %v", sorted(a), sorted(b), overlaps)
		}

		if a.IsSubsetOfSet(b) != subset {
			t.Fatalf("IsSubsetOfSet(%v, %v) != %v", sorted(a), sorted(b), subset)
		}

		if a.IsProperSubsetOfSet(b) != (subset && !equal) {
			t.Fatalf("IsProperSubsetOfSet(%v, %v) != %v", sorted(a), sorted(b), subset && !equal)
		}

		if a.IsSupersetOfSet(b) != superset {
			t.Fatalf("IsSupersetOfSet(%v, %v) != %v", sorted(a), sorted(b), superset)
		}

		if a.IsProperSupersetOfSet(b) != (superset && !equal) {
			t.Fatalf("IsProperSupersetOfSet(%v, %v) != %v", sorted(a), sorted(b), superset && !equal)
		}
	}
}

// TestSliceArgumentsAreSets checks that a slice argument stands for the set of
// its distinct values, so repeating a value does not change a result.
func TestSliceArgumentsAreSets(t *testing.T) {
	one, oneTwo := NewHashSet([]int{1}), NewHashSet([]int{1, 2})

	if !one.SetEquals([]int{1, 1}) || oneTwo.SetEquals([]int{1, 1}) {
		t.Error("SetEquals should treat the slice as a set")
	}

	if !one.IsSupersetOf([]int{1, 1}) {
		t.Error("IsSupersetOf should treat the slice as a set")
	}

	if !oneTwo.IsProperSupersetOf([]int{1, 1}) {
		t.Error("IsProperSupersetOf should treat the slice as a set")
	}

	// expectSet stops the test, so this check stays last.
	sym := NewHashSet([]int{1})
	sym.SymmetricExceptWith([]int{1, 1, 2, 2})
	expectSet(t, sym, 2)
}

// TestEmptySetIsNotAProperSuperset checks that the empty set is not a proper
// superset of any set, the empty set included.
func TestEmptySetIsNotAProperSuperset(t *testing.T) {
	var zero HashSet[int] // nil map

	for _, empty := range []HashSet[int]{NewHashSet[int](nil), zero} {
		if empty.IsProperSupersetOf(nil) || empty.IsProperSupersetOf([]int{1}) {
			t.Error("IsProperSupersetOf: the empty set is not a proper superset of any set")
		}

		if empty.IsProperSupersetOfSet(nil) || empty.IsProperSupersetOfSet(NewHashSet([]int{1})) {
			t.Error("IsProperSupersetOfSet: the empty set is not a proper superset of any set")
		}
	}
}

// TestDifferential checks every slice form and every Set form exhaustively
// over small inputs. The receivers are every subset of {0, 1, 2}; the other
// collections are every slice of length 0 through 4 drawn from the same values,
// repeats included. Each slice form must give the same answer as its Set form
// applied to NewHashSet of the slice, and each Set form must give the same
// answer as a reference computed from plain maps, also when it is applied to
// its own receiver.
func TestDifferential(t *testing.T) {
	universe := []int{0, 1, 2}

	var receivers [][]int

	for mask := 0; mask < 1<<len(universe); mask++ {
		var r []int

		for i, v := range universe {
			if mask&(1<<i) != 0 {
				r = append(r, v)
			}
		}

		receivers = append(receivers, r)
	}

	others := [][]int{nil}

	for n, level := 1, [][]int{nil}; n <= 4; n++ {
		var next [][]int

		for _, s := range level {
			for _, v := range universe {
				next = append(next, append(append([]int(nil), s...), v))
			}
		}

		others = append(others, next...)
		level = next
	}

	// The reference works on plain maps and calls no HashSet method.
	toMap := func(values []int) map[int]bool {
		m := make(map[int]bool)

		for _, v := range values {
			m[v] = true
		}

		return m
	}

	// within reports whether every element of a is in b.
	within := func(a, b map[int]bool) bool {
		for v := range a {
			if !b[v] {
				return false
			}
		}

		return true
	}

	// pick returns, in ascending order, the values whose membership in a and
	// b satisfies keep.
	pick := func(a, b map[int]bool, keep func(inA, inB bool) bool) []int {
		var out []int

		for _, v := range universe {
			if keep(a[v], b[v]) {
				out = append(out, v)
			}
		}

		return out
	}

	type query struct {
		name  string
		slice func(hs HashSet[int], other []int) bool
		set   func(hs, other HashSet[int]) bool
		want  func(a, b map[int]bool) bool
	}

	queries := []query{
		{
			name:  "SetEquals",
			slice: func(hs HashSet[int], other []int) bool { return hs.SetEquals(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.SetEqualsSet(other) },
			want:  func(a, b map[int]bool) bool { return within(a, b) && within(b, a) },
		},
		{
			name:  "Overlaps",
			slice: func(hs HashSet[int], other []int) bool { return hs.Overlaps(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.OverlapsSet(other) },
			want: func(a, b map[int]bool) bool {
				return len(pick(a, b, func(inA, inB bool) bool { return inA && inB })) > 0
			},
		},
		{
			name:  "IsSubsetOf",
			slice: func(hs HashSet[int], other []int) bool { return hs.IsSubsetOf(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.IsSubsetOfSet(other) },
			want:  func(a, b map[int]bool) bool { return within(a, b) },
		},
		{
			name:  "IsProperSubsetOf",
			slice: func(hs HashSet[int], other []int) bool { return hs.IsProperSubsetOf(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.IsProperSubsetOfSet(other) },
			want:  func(a, b map[int]bool) bool { return within(a, b) && !within(b, a) },
		},
		{
			name:  "IsSupersetOf",
			slice: func(hs HashSet[int], other []int) bool { return hs.IsSupersetOf(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.IsSupersetOfSet(other) },
			want:  func(a, b map[int]bool) bool { return within(b, a) },
		},
		{
			name:  "IsProperSupersetOf",
			slice: func(hs HashSet[int], other []int) bool { return hs.IsProperSupersetOf(other) },
			set:   func(hs, other HashSet[int]) bool { return hs.IsProperSupersetOfSet(other) },
			want:  func(a, b map[int]bool) bool { return within(b, a) && !within(a, b) },
		},
	}

	type mutation struct {
		name  string
		slice func(hs HashSet[int], other []int)
		set   func(hs, other HashSet[int])
		want  func(a, b map[int]bool) []int
	}

	mutations := []mutation{
		{
			name:  "UnionWith",
			slice: func(hs HashSet[int], other []int) { hs.UnionWith(other) },
			set:   func(hs, other HashSet[int]) { hs.UnionWithSet(other) },
			want: func(a, b map[int]bool) []int {
				return pick(a, b, func(inA, inB bool) bool { return inA || inB })
			},
		},
		{
			name:  "IntersectWith",
			slice: func(hs HashSet[int], other []int) { hs.IntersectWith(other) },
			set:   func(hs, other HashSet[int]) { hs.IntersectWithSet(other) },
			want: func(a, b map[int]bool) []int {
				return pick(a, b, func(inA, inB bool) bool { return inA && inB })
			},
		},
		{
			name:  "ExceptWith",
			slice: func(hs HashSet[int], other []int) { hs.ExceptWith(other) },
			set:   func(hs, other HashSet[int]) { hs.ExceptWithSet(other) },
			want: func(a, b map[int]bool) []int {
				return pick(a, b, func(inA, inB bool) bool { return inA && !inB })
			},
		},
		{
			name:  "SymmetricExceptWith",
			slice: func(hs HashSet[int], other []int) { hs.SymmetricExceptWith(other) },
			set:   func(hs, other HashSet[int]) { hs.SymmetricExceptWithSet(other) },
			want: func(a, b map[int]bool) []int {
				return pick(a, b, func(inA, inB bool) bool { return inA != inB })
			},
		},
	}

	// Mismatches are gathered per method, so one run names every method that
	// disagrees instead of stopping at the first.
	type mismatch struct {
		count int
		first string
	}

	mismatches := make(map[string]*mismatch)
	comparisons := 0

	check := func(ok bool, name, format string, args ...interface{}) {
		comparisons++

		if ok {
			return
		}

		m := mismatches[name]

		if m == nil {
			m = &mismatch{first: fmt.Sprintf(format, args...)}
			mismatches[name] = m
		}

		m.count++
	}

	for _, r := range receivers {
		a := toMap(r)

		for _, o := range others {
			b := toMap(o)

			for _, q := range queries {
				bySlice := q.slice(NewHashSet(r), o)
				bySet := q.set(NewHashSet(r), NewHashSet(o))
				want := q.want(a, b)

				check(bySlice == bySet, q.name,
					"receiver %v, slice %v: %v, but %sSet gives %v", r, o, bySlice, q.name, bySet)
				check(bySet == want, q.name+"Set",
					"receiver %v, set of %v: %v, want %v", r, o, bySet, want)
			}

			for _, m := range mutations {
				bySlice, bySet := NewHashSet(r), NewHashSet(r)
				m.slice(bySlice, o)
				m.set(bySet, NewHashSet(o))
				want := m.want(a, b)

				check(equalInts(sorted(bySlice), sorted(bySet)), m.name,
					"receiver %v, slice %v: %v, but %sSet gives %v", r, o, sorted(bySlice), m.name, sorted(bySet))
				check(equalInts(sorted(bySet), want), m.name+"Set",
					"receiver %v, set of %v: %v, want %v", r, o, sorted(bySet), want)
			}
		}

		for _, q := range queries {
			hs := NewHashSet(r)
			got, want := q.set(hs, hs), q.want(a, a)

			check(got == want, q.name+"Set", "receiver %v with itself: %v, want %v", r, got, want)
		}

		for _, m := range mutations {
			hs := NewHashSet(r)
			m.set(hs, hs)
			want := m.want(a, a)

			check(equalInts(sorted(hs), want), m.name+"Set",
				"receiver %v with itself: %v, want %v", r, sorted(hs), want)
		}
	}

	names := make([]string, 0, len(mismatches))

	for name := range mismatches {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		m := mismatches[name]
		t.Errorf("%s: %d mismatches, first: %s", name, m.count, m.first)
	}

	t.Logf("%d receivers, %d slices, %d comparisons", len(receivers), len(others), comparisons)
}
