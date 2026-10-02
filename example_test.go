// SPDX-License-Identifier: MIT
// Copyright (c) 2021-2026 The go2cs Authors

package hashset_test

import (
	"fmt"
	"sort"

	"github.com/ritchiecarroll/hashset"
)

func ExampleNewHashSet() {
	hs := hashset.NewHashSet([]string{"go", "cs", "go"})

	fmt.Println(len(hs), hs.Contains("go"), hs.Contains("ts"))
	// Output: 2 true false
}

func ExampleHashSet_Add() {
	hs := hashset.NewHashSet[int](nil)

	fmt.Println(hs.Add(1), hs.Add(1))
	// Output: true false
}

func ExampleHashSet_UnionWithSet() {
	hs := hashset.NewHashSet([]int{1, 2})
	hs.UnionWithSet(hashset.NewHashSet([]int{2, 3}))

	keys := hs.Keys()
	sort.Ints(keys)
	fmt.Println(keys)
	// Output: [1 2 3]
}

func ExampleHashSet_IntersectWith() {
	hs := hashset.NewHashSet([]int{1, 2, 3, 4})
	hs.IntersectWith([]int{2, 4, 6})

	keys := hs.Keys()
	sort.Ints(keys)
	fmt.Println(keys)
	// Output: [2 4]
}

func ExampleHashSet_IsSubsetOfSet() {
	small := hashset.NewHashSet([]int{1, 2})
	large := hashset.NewHashSet([]int{1, 2, 3})

	fmt.Println(small.IsSubsetOfSet(large), small.IsProperSubsetOfSet(large), large.IsSubsetOfSet(small))
	// Output: true true false
}
