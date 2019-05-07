package deepmerge

import (
	"fmt"
	"reflect"
)

func merge(m1, m2, fptr interface{}) interface{} {
	m1_t := reflect.ValueOf(m1)
	m2_t := reflect.ValueOf(m2)
	ret_map := reflect.MakeMap(m1_t.Type())

	cp_m1 := reflect.New(m1_t.Type()).Elem()
	cp_m2 := reflect.New(m2_t.Type()).Elem()
	translateRecursive(cp_m1, m1_t)
	translateRecursive(cp_m2, m2_t)
	fn := reflect.ValueOf(fptr).Elem()
	//fmt.Println(cp_m1, cp_m1.Interface(), cp_m2, ret_map, fn)
	//fmt.Printf("%T | %T | %T | %T | %T\n", cp_m1, cp_m1.Interface(), cp_m2, ret_map, fn)
	for _, k := range cp_m1.MapKeys() {
		v := cp_m1.MapIndex(k)
		o_v := cp_m2.MapIndex(k)
		fmt.Println(k, v, o_v)
		if v.Kind() == reflect.Map || o_v.Kind() == reflect.Map {
			yy := merge(v.Interface(), o_v.Interface(), fptr)
			fmt.Printf("yy: %[1]v | %[1]T\n", yy)
			ret_map.SetMapIndex(k, reflect.ValueOf(yy))
		} else {
			if v.Kind() == reflect.Invalid {
				ret_map.SetMapIndex(k, o_v)
				continue
			}
			if o_v.Kind() == reflect.Invalid {
				ret_map.SetMapIndex(k, v)
				continue
			}

			fmt.Println(v.Kind(), o_v.Kind())
			in := []reflect.Value{v, o_v}
			fmt.Printf("%[1]v, %[1]T\n", in)
			zz := fn.Call(in)
			fmt.Printf("Rsult: %v\n", zz[0])
			ret_map.SetMapIndex(k, zz[0])
		}
	}
	fmt.Printf("returning Map: %[1]v | %[1]T\n", ret_map)
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
