# hashset

A generic set for Go, built on a Go map, with a method surface that mirrors
.NET's `HashSet<T>`. It has no dependencies outside the Go standard library.

```go
type HashSet[T comparable] map[T]void
```

A `HashSet` is a map, so `len(hs)` and `for v := range hs` work directly.
Element order is not defined. The zero value is a nil map: reads work, but
adding to it panics, so create sets with `NewHashSet`.

## Install

```sh
go get github.com/ritchiecarroll/hashset
```

Requires Go 1.18 or later.

## Usage

```go
package main

import (
	"fmt"

	"github.com/ritchiecarroll/hashset"
)

func main() {
	langs := hashset.NewHashSet([]string{"go", "cs", "go"})

	langs.Add("ts")
	langs.Remove("cs")

	fmt.Println(len(langs), langs.Contains("go")) // 2 true

	other := hashset.NewHashSet([]string{"go", "rs"})
	langs.IntersectWithSet(other)

	fmt.Println(langs.Keys()) // [go]
}
```

## Methods

Most operations come in two forms: one takes a slice, and one with a `Set`
suffix takes another `HashSet`. A slice argument stands for the set of its
distinct values, so repeating a value in it does not change the result.

| Method | Purpose |
|:--|:--|
| `NewHashSet(items)` | New set holding the distinct items of a slice |
| `Add`, `Remove`, `RemoveWhere` | Add or remove elements; report what changed |
| `Contains`, `IsEmpty`, `Keys`, `Clear` | Query, copy out, or empty the set |
| `UnionWith`, `UnionWithSet` | Add every element of the other collection |
| `IntersectWith`, `IntersectWithSet` | Keep only elements also in the other collection |
| `ExceptWith`, `ExceptWithSet` | Remove every element of the other collection |
| `SymmetricExceptWith`, `SymmetricExceptWithSet` | Keep elements in exactly one of the two |
| `SetEquals`, `SetEqualsSet` | Same elements? |
| `Overlaps`, `OverlapsSet` | Any element in common? |
| `IsSubsetOf`, `IsSubsetOfSet`, `IsProperSubsetOf`, `IsProperSubsetOfSet` | Subset tests |
| `IsSupersetOf`, `IsSupersetOfSet`, `IsProperSupersetOf`, `IsProperSupersetOfSet` | Superset tests |

## Thread safety

The methods are not thread-safe. To use a set from more than one goroutine,
hold a lock around every call.

## Origin

HashSet was extracted from [go2cs](https://github.com/ritchiecarroll/go2cs), a
Go to C# transpiler.

## License

[MIT](LICENSE)
