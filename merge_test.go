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
	"reflect"
	"testing"
)

func testcase(t *testing.T, name string, first, second, expected map[string]interface{}, expectedErr string) {
	actual, err := MergeTables(first, second)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"testcase failed: %s\n  expected result: %#v\n    actual result: %#v\n",
			name, expected, actual,
		)
	}
	var actualErr string
	if err != nil {
		actualErr = err.Error()
	}
	if expectedErr != actualErr {
		t.Errorf(
			"testcase failed: %s\n  expected error: %#v\n    actual error: %#v\n",
			name, expectedErr, actualErr,
		)
	}
}

func TestMergeTables(t *testing.T) {
	testcase(t, "merge tables with scalars",
		map[string]interface{}{"bool1": true, "int1": 42, "str1": "hallo"},
		map[string]interface{}{"bool2": false, "int1": 23, "str1": 5},
		map[string]interface{}{"bool1": true, "bool2": false, "int1": 23, "str1": 5},
		"",
	)
	testcase(t, "merge tables with arrays",
		map[string]interface{}{
			"scalar1": false,
			"array":   []interface{}{1, 2, 3},
		},
		map[string]interface{}{
			"scalar2": true,
			"array":   []interface{}{4, 5, 6},
		},
		map[string]interface{}{
			"scalar1": false,
			"scalar2": true,
			"array":   []interface{}{1, 2, 3, 4, 5, 6},
		},
		"",
	)
	testcase(t, "merge tables with tables",
		map[string]interface{}{
			"foo": map[string]interface{}{"a": 1, "b": 2},
			"bar": map[string]interface{}{"c": 3, "d": 4},
		},
		map[string]interface{}{
			"foo": map[string]interface{}{"a": 10, "c": 30},
		},
		map[string]interface{}{
			"foo": map[string]interface{}{"a": 10, "b": 2, "c": 30},
			"bar": map[string]interface{}{"c": 3, "d": 4},
		},
		"",
	)
	testcase(t, "type mismatch: scalar <> array",
		map[string]interface{}{"foo": map[string]interface{}{"bar": 5, "baz": 23}},
		map[string]interface{}{"foo": map[string]interface{}{"bar": []interface{}{1, 2, 3}}},
		nil,
		"type mismatch in .Vars.foo.bar: scalar <> array",
	)
}
