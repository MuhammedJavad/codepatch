# Array Utilities (Go 1.23)

A collection of functional array utilities for Go with generic support, providing common operations like mapping, filtering, finding, and aggregating.

## Features

- **Generic Support**: Works with any type using Go generics
- **Functional Style**: Immutable operations that don't modify original slices
- **Comprehensive**: Map, filter, find, aggregate, and utility functions
- **Well Tested**: Full test coverage with standard library testing
- **Zero Dependencies**: Pure Go with no external dependencies

## Installation

```bash
go get github.com/MuhammedJavad/codepatch/array
```

## Quick Start

```go
package main

import (
    "fmt"
    "log/slog"
    "github.com/MuhammedJavad/codepatch/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Map: transform each element
    doubled := array.Project(numbers, func(val, index int) int {
        return val * 2
    })
    fmt.Println("Doubled:", doubled) // [2 4 6 8 10 12 14 16 18 20]
    
    // Filter: keep only even numbers
    evens := array.Filter(numbers, func(val *int) bool {
        return *val%2 == 0
    })
    fmt.Println("Evens:", evens) // [2 4 6 8 10]
    
    // Find: get first element matching condition
    found := array.Find(numbers, func(val int) bool {
        return val > 5
    })
    if found != nil {
        fmt.Println("First > 5:", *found) // 6
    }
    
    // Check conditions
    allPositive := array.All(numbers, func(val int) bool {
        return val > 0
    })
    fmt.Println("All positive:", allPositive) // true
    
    hasEven := array.Any(numbers, func(val int) bool {
        return val%2 == 0
    })
    fmt.Println("Has even:", hasEven) // true
}
```

## API Reference

### Mapping Operations

#### `Project[TIn, TOut](arr []TIn, selector func(TIn, int) TOut) []TOut`
Transform each element using a selector function.

```go
squares := array.Project([]int{1, 2, 3}, func(val, index int) int {
    return val * val
})
// Result: [1, 4, 9]
```

#### `ProjectErr[TIn, TOut](arr []TIn, selector func(TIn, int) (*TOut, error)) ([]TOut, error)`
Transform with error handling. Returns error if any transformation fails.

```go
result, err := array.ProjectErr([]string{"1", "2", "abc"}, func(val string, index int) (*int, error) {
    return strconv.Atoi(val)
})
// Error on "abc"
```

#### `FilterAndProject[TIn, TOut](arr []TIn, selector func(TIn, int) *TOut) []TOut`
Map and filter in one operation. Nil results are filtered out.

```go
validNumbers := array.FilterAndProject([]string{"1", "abc", "3"}, func(val string, index int) *int {
    if num, err := strconv.Atoi(val); err == nil {
        return &num
    }
    return nil
})
// Result: [1, 3]
```

### Filtering Operations

#### `Filter[T](arr []T, predicate func(*T) bool) []T`
Keep elements that match the predicate.

```go
evens := array.Filter([]int{1, 2, 3, 4}, func(val *int) bool {
    return *val%2 == 0
})
// Result: [2, 4]
```

### Finding Operations

#### `Find[T](arr []T, predicate func(T) bool) *T`
Find first element matching predicate. Returns nil if not found.

```go
found := array.Find([]string{"hello", "world", "go"}, func(val string) bool {
    return len(val) > 3
})
// Result: &"hello"
```

### Conditional Operations

#### `All[T](arr []T, predicate func(T) bool) bool`
Check if all elements match predicate.

```go
allPositive := array.All([]int{1, 2, 3}, func(val int) bool {
    return val > 0
})
// Result: true
```

#### `Any[T](arr []T, predicate func(T) bool) bool`
Check if any element matches predicate.

```go
hasEven := array.Any([]int{1, 3, 5, 6}, func(val int) bool {
    return val%2 == 0
})
// Result: true
```

#### `Contains[T comparable](slice []T, element T) bool`
Check if slice contains element using equality comparison.

```go
hasGo := array.Contains([]string{"go", "rust", "python"}, "go")
// Result: true
```

### Error Handling

#### `AnyErr[T](arr []T, predicate func(T) error) error`
Apply predicate to each element, return first error encountered.

```go
err := array.AnyErr([]string{"file1.txt", "file2.txt"}, func(filename string) error {
    return os.Stat(filename)
})
// Returns first file not found error
```

### Aggregation Operations

#### `Sum[T, N Numbers](arr []T, selector func(T) N) N`
Sum values using selector function. Works with numeric types.

```go
type Product struct {
    Name  string
    Price float64
}

products := []Product{
    {Name: "Laptop", Price: 999.99},
    {Name: "Mouse", Price: 29.99},
}

total := array.Sum(products, func(p Product) float64 {
    return p.Price
})
// Result: 1029.98
```

#### `Chunk[T](arr []T, chunkSize int) [][]T`
Split array into chunks of specified size.

```go
chunks := array.Chunk([]int{1, 2, 3, 4, 5}, 2)
// Result: [[1, 2], [3, 4], [5]]
```

### Map Operations

#### `ToMap[TKey comparable, TValue](in []TValue, keySelector func(TValue) TKey) map[TKey]TValue`
Convert slice to map using key selector.

```go
type User struct {
    ID   int
    Name string
}

users := []User{{1, "Alice"}, {2, "Bob"}}
userMap := array.ToMap(users, func(u User) int { return u.ID })
// Result: map[1:{1 Alice} 2:{2 Bob}]
```

#### `ProjectMap[TKey comparable, TIn, TOut](m map[TKey]TIn, selector func(TIn, TKey) TOut) []TOut`
Transform map entries to slice.

```go
scores := map[string]int{"Alice": 95, "Bob": 87}
names := array.ProjectMap(scores, func(score int, name string) string {
    return name
})
// Result: ["Alice", "Bob"] (order not guaranteed)
```

#### `FlatMap[TKey comparable, TIn, TC, TOut](m map[TKey]TIn, collectionSelector func(TIn) []TC, resultSelector func(TC) TOut) []TOut`
Flatten map values into single slice.

```go
departments := map[string][]string{
    "Engineering": {"Alice", "Bob"},
    "Marketing":   {"Charlie"},
}
allNames := array.FlatMap(departments, 
    func(employees []string) []string { return employees },
    func(name string) string { return name },
)
// Result: ["Alice", "Bob", "Charlie"] (order not guaranteed)
```

### Flattening Operations

#### `Flat[TIn, TC, TOut](arr []TIn, collectionSelector func(TIn) []TC, resultSelector func(TIn, TC) *TOut) []TOut`
Flatten nested collections, filtering out nil results.

```go
type Order struct {
    Items []string
}

orders := []Order{
    {Items: []string{"book", "pen"}},
    {Items: []string{"laptop"}},
}

allItems := array.Flat(orders,
    func(order Order) []string { return order.Items },
    func(order Order, item string) *string { return &item },
)
// Result: ["book", "pen", "laptop"]
```

### Utility Functions

#### `IsEmpty[T](arr []T) bool`
Check if slice is empty (nil or zero length).

```go
empty := array.IsEmpty([]int{})
// Result: true

empty = array.IsEmpty(nil)
// Result: true
```

## Type Constraints

The package uses Go's type constraints for numeric operations:

```go
type Numbers interface {
    int | int8 | int16 | int32 | int64 | float32 | float64
}
```

## Performance Notes

- All operations create new slices; original data is never modified
- Operations are optimized for readability and correctness
- For very large datasets, consider streaming approaches
- `for range` loops are used for better performance and bounds checking

## Examples

### Working with Structs

```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
    {Name: "Charlie", Age: 35},
}

// Find adults
adults := array.Filter(people, func(p *Person) bool {
    return p.Age >= 18
})

// Get names
names := array.Project(people, func(p Person, index int) string {
    return p.Name
})

// Check if all are adults
allAdults := array.All(people, func(p Person) bool {
    return p.Age >= 18
})
```

### Error Handling Patterns

```go
// Process files with error handling
files := []string{"file1.txt", "file2.txt", "invalid.txt"}

results := array.FilterAndProject(files, func(filename string, index int) *string {
    if _, err := os.Stat(filename); err == nil {
        return &filename
    }
    return nil
})
// Only valid files in results
```

### Complex Transformations

```go
type Order struct {
    ID     int
    Items  []Item
    Total  float64
}

type Item struct {
    Name  string
    Price float64
}

orders := []Order{/* ... */}

// Get all item names from all orders
allItemNames := array.Flat(orders,
    func(order Order) []Item { return order.Items },
    func(order Order, item Item) *string { return &item.Name },
)

// Calculate total revenue
totalRevenue := array.Sum(orders, func(order Order) float64 {
    return order.Total
})
```

## Testing

Run tests with:

```bash
go test ./...
```

The package includes comprehensive tests for all functions with edge cases and error conditions.

## License

This package is part of the codepatch project. See the main project for license information.
