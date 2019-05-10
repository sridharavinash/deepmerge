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
	d := &DeepMerge{}
    mergedMap := d.Merge(map1, map2, &add)
```

Expected Output(`mergedMap`):
```go
	 map[string]int{
		"a": 13,
		"b": 6,
		"c": 14,
	}
```

There are more examples in `deepmerge_test.go`
