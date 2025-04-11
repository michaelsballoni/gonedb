package test

import (
	"slices"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestIdsToSqlIn(t *testing.T) {
	_, e := gonedb.IdsToSqlIn([]int64{})
	if e == nil {
		t.Error("IdsToSqlIn(Empty call not errored)")
	}

	testCases := []struct {
		name     string
		input    []int64
		expected string
	}{
		{"One ID", []int64{0}, "0"},
		{"Two ID", []int64{0, 1}, "0,1"},
		{"Three ID", []int64{0, 1, 2}, "0,1,2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			s, e := gonedb.IdsToSqlIn(tc.input)
			if e != nil {
				t.Errorf("Call errored")
			}
			if s != tc.expected {
				t.Errorf("Called failed: is %s - expected %s", s, tc.expected)
			}
		})
	}
}

func TestStringToIds(t *testing.T) {
	_, e := gonedb.StringToIds("0,1,2,foo,bar", ',')
	if e == nil {
		t.Error("StringToIds(Bad data not errored)")
	}

	testCases := []struct {
		name     string
		input    string
		expected []int64
	}{
		{"No IDs", "", []int64{}},
		{"One ID", "0", []int64{0}},
		{"Two ID", "0,1", []int64{0, 1}},
		{"Three ID", "0,1,2", []int64{0, 1, 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			v, e := gonedb.StringToIds(tc.input, ',')
			if e != nil {
				t.Errorf("Call errored")
			}
			if !slices.Equal(v, tc.expected) {
				t.Errorf("Called failed: is %v - expected %v", v, tc.expected)
			}
		})
	}
}
