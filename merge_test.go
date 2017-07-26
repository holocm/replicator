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
	actual, err := MergeMaps(first, second)
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

type tbl map[string]interface{}
type ary []interface{}

func TestMergeMaps(t *testing.T) {
	testcase(t, "merge tables with scalars",
		tbl{"bool1": true, "int1": 42, "str1": "hallo"},
		tbl{"bool2": false, "int1": 23, "str1": 5},
		tbl{"bool1": true, "bool2": false, "int1": 23, "str1": 5},
		"",
	)
	testcase(t, "merge tables with arrays",
		tbl{
			"scalar1": false,
			"array":   ary{1, 2, 3},
		},
		tbl{
			"scalar2": true,
			"array":   ary{4, 5, 6},
		},
		tbl{
			"scalar1": false,
			"scalar2": true,
			"array":   ary{1, 2, 3, 4, 5, 6},
		},
		"",
	)
	testcase(t, "merge tables with tables",
		tbl{
			"foo": tbl{"a": 1, "b": 2},
			"bar": tbl{"c": 3, "d": 4},
		},
		tbl{
			"foo": tbl{"a": 10, "c": 30},
		},
		tbl{
			"foo": tbl{"a": 10, "b": 2, "c": 30},
			"bar": tbl{"c": 3, "d": 4},
		},
		"",
	)
	testcase(t, "type mismatch: scalar <> array",
		tbl{"foo": tbl{"bar": 5, "baz": 23}},
		tbl{"foo": tbl{"bar": ary{1, 2, 3}}},
		nil,
		"type mismatch in .Vars.foo: scalar <> array",
	)
}
