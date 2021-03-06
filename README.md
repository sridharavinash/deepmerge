Inspired by `deep_merge` in rails.

Typical usage:

```go
	map1 := map[string]int{
		"a": 10,
		"b": 3,
		"c": 12,
	}

	map2 := map[string]int{
		"a": 3,
		"b": 3,
		"c": 2,
	}

    add := func(a, b int) int { return a + b }
    d := &DeepMerge{
         map1: map1,
         map2: map2,
    }
    mergedMap, err := d.Merge(&add)
```

Expected Output(`mergedMap`):
```go
	 map[string]int{
		"a": 13,
		"b": 6,
		"c": 14,
	}
```

There are more examples in [deepmerge_test.go](https://github.com/sridharavinash/deepmerge/blob/master/deepmerge_test.go)

## Installation


```bash
go get -u github.com/sridharavinash/deepmerge
```

## Tests

```bash
go test ./...
```
