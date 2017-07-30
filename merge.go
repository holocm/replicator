/*******************************************************************************
*
* Copyright 2017 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT ANY
* WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
* A PARTICULAR PURPOSE. See the GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"fmt"
	"reflect"
)

//MergeError is the type of error that's returned by MergeTables.
type MergeError struct {
	Location string
	Kind     string
	Message  string
}

//Error implements the error interface.
func (e MergeError) Error() string {
	return fmt.Sprintf("%s in .Vars%s: %s", e.Kind, e.Location, e.Message)
}

//MergeTables merges two maps in a way that is useful for Replicator.
//See README for details.
func MergeTables(first, second map[string]interface{}) (map[string]interface{}, error) {
	result, me := doMergeTables(first, second)
	if me == nil {
		return result, nil
	}
	return result, me
}

func doMergeTables(first, second map[string]interface{}) (map[string]interface{}, *MergeError) {
	if first == nil {
		return second, nil
	}
	if second == nil {
		return first, nil
	}
	//After this point, first != nil && second != nil.

	result := map[string]interface{}{}

	for key, val1 := range first {
		result[key] = val1
	}

	for key, val2 := range second {
		val1, exists := result[key]
		if exists {
			var err *MergeError
			result[key], err = mergeValues(
				reflect.ValueOf(val1),
				reflect.ValueOf(val2),
			)
			if err != nil {
				return nil, &MergeError{
					Location: "." + key + err.Location,
					Kind:     err.Kind,
					Message:  err.Message,
				}
			}
		} else {
			result[key] = val2
		}
	}

	return result, nil
}

func mergeValues(first, second reflect.Value) (interface{}, *MergeError) {
	//deference pointers, if any
	switch first.Kind() {
	case reflect.Ptr, reflect.Interface:
		return mergeValues(first.Elem(), second)
	}
	switch second.Kind() {
	case reflect.Ptr, reflect.Interface:
		return mergeValues(first, second.Elem())
	}

	kind1, err := simplifiedKindOf(first)
	if err != nil {
		return nil, err
	}
	kind2, err := simplifiedKindOf(second)
	if err != nil {
		return nil, err
	}

	if kind1 != kind2 {
		return nil, &MergeError{
			Kind:    "type mismatch",
			Message: kind1 + " <> " + kind2,
		}
	}

	switch kind1 {
	case "scalar":
		//second value takes precedence over first value
		return second.Interface(), nil
	case "array":
		//arrays/slices are concatenated
		len1 := first.Len()
		len2 := second.Len()
		result := make([]interface{}, 0, len1+len2)
		for idx := 0; idx < len1; idx++ {
			result = append(result, first.Index(idx).Interface())
		}
		for idx := 0; idx < len2; idx++ {
			result = append(result, second.Index(idx).Interface())
		}
		return result, nil
	case "table":
		//maps are merged recursively with MergeTables, but they need to be
		//converted to map[string]interface{} first
		map1, err := coerceMap(first)
		if err != nil {
			return nil, err
		}
		map2, err := coerceMap(second)
		if err != nil {
			return nil, err
		}
		return doMergeTables(map1, map2)
	default:
		panic("unreachable: " + kind1)
	}
}

func simplifiedKindOf(rv reflect.Value) (string, *MergeError) {
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr, reflect.Struct, reflect.UnsafePointer:
		return "", &MergeError{Kind: "unsupported value type", Message: rv.Type().String()}
	case reflect.Array, reflect.Slice:
		return "array", nil
	case reflect.Map:
		return "table", nil
	default:
		return "scalar", nil
	}
}

func coerceMap(rv reflect.Value) (map[string]interface{}, *MergeError) {
	if rv.Type().Key().String() != "string" {
		return nil, &MergeError{
			Kind:    "unsupported key type",
			Message: rv.Type().Key().String(),
		}
	}
	result := make(map[string]interface{}, rv.Len())
	for _, key := range rv.MapKeys() {
		result[key.Interface().(string)] = rv.MapIndex(key).Interface()
	}
	return result, nil
}
