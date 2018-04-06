/*
This file was originally written for the grules library.
https://github.com/huttotw/grules

Copyright © 2018 Trevor Hutto
Licensed under the Apache License, Version 2.0 (the "License")
*/
package jules

import (
	"testing"
)

type testCase struct {
	args     []interface{}
	expected bool
}

func TestEqual(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: true},
		testCase{args: []interface{}{"a", "b"}, expected: false},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: true},
		testCase{args: []interface{}{float64(1), float64(0)}, expected: false},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(0.1)}, expected: false},
	}

	for i, c := range cases {
		res := equal(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		equal("benchmark", "benchmark")
	}
}

func TestNotEqual(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: false},
		testCase{args: []interface{}{"a", "b"}, expected: true},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: false},
		testCase{args: []interface{}{float64(1), float64(0)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: false},
		testCase{args: []interface{}{float64(1.1), float64(0.1)}, expected: true},
	}

	for i, c := range cases {
		res := notEqual(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkNotEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		notEqual("benchmark", "not-benchmark")
	}
}

func TestLessThan(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: false},
		testCase{args: []interface{}{"a", "b"}, expected: true},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: false},
		testCase{args: []interface{}{float64(0), float64(1)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: false},
		testCase{args: []interface{}{float64(1.1), float64(1.2)}, expected: true},
	}

	for i, c := range cases {
		res := lessThan(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkLessThan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lessThan(0, 1)
	}
}

func TestLessThanEqual(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: true},
		testCase{args: []interface{}{"a", "b"}, expected: true},
		testCase{args: []interface{}{"c", "b"}, expected: false},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: true},
		testCase{args: []interface{}{float64(0), float64(1)}, expected: true},
		testCase{args: []interface{}{float64(1), float64(0)}, expected: false},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.2)}, expected: true},
		testCase{args: []interface{}{float64(1.2), float64(1.1)}, expected: false},
	}

	for i, c := range cases {
		res := lessThanEqual(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkLessThanEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lessThanEqual(0, 0)
	}
}

func TestGreaterThan(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: false},
		testCase{args: []interface{}{"b", "a"}, expected: true},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: false},
		testCase{args: []interface{}{float64(1), float64(0)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: false},
		testCase{args: []interface{}{float64(1.2), float64(1.1)}, expected: true},
	}

	for i, c := range cases {
		res := greaterThan(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkGreaterThan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		greaterThan(1, 0)
	}
}

func TestGreaterThanEqual(t *testing.T) {
	cases := []testCase{
		testCase{args: []interface{}{"a", "a"}, expected: true},
		testCase{args: []interface{}{"a", "b"}, expected: false},
		testCase{args: []interface{}{"c", "b"}, expected: true},
		testCase{args: []interface{}{float64(1), float64(1)}, expected: true},
		testCase{args: []interface{}{float64(0), float64(1)}, expected: false},
		testCase{args: []interface{}{float64(1), float64(0)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.1)}, expected: true},
		testCase{args: []interface{}{float64(1.1), float64(1.2)}, expected: false},
		testCase{args: []interface{}{float64(1.2), float64(1.1)}, expected: true},
	}

	for i, c := range cases {
		res := greaterThanEqual(c.args[0], c.args[1])
		if res != c.expected {
			t.Fatalf("expected case %d to be %v, got %v", i, c.expected, res)
		}
	}
}

func BenchmarkGreaterThanEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		greaterThanEqual(0, 0)
	}
}
