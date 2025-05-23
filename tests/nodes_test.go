package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestGetParentsNodeIds(t *testing.T) {
	db := GetTestDb("TestGetParentsNodeIds.db")
	defer db.Close()

	nameStrId1 := GetTestStringId(db, "foobar")
	typeStrId1 := GetTestStringId(db, "my type")

	node1, node1_err := gonedb.Nodes.Create(db, 0, nameStrId1, typeStrId1)
	AssertNoError(node1_err)
	AssertTrue(node1.Id > 0)
	AssertEqual(0, node1.ParentId)
	AssertEqual(nameStrId1, node1.NameStringId)
	AssertEqual(typeStrId1, node1.TypeStringId)

	node1b, node1b_err := gonedb.Nodes.Get(db, node1.Id)
	AssertNoError(node1b_err)
	AssertEqual(node1.Id, node1b.Id)
	AssertEqual(node1.ParentId, node1b.ParentId)
	AssertEqual(node1.NameStringId, node1b.NameStringId)
	AssertEqual(node1.TypeStringId, node1b.TypeStringId)

	nameStrId2 := GetTestStringId(db, "bletmonkey")
	typeStrId2 := GetTestStringId(db, "my other type")

	node2, node2_err := gonedb.Nodes.Create(db, node1.Id, nameStrId2, typeStrId2)
	AssertNoError(node2_err)
	AssertTrue(node2.Id > 0)
	AssertEqual(node1.Id, node2.ParentId)
	AssertEqual(nameStrId2, node2.NameStringId)
	AssertEqual(typeStrId2, node2.TypeStringId)

	node2b, node2b_err := gonedb.Nodes.Get(db, node2.Id)
	AssertNoError(node2b_err)
	AssertEqual(node2.Id, node2b.Id)
	AssertEqual(node2.ParentId, node2b.ParentId)
	AssertEqual(node2.NameStringId, node2b.NameStringId)
	AssertEqual(node2.TypeStringId, node2b.TypeStringId)

	null_parent_nodes, null_parent_nodes_err := gonedb.Nodes.GetParentsNodeIds(db, 0)
	AssertNoError(null_parent_nodes_err)
	AssertEqual(0, len(null_parent_nodes))

	root_parent_nodes, root_parent_nodes_err := gonedb.Nodes.GetParentsNodeIds(db, node1.Id)
	AssertNoError(root_parent_nodes_err)
	AssertEqual(0, len(root_parent_nodes))

	child_parent_nodes, err11 := gonedb.Nodes.GetParentsNodeIds(db, node2.Id)
	AssertNoError(err11)
	AssertEqual(1, len(child_parent_nodes))
	AssertEqual(node2.ParentId, child_parent_nodes[0])
}

func TestNodeCopy(t *testing.T) {
	db := GetTestDb("TestNodeCopy.db")
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

	// create a new root to copy the child into
	new_root_node, new_root_node_err := gonedb.Nodes.Create(db, 0, GetTestStringId(db, "newbar"), GetTestStringId(db, "newroot"))
	AssertNoError(new_root_node_err)

	// copy the child into the new root
	copy_node_id, copy_err := gonedb.Nodes.Copy(db, child_node.Id, new_root_node.Id)
	AssertNoError(copy_err)

	// get the copy
	got_copy_node, got_copy_err := gonedb.Nodes.Get(db, copy_node_id)
	AssertNoError(got_copy_err)
	AssertEqual(copy_node_id, got_copy_node.Id)
	AssertEqual(new_root_node.Id, got_copy_node.ParentId) // null
	AssertEqual(child_node.NameStringId, got_copy_node.NameStringId)
	AssertEqual(child_node.TypeStringId, got_copy_node.TypeStringId)

	// get the new grandchildren
	copy_grandchildren, copy_grandchildren_err := gonedb.Nodes.GetChildren(db, copy_node_id)
	AssertNoError(copy_grandchildren_err)
	AssertEqual(2, len(copy_grandchildren))

	// compare new and old grandchildren
	AssertEqual(copy_node_id, copy_grandchildren[0].ParentId)
	AssertEqual(grandchild_node1.NameStringId, copy_grandchildren[0].NameStringId)
	AssertEqual(grandchild_node1.TypeStringId, copy_grandchildren[0].TypeStringId)

	AssertEqual(copy_node_id, copy_grandchildren[1].ParentId)
	AssertEqual(grandchild_node2.NameStringId, copy_grandchildren[1].NameStringId)
	AssertEqual(grandchild_node2.TypeStringId, copy_grandchildren[1].TypeStringId)
}

func TestNodeMove(t *testing.T) {
	db := GetTestDb("TestNodeMove.db")
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

	// create a new root to move the child into
	new_root_node, new_root_node_err := gonedb.Nodes.Create(db, 0, GetTestStringId(db, "newbar"), GetTestStringId(db, "newroot"))
	AssertNoError(new_root_node_err)

	// move the child into the new root
	move_err := gonedb.Nodes.Move(db, child_node.Id, new_root_node.Id)
	AssertNoError(move_err)

	// get the copy
	got_child_node, got_move_err := gonedb.Nodes.Get(db, child_node.Id)
	AssertNoError(got_move_err)
	AssertEqual(child_node.Id, got_child_node.Id)
	AssertEqual(new_root_node.Id, got_child_node.ParentId) // null
	AssertEqual(child_node.NameStringId, got_child_node.NameStringId)
	AssertEqual(child_node.TypeStringId, got_child_node.TypeStringId)

	// get the new grandchildren
	move_grandchildren, move_grandchildren_err := gonedb.Nodes.GetChildren(db, got_child_node.Id)
	AssertNoError(move_grandchildren_err)
	AssertEqual(2, len(move_grandchildren))

	// compare new and old grandchildren
	AssertEqual(got_child_node.Id, move_grandchildren[0].ParentId)
	AssertEqual(grandchild_node1.NameStringId, move_grandchildren[0].NameStringId)
	AssertEqual(grandchild_node1.TypeStringId, move_grandchildren[0].TypeStringId)

	AssertEqual(got_child_node.Id, move_grandchildren[1].ParentId)
	AssertEqual(grandchild_node2.NameStringId, move_grandchildren[1].NameStringId)
	AssertEqual(grandchild_node2.TypeStringId, move_grandchildren[1].TypeStringId)
}

func TestNodeRename(t *testing.T) {
	db := GetTestDb("TestNodeRename.db")
	defer db.Close()

	// set up root, kid, and grandkids
	root_node, root_node_err := gonedb.Nodes.Create(db, 0, GetTestStringId(db, "foobar"), GetTestStringId(db, "root"))
	AssertNoError(root_node_err)

	child_node, child_node_err := gonedb.Nodes.Create(db, root_node.Id, GetTestStringId(db, "bletmonkey"), GetTestStringId(db, "child"))
	AssertNoError(child_node_err)

	// rename the child in place
	ren_err := gonedb.Nodes.Rename(db, child_node.Id, GetTestStringId(db, "something else"))
	AssertNoError(ren_err)

	// get the renamed
	got_child_node, got_move_err := gonedb.Nodes.Get(db, child_node.Id)
	AssertNoError(got_move_err)
	AssertEqual(child_node.Id, got_child_node.Id)
	AssertEqual(child_node.ParentId, got_child_node.ParentId) // null
	AssertEqual(GetTestStringId(db, "something else"), got_child_node.NameStringId)
	AssertEqual(child_node.TypeStringId, got_child_node.TypeStringId)
}

func TestNodeRemove(t *testing.T) {
	db := GetTestDb("TestNodeRename.db")
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

	move_grandchildren, move_grandchildren_err := gonedb.Nodes.GetChildren(db, child_node.Id)
	AssertNoError(move_grandchildren_err)
	AssertEqual(2, len(move_grandchildren))

	// rename the child in place
	rem_err := gonedb.Nodes.Remove(db, child_node.Id)
	AssertNoError(rem_err)

	// get the renamed
	_, get_rem_err := gonedb.Nodes.Get(db, child_node.Id)
	AssertError(get_rem_err)

	// get the grandchildren
	grandkids, rem_grandchildren_err := gonedb.Nodes.GetChildren(db, child_node.Id)
	AssertNoError(rem_grandchildren_err)
	AssertEqual(0, len(grandkids))

	_, grandchild_node1_err_after := gonedb.Nodes.Get(db, grandchild_node1.Id)
	AssertError(grandchild_node1_err_after)

	_, grandchild_node2_err_after := gonedb.Nodes.Get(db, grandchild_node2.Id)
	AssertError(grandchild_node2_err_after)
}

func TestNodePayload(t *testing.T) {
	db := GetTestDb("TestNodePayload.db")
	defer db.Close()

	// set up root, kid, and grandkids
	root_node, root_node_err := gonedb.Nodes.Create(db, 0, GetTestStringId(db, "foobar"), GetTestStringId(db, "root"))
	AssertNoError(root_node_err)

	payload, err := gonedb.Nodes.GetPayload(db, root_node.Id)
	AssertNoError(err)
	AssertEqual("", payload)

	AssertNoError(gonedb.Nodes.SetPayload(db, root_node.Id, "foobar"))

	payload, err = gonedb.Nodes.GetPayload(db, root_node.Id)
	AssertNoError(err)
	AssertEqual("foobar", payload)

	AssertNoError(gonedb.Nodes.SetPayload(db, root_node.Id, "blet monkey"))

	payload, err = gonedb.Nodes.GetPayload(db, root_node.Id)
	AssertNoError(err)
	AssertEqual("blet monkey", payload)

	AssertNoError(gonedb.Nodes.Remove(db, root_node.Id))

	AssertError(gonedb.Nodes.SetPayload(db, root_node.Id, "post remove"))
}
