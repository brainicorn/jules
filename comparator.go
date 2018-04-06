/*
This file was originally written for the grules library.
https://github.com/huttotw/grules

Copyright © 2018 Trevor Hutto
Licensed under the Apache License, Version 2.0 (the "License")
*/
package jules

import (
	"reflect"
)

var (
	comparators = map[string]Comparator{
		"eq":     equal,
		"neq":    notEqual,
		"gt":     greaterThan,
		"gte":    greaterThanEqual,
		"lt":     lessThan,
		"lte":    lessThanEqual,
	}
)

// Comparator is a function that should evaluate two values and return
// the true if the comparison is true, or false if the comparison is
// false
type Comparator func(a, b interface{}) bool

// equal will return true if a == b
func equal(a, b interface{}) bool {
	return a == b
}

// notEqual will return true if a != b
func notEqual(a, b interface{}) bool {
	return !equal(a, b)
}

// lessThan will return true if a < b
func lessThan(a, b interface{}) bool {
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// Make sure the types are the same
	if ta != tb {
		return false
	}

	// We have already checked that each argument is the same type
	// so it is safe to only check the first argument
	switch ta.Kind() {
	case reflect.String:
		return a.(string) < b.(string)
	case reflect.Float64:
		return a.(float64) < b.(float64)
	}

	return false
}

// lessThanEqual will return true if a <= b
func lessThanEqual(a, b interface{}) bool {
	// If the values are equal, no more work necessary
	if a == b {
		return true
	}

	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// Make sure the types are the same
	if ta != tb {
		return false
	}

	// We have already checked that each argument is the same type
	// so it is safe to only check the first argument
	switch ta.Kind() {
	case reflect.String:
		return a.(string) <= b.(string)
	case reflect.Float64:
		return a.(float64) <= b.(float64)
	}

	return false
}

// greaterThan will return true if a > b
func greaterThan(a, b interface{}) bool {
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// Make sure the types are the same
	if ta != tb {
		return false
	}

	// We have already checked that each argument is the same type
	// so it is safe to only check the first argument
	switch ta.Kind() {
	case reflect.String:
		return a.(string) > b.(string)
	case reflect.Float64:
		return a.(float64) > b.(float64)
	}

	return false
}

// greaterThanEqual will return true if a >= b
func greaterThanEqual(a, b interface{}) bool {
	// If the values are equal, no more work necessary
	if a == b {
		return true
	}

	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)

	// Make sure the types are the same
	if ta != tb {
		return false
	}

	// We have already checked that each argument is the same type
	// so it is safe to only check the first argument
	switch ta.Kind() {
	case reflect.String:
		return a.(string) >= b.(string)
	case reflect.Float64:
		return a.(float64) >= b.(float64)
	}

	return false
}
