package deepmerge

import (
	"errors"
	"fmt"
	"reflect"
)

// DeepMerge instantiates initial counters / keys for traversal
type DeepMerge struct {
	// Stores the keys that we have processed as we iterate the maps
	seenKeys map[interface{}]bool

	// Keeps track of nested parent keys
	parentKey reflect.Value
}

// Merge merges 2 maps by applying a fptr (function pointer) to the values
// Example: Merge(map1, map2, &my_func)
// where my_func := func(a,b int) int { return a + b }
func (d DeepMerge) Merge(m1, m2, fptr interface{}) (interface{}, error) {

	if m1 == nil && m2 != nil {
		return m2, nil
	}

	if m1 != nil && m2 == nil {
		return m1, nil
	}

	m1_t := reflect.ValueOf(m1).Type()
	m2_t := reflect.ValueOf(m2).Type()
	if m1_t != m2_t {
		return nil, errors.New("Maps have to be of the same type")
	}

	d.seenKeys = make(map[interface{}]bool)
	return d.merge(m1, m2, fptr), nil
}

func (d DeepMerge) merge(m1, m2, fptr interface{}) interface{} {
	// Lets keep track of the keys from the maps were iterating
	var allKeys []reflect.Value

	m1_t := reflect.ValueOf(m1)
	m2_t := reflect.ValueOf(m2)

	// This will store our final merged map
	ret_map := reflect.MakeMap(m1_t.Type())

	cp_m1 := reflect.New(m1_t.Type()).Elem()
	cp_m2 := reflect.New(m2_t.Type()).Elem()

	// Copy over the map values to the placeholder maps
	// that will perform the fptr function operations
	translateRecursive(cp_m1, m1_t)
	translateRecursive(cp_m2, m2_t)

	// Lets find out what type of function we have
	// so that we can call it
	fn := reflect.ValueOf(fptr).Elem()

	allKeys = append(allKeys, cp_m1.MapKeys()...)
	allKeys = append(allKeys, cp_m2.MapKeys()...)

	// For each key we'll run the function block
	for _, k := range allKeys {
		// If we've already processed the key we'll skip
		if _, ok := d.seenKeys[k.Interface()]; ok {
			continue
		}
		// If we're traversing a parent_key, then we'll need to add that key
		if (d.parentKey.IsValid()) && (d.parentKey.Len() != 0) {
			keyplus := fmt.Sprintf("%v_%v", k.Interface(), d.parentKey.Interface())
			d.seenKeys[keyplus] = true

		} else {
			d.seenKeys[k.Interface()] = true
		}

		// Get the value of the key from each map
		v := cp_m1.MapIndex(k)
		o_v := cp_m2.MapIndex(k)

		// If we have a map, lets iterative through
		// recursively
		if v.Kind() == reflect.Map && o_v.Kind() == reflect.Map {
			d.parentKey = k
			yy := d.merge(v.Interface(), o_v.Interface(), fptr)
			ret_map.SetMapIndex(k, reflect.ValueOf(yy))
		} else {
			// If any of the keys traversed is invalid
			// we'll ignore it and update the map to
			// the values of the other key
			if !v.IsValid() && o_v.IsValid() {
				ret_map.SetMapIndex(k, o_v)
				continue
			}
			if !o_v.IsValid() && v.IsValid() {
				ret_map.SetMapIndex(k, v)
				continue
			}

			// Everything looks good to call the function pointer
			in := []reflect.Value{v, o_v}
			zz := fn.Call(in)
			// We'll only take the first return value
			ret_map.SetMapIndex(k, zz[0])
		}
	}
	return ret_map.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}
}
