package test

import (
	"fmt"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestIdsToSqlIn(t *testing.T) {
	_, e := gonedb.NodeUtils.IdsToSqlIn([]int64{})
	AssertTrue(e != nil)

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
		s, e := gonedb.NodeUtils.IdsToSqlIn(tc.input)
		AssertNoError(e)
		AssertEqual(tc.expected, s)
	}
}

func TestStringToIds(t *testing.T) {
	_, e := gonedb.NodeUtils.StringToIds("0,1,2,foo,bar")
	AssertTrue(e != nil)

	testCases := []struct {
		name     string
		input    string
		expected []int64
	}{
		{"No IDs", "", []int64{}},
		{"One ID", "0", []int64{0}},
		{"Two IDs", "0/1/", []int64{0, 1}},
		{"Three IDs", "0/1/2/", []int64{0, 1, 2}},
	}
	for _, tc := range testCases {
		s, e := gonedb.NodeUtils.StringToIds(tc.input)
		AssertNoError(e)
		AssertEqual(fmt.Sprintf("%v", tc.expected), fmt.Sprintf("%v", s))
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
		output := gonedb.NodeUtils.IdsToParentsStr(tc.input)
		AssertEqual(tc.expected, output)
	}
}

func TestGetLasttPathId(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"", 0},
		{"/", 0},
		{"/1", 1},
		{"/1/", 1},
		{"/1/2", 2},
		{"/1/2/", 2},
		{"/1/10/", 10},
		{"/20/10", 10},
		{"/20/10/", 10},
		{"/20/100", 100},
		{"/20/100/", 100},
	}
	for _, tc := range testCases {
		output, err := gonedb.NodeUtils.GetLasttPathId(tc.input)
		AssertNoError(err)
		AssertEqual(tc.expected, output)
	}
}
