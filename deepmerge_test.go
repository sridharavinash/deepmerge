package deepmerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_map1_is_nil(t *testing.T) {

	map2 := map[string]int{
		"a": 10,
		"b": 3,
		"c": 12,
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, err := d.Merge(nil, map2, &sub)
	assert.Equal(t, got, map2)
	assert.Nil(t, err)
}

func Test_map2_is_nil(t *testing.T) {

	map1 := map[string]int{
		"a": 10,
		"b": 3,
		"c": 12,
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, err := d.Merge(map1, nil, &sub)
	assert.Equal(t, got, map1)
	assert.Nil(t, err)
}

func Test_fptr_is_basic_update(t *testing.T) {

	map1 := map[string]int{
		"a": 10,
		"b": 3,
		"c": 12,
	}

	map2 := map[string]int{
		"a": 20,
		"b": 32,
		"c": 22,
	}
	fp := func(a, b int) int { return b }
	d := &DeepMerge{}
	got, err := d.Merge(map1, map2, &fp)
	assert.Equal(t, got, map2)
	assert.Nil(t, err)
}

func Test_maps_are_not_the_same_type(t *testing.T) {

	map1 := map[string]map[string]int{
		"a": map[string]int{"a1": 1, "a2": 2, "a3": -1},
		"b": map[string]int{"a1": 1, "a2": 2},
	}

	map2 := map[string]int{
		"a": 10,
		"b": 3,
		"c": 12,
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, err := d.Merge(map1, map2, &sub)
	assert.Nil(t, got)
	assert.Error(t, err)
}

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
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &add)
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
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &sub)
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
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}

func Test_nested_map1_missing_key(t *testing.T) {

	map1 := map[string]map[string]int{
		"a": map[string]int{"a1": 1, "a2": 2, "a3": -1},
	}

	map2 := map[string]map[string]int{
		"a": map[string]int{"a1": 3, "a2": 1},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	expect := map[string]map[string]int{
		"a": map[string]int{"a1": -2, "a2": 1, "a3": -1},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}
func Test_nested_both_maps_missing_key(t *testing.T) {

	map1 := map[string]map[string]int{
		"a": map[string]int{"a1": 1, "a2": 2, "a3": -1},
		"b": map[string]int{"b1": 1, "b2": 2},
	}

	map2 := map[string]map[string]int{
		"a": map[string]int{"a1": 3, "a2": 1},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	expect := map[string]map[string]int{
		"a": map[string]int{"a1": -2, "a2": 1, "a3": -1},
		"b": map[string]int{"b1": 1, "b2": 2},
		"c": map[string]int{"c1": 1, "c2": 2},
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}

func Test_nested_same_sub_keys(t *testing.T) {

	map1 := map[string]map[string]int{
		"a": map[string]int{"a1": 1, "a2": 2, "a3": -1},
		"b": map[string]int{"a1": 1, "a2": 2},
	}

	map2 := map[string]map[string]int{
		"a": map[string]int{"a1": 3, "a2": 1},
		"b": map[string]int{"a1": 1, "a2": 2},
		"c": map[string]int{"a1": 1, "a2": 2},
	}

	expect := map[string]map[string]int{
		"a": map[string]int{"a1": -2, "a2": 1, "a3": -1},
		"b": map[string]int{"a1": 0, "a2": 0},
		"c": map[string]int{"a1": 1, "a2": 2},
	}

	sub := func(a, b int) int { return a - b }
	d := &DeepMerge{}
	got, _ := d.Merge(map1, map2, &sub)
	assert.Equal(t, expect, got)
}
