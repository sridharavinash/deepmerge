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
