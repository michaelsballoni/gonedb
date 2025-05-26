package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestNodePaths(t *testing.T) {
	db := GetTestDb("TestNodePaths.db")
	defer db.Close()

	// set up root, kid, and grandkids
	root_node, root_node_err := gonedb.Nodes.Create(db, 0, GetTestStringId(db, "foobar"), GetTestStringId(db, "root"))
	AssertNoError(root_node_err)

	child_node, child_node_err := gonedb.Nodes.Create(db, root_node.Id, GetTestStringId(db, "bletmonkey"), GetTestStringId(db, "child"))
	AssertNoError(child_node_err)

	grandchild_node1, grandchild_node1_err := gonedb.Nodes.Create(db, child_node.Id, GetTestStringId(db, "funkadelic"), GetTestStringId(db, "grandkid"))
	AssertNoError(grandchild_node1_err)

	grandchild_node2, grandchild_node2_err := gonedb.Nodes.Create(db, child_node.Id, GetTestStringId(db, "superfly"), GetTestStringId(db, "grandkid"))
	AssertNoError(grandchild_node2_err)

	gc1, err1 := gonedb.NodePaths.GetNodes(db, grandchild_node1)
	AssertNoError(err1)
	AssertEqual(3, len(gc1))

	gc2, err2 := gonedb.NodePaths.GetNodes(db, grandchild_node2)
	AssertNoError(err2)
	AssertEqual(3, len(gc2))
}
