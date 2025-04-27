package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestGetParentsNodeIds(t *testing.T) {
	db := GetTestDb("TestGetParentsNodeIds.db")
	defer db.Close()

	nameStrId1 := GetStringId(t, db, "foobar")
	typeStrId1 := GetStringId(t, db, "my type")

	node1, node1_err := gonedb.Nodes.Create(db, 0, nameStrId1, typeStrId1)
	AssertNoError(t, node1_err)
	AssertTrue(t, node1.Id > 0)
	AssertEqual(t, 0, node1.ParentId)
	AssertEqual(t, nameStrId1, node1.NameStringId)
	AssertEqual(t, typeStrId1, node1.TypeStringId)

	node1b, node1b_err := gonedb.Nodes.Get(db, node1.Id)
	AssertNoError(t, node1b_err)
	AssertEqual(t, node1.Id, node1b.Id)
	AssertEqual(t, node1.ParentId, node1b.ParentId)
	AssertEqual(t, node1.NameStringId, node1b.NameStringId)
	AssertEqual(t, node1.TypeStringId, node1b.TypeStringId)

	nameStrId2 := GetStringId(t, db, "bletmonkey")
	typeStrId2 := GetStringId(t, db, "my other type")

	node2, node2_err := gonedb.Nodes.Create(db, node1.Id, nameStrId2, typeStrId2)
	AssertNoError(t, node2_err)
	AssertTrue(t, node2.Id > 0)
	AssertEqual(t, node1.Id, node2.ParentId)
	AssertEqual(t, nameStrId2, node2.NameStringId)
	AssertEqual(t, typeStrId2, node2.TypeStringId)

	node2b, node2b_err := gonedb.Nodes.Get(db, node2.Id)
	AssertNoError(t, node2b_err)
	AssertEqual(t, node2.Id, node2b.Id)
	AssertEqual(t, node2.ParentId, node2b.ParentId)
	AssertEqual(t, node2.NameStringId, node2b.NameStringId)
	AssertEqual(t, node2.TypeStringId, node2b.TypeStringId)

	null_parent_nodes, null_parent_nodes_err := gonedb.Nodes.GetParentsNodeIds(db, 0)
	AssertNoError(t, null_parent_nodes_err)
	AssertEqual(t, 0, len(null_parent_nodes))

	root_parent_nodes, root_parent_nodes_err := gonedb.Nodes.GetParentsNodeIds(db, node1.Id)
	AssertNoError(t, root_parent_nodes_err)
	AssertEqual(t, 1, len(root_parent_nodes))
	AssertEqual(t, node1.ParentId, root_parent_nodes[0])

	/*
		child_parent_nodes, err11 := gonedb.Nodes.GetParentsNodeIds(db, node2.Id)
		AssertNoError(t, err11)
		AssertEqual(t, 2, len(child_parent_nodes))
		AssertEqual(t, node2.ParentId, child_parent_nodes[0])
		AssertEqual(t, node2.Id, child_parent_nodes[1])
	*/
}

// FORNOW - Test...
// GetPayload / SetPayload
