package deepmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_basicAddition_depth1(t *testing.T) {
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

	expect := map[string]int{
		"a": 13,
		"b": 6,
		"c": 14,
	}

	add := func(a, b int) int { return a + b }

	got := merge(map1, map2, &add)
	assert.Equal(t, expect, got)
}

func Test_basicsubtraction_depth1(t *testing.T) {
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

	expect := map[string]int{
		"a": 7,
		"b": 0,
		"c": 10,
	}

	sub := func(a, b int) int { return a - b }

	got := merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}

func Test_basicsubtraction_nested(t *testing.T) {

	map1 := map[string]map[string]int{
		"a": map[string]int{"a1": 1, "a2": 2, "a3": -1},
		"b": map[string]int{"b1": 1, "b2": 2},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	map2 := map[string]map[string]int{
		"a": map[string]int{"a1": 3, "a2": 1},
		"b": map[string]int{"b1": 1, "b2": 2},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	expect := map[string]map[string]int{
		"a": map[string]int{"a1": -2, "a2": 1, "a3": -1},
		"b": map[string]int{"b1": 0, "b2": 0},
		"c": map[string]int{"c1": 0, "c2": 0},
	}

	sub := func(a, b int) int { return a - b }

	got := merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}
