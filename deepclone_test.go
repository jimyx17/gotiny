// deepcopy makes deep copies of things. A standard copy will copy the
// pointers: deep copy copies the values pointed to.  Unexported field
// values are not copied.
//
// Copyright (c)2014-2016, Joel Scoble (github.com/mohae), all rights reserved.
// License: MIT, for more details check the included LICENSE file.
package gotiny_test

import (
	"reflect"
	"unsafe"
)

// Copy creates a deep copy of whatever is passed to it and returns the copy
// in an interface{}.  The returned value will need to be asserted to the
// correct type.
func DeepClone(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	// Make the interface a reflect.Value
	original := reflect.ValueOf(src)

	// Make a copy of the same type as the original.
	cpy := reflect.New(original.Type()).Elem()

	// Recursively copy the original.
	copyRecursive(original, cpy)

	//fmt.Println("original ", original)
	//fmt.Println("copy", cpy)
	// Return the copy as an interface.
	return cpy.Interface()
}

// copyRecursive does the actual copying of the interface. It currently has
// limited support for what it can handle. Add as needed.
func copyRecursive(original, cpy reflect.Value) {
	// handle according to original's Kind
	switch original.Kind() {
	case reflect.Ptr:
		// Get the actual value being pointed to.
		originalValue := original.Elem()

		// if  it isn't valid, return.
		if !originalValue.IsValid() {
			return
		}
		cpy.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, cpy.Elem())

	case reflect.Interface:
		// If this is a nil, don't do anything
		if original.IsNil() {
			return
		}
		// Get the value for the interface, not the pointer.
		originalValue := original.Elem()

		// Get the value by calling Elem().
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		cpy.Set(copyValue)

	case reflect.Struct:
		for i := 0; i < original.NumField(); i++ {
			field := cpy.Type().Field(i)
			copyRecursive(original.Field(i),
				reflect.NewAt(field.Type, unsafe.Pointer(cpy.Field(i).UnsafeAddr())).Elem())
		}

	case reflect.Slice:
		if original.IsNil() {
			return
		}
		// Make a new slice and copy each element.
		cpy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i))
		}

	case reflect.Map:
		if original.IsNil() {
			return
		}
		cpy.Set(reflect.MakeMap(original.Type()))
		for _, oKey := range original.MapKeys() {
			cKey := reflect.New(oKey.Type()).Elem()
			oVal := original.MapIndex(oKey)
			cVal := reflect.New(oVal.Type()).Elem()
			copyRecursive(oKey, cKey)
			copyRecursive(oVal, cVal)
			cpy.SetMapIndex(cKey, cVal)
		}
	case reflect.Array:
		for i := 0; i < original.Len(); i++ {
			copyRecursive(original.Index(i), cpy.Index(i))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cpy.SetInt(original.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		cpy.SetUint(original.Uint())
	case reflect.Bool:
		cpy.SetBool(original.Bool())
	case reflect.Float32, reflect.Float64:
		cpy.SetFloat(original.Float())
	case reflect.Complex64, reflect.Complex128:
		cpy.SetComplex(original.Complex())
	case reflect.String:
		cpy.SetString(original.String())
	case reflect.UnsafePointer:
		cpy.SetPointer(unsafe.Pointer(original.Pointer()))
	}
}
