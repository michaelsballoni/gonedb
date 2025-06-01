package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestProps(t *testing.T) {
	db := GetTestDb("TestProps.db")
	defer db.Close()

	// create the test nodes
	var from_node, to_node gonedb.Node
	var err error
	from_node, err = gonedb.Nodes.Create(db, 0, GetTestStringId(db, "from"), 0)
	AssertNoError(err)
	to_node, err = gonedb.Nodes.Create(db, 0, GetTestStringId(db, "to"), 0)
	AssertNoError(err)

	// set a foo / bar property and get it back
	set_err := gonedb.Props.Set(db, gonedb.NodeItemTypeId, from_node.Id, GetTestStringId(db, "foo"), GetTestStringId(db, "bar"))
	AssertNoError(set_err)
	from_prop_value_string_id, from_err := gonedb.Props.Get(db, gonedb.NodeItemTypeId, from_node.Id, GetTestStringId(db, "foo"))
	AssertNoError(from_err)
	AssertEqual(GetTestStringId(db, "bar"), from_prop_value_string_id)

	// set a new baz value on the foo property
	set_err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, from_node.Id, GetTestStringId(db, "foo"), GetTestStringId(db, "baz"))
	AssertNoError(set_err)
	from_prop_value_string_id, from_err = gonedb.Props.Get(db, gonedb.NodeItemTypeId, from_node.Id, GetTestStringId(db, "foo"))
	AssertNoError(from_err)
	AssertEqual(GetTestStringId(db, "baz"), from_prop_value_string_id)

	// set blet / summary / abbra properties on the to node
	gonedb.Props.Set(db, gonedb.NodeItemTypeId, to_node.Id, GetTestStringId(db, "blet"), GetTestStringId(db, "monkey"))
	gonedb.Props.Set(db, gonedb.NodeItemTypeId, to_node.Id, GetTestStringId(db, "something"), GetTestStringId(db, "else"))
	gonedb.Props.Set(db, gonedb.NodeItemTypeId, to_node.Id, GetTestStringId(db, "abbra"), GetTestStringId(db, "caddabra"))
	to_props, to_err := gonedb.Props.GetAll(db, gonedb.NodeItemTypeId, to_node.Id)
	AssertNoError(to_err)
	AssertEqual(3, len(to_props))
	to_summary, sum_err := gonedb.Props.Summarize(db, to_props)
	AssertNoError(sum_err)
	AssertEqual("abbra caddabra\nblet monkey\nsomething else", to_summary)

	// erase the blet property
	gonedb.Props.Set(db, gonedb.NodeItemTypeId, to_node.Id, GetTestStringId(db, "blet"), -1)
	_, from_err = gonedb.Props.Get(db, gonedb.NodeItemTypeId, from_node.Id, GetTestStringId(db, "foo"))
	AssertNoError(from_err)
	to_props, to_err = gonedb.Props.GetAll(db, gonedb.NodeItemTypeId, to_node.Id)
	AssertNoError(to_err)
	AssertEqual(2, len(to_props))
	to_summary, sum_err = gonedb.Props.Summarize(db, to_props)
	AssertNoError(sum_err)
	AssertEqual("abbra caddabra\nsomething else", to_summary)

	// erase the rest of the properties
	gonedb.Props.Set(db, gonedb.NodeItemTypeId, to_node.Id, -1, -1)
	to_props, to_err = gonedb.Props.GetAll(db, gonedb.NodeItemTypeId, to_node.Id)
	AssertNoError(to_err)
	AssertEqual(0, len(to_props))
}
