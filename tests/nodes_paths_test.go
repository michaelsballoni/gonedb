package test

import (
	"strings"
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
	AssertEqual(root_node, gc1[0])
	AssertEqual(child_node, gc1[1])
	AssertEqual(grandchild_node1, gc1[2])

	gc2, err2 := gonedb.NodePaths.GetNodes(db, grandchild_node2)
	AssertNoError(err2)
	AssertEqual(3, len(gc2))
	AssertEqual(root_node, gc2[0])
	AssertEqual(child_node, gc2[1])
	AssertEqual(grandchild_node2, gc2[2])

	str1, strErr1 := gonedb.NodePaths.GetStrs(db, grandchild_node1)
	AssertNoError(strErr1)
	AssertEqual("foobar/bletmonkey/funkadelic", strings.Join(str1, "/"))

	arr1star, arrErr := gonedb.NodePaths.GetStrNodes(db, strings.Split("foobar/bletmonkey/funkadelic", "/"))
	AssertNoError(arrErr)
	AssertTrue(arr1star != nil)
	arr1 := *arr1star
	AssertEqual(3, len(arr1))
	AssertEqual(root_node, arr1[0])
	AssertEqual(child_node, arr1[1])
	AssertEqual(grandchild_node1, arr1[2])
}
