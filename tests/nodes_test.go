package test

import (
	"slices"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestIdsToSqlIn(t *testing.T) {
	_, e := gonedb.Nodes.IdsToSqlIn([]int64{})
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

			s, e := gonedb.Nodes.IdsToSqlIn(tc.input)
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
	_, e := gonedb.Nodes.StringToIds("0,1,2,foo,bar")
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
		{"Two IDs", "0,1", []int64{0, 1}},
		{"Three IDs", "0,1,2", []int64{0, 1, 2}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, e := gonedb.Nodes.StringToIds(tc.input)
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
			output := gonedb.Nodes.IdsToParentsStr(tc.input)
			if output != tc.expected {
				t.Errorf("Called failed: is %v - expected %v", output, tc.expected)
			}
		})
	}
}

func TestGetParentsNodeIds(t *testing.T) {
	db := GetTestDb("TestGetParentsNodeIds.db")
	defer db.Close()

	nameStrId1, err1 := gonedb.Strings.GetId(db, "foobar")
	AssertNoError(t, err1)

	typeStrId1, err2 := gonedb.Strings.GetId(db, "my type")
	AssertNoError(t, err2)

	node1, err3 := gonedb.Nodes.Create(db, 0, nameStrId1, typeStrId1)
	AssertNoError(t, err3)
	AssertTrue(t, node1.Id > 0)
	AssertEqual(t, 0, node1.ParentId)
	AssertEqual(t, nameStrId1, node1.NameStringId)
	AssertEqual(t, typeStrId1, node1.TypeStringId)

	node1b, err4 := gonedb.Nodes.Get(db, node1.Id)
	AssertNoError(t, err4)
	AssertEqual(t, node1.Id, node1b.Id)
	AssertEqual(t, node1.ParentId, node1b.ParentId)
	AssertEqual(t, node1.NameStringId, node1b.NameStringId)
	AssertEqual(t, node1.TypeStringId, node1b.TypeStringId)

	nameStrId2, err5 := gonedb.Strings.GetId(db, "bletmonkey")
	AssertNoError(t, err5)

	typeStrId2, err6 := gonedb.Strings.GetId(db, "my other type")
	AssertNoError(t, err6)

	node2, err7 := gonedb.Nodes.Create(db, node1.Id, nameStrId2, typeStrId2)
	AssertNoError(t, err7)
	AssertTrue(t, node2.Id > 0)
	AssertEqual(t, node1.Id, node2.ParentId)
	AssertEqual(t, nameStrId2, node2.NameStringId)
	AssertEqual(t, typeStrId2, node2.TypeStringId)

	node2b, err8 := gonedb.Nodes.Get(db, node2.Id)
	AssertNoError(t, err8)
	AssertEqual(t, node2.Id, node2b.Id)
	AssertEqual(t, node2.ParentId, node2b.ParentId)
	AssertEqual(t, node2.NameStringId, node2b.NameStringId)
	AssertEqual(t, node2.TypeStringId, node2b.TypeStringId)

	/* FORNOW
	null_parent_nodes, err9a := gonedb.Nodes.GetParentsNodeIds(db, 0)
	AssertNoError(t, err9a)
	AssertEqual(t, 0, len(null_parent_nodes))

	parent_nodes, err9 := gonedb.Nodes.GetParentsNodeIds(db, node2.Id)
	AssertNoError(t, err9)
	AssertEqual(t, 1, len(parent_nodes))
	AssertEqual(t, node2.ParentId, parent_nodes[0])
	*/
}
