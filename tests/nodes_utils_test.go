package test

import (
	"slices"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestIdsToSqlIn(t *testing.T) {
	_, e := gonedb.NodeUtils.IdsToSqlIn([]int64{})
	if e == nil {
		t.Error("IdsToSqlIn(Empty call not errored)")
	}

	testCases := []struct {
		name     string
		input    []int64
		expected string
	}{
		{"One ID", []int64{0}, "0"},
		{"Two IDs", []int64{0, 1}, "0,1"},
		{"Three IDs", []int64{0, 1, 2}, "0,1,2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			s, e := gonedb.NodeUtils.IdsToSqlIn(tc.input)
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
	_, e := gonedb.NodeUtils.StringToIds("0,1,2,foo,bar")
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
		{"Two IDs", "0/1", []int64{0, 1}},
		{"Three IDs", "0/1/2", []int64{0, 1, 2}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, e := gonedb.NodeUtils.StringToIds(tc.input)
			if e != nil {
				t.Errorf("Call errored: %v", e)
			}
			if !slices.Equal(v, tc.expected) {
				t.Errorf("Called failed: is %v - expected %v", v, tc.expected)
			}
		})
	}
}

func TestIdsToParentsStr(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int64
		expected string
	}{
		{"No IDs", []int64{}, ""},
		{"One ID (null)", []int64{0}, ""},
		{"One ID (one)", []int64{1}, "1/"},
		{"Two IDs", []int64{0, 1}, "1/"},
		{"Three IDs", []int64{0, 1, 2}, "1/2/"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := gonedb.NodeUtils.IdsToParentsStr(tc.input)
			if output != tc.expected {
				t.Errorf("Called failed: is %v - expected %v", output, tc.expected)
			}
		})
	}
}
