package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestLinks(t *testing.T) {
	db := GetTestDb("TestLinks.db")
	defer db.Close()

	var from_node, to_node gonedb.Node
	var err error
	from_node, err = gonedb.Nodes.Create(db, 0, GetTestStringId(db, "from"), 0)
	AssertNoError(err)
	to_node, err = gonedb.Nodes.Create(db, 0, GetTestStringId(db, "to"), 0)
	AssertNoError(err)
	type_string_id, err := gonedb.Strings.GetId(db, "my type")
	AssertNoError(err)

	link, err := gonedb.Links.Create(db, from_node.Id, to_node.Id, type_string_id)
	AssertNoError(err)
	AssertEqual(from_node.Id, link.FromNodeId)
	AssertEqual(to_node.Id, link.ToNodeId)
	AssertEqual(type_string_id, link.TypeStringId)

	var payload string
	payload, err = gonedb.Links.GetPayload(db, link.Id)
	AssertNoError(err)
	AssertEqual("", payload)

	err = gonedb.Links.SetPayload(db, link.Id, "foobar")
	AssertNoError(err)

	payload, err = gonedb.Links.GetPayload(db, link.Id)
	AssertNoError(err)
	AssertEqual("foobar", payload)

	var out_links []gonedb.Link
	out_links, err = gonedb.Links.GetOutLinks(db, from_node.Id)
	AssertNoError(err)
	AssertEqual(1, len(out_links))
	AssertEqual(from_node.Id, out_links[0].FromNodeId)
	AssertEqual(to_node.Id, out_links[0].ToNodeId)
	AssertEqual(type_string_id, out_links[0].TypeStringId)

	var to_links []gonedb.Link
	to_links, err = gonedb.Links.GetToLinks(db, to_node.Id)
	AssertNoError(err)
	AssertEqual(1, len(to_links))
	AssertEqual(from_node.Id, to_links[0].FromNodeId)
	AssertEqual(to_node.Id, to_links[0].ToNodeId)
	AssertEqual(type_string_id, to_links[0].TypeStringId)

	var link_gotten gonedb.Link
	link_gotten, err = gonedb.Links.Get(db, link.Id)
	AssertNoError(err)
	AssertEqual(link.FromNodeId, link_gotten.FromNodeId)
	AssertEqual(link.ToNodeId, link_gotten.ToNodeId)
	AssertEqual(link.TypeStringId, link_gotten.TypeStringId)

	err = gonedb.Links.Remove(db, link.Id)
	AssertNoError(err)

	link, err = gonedb.Links.Create(db, from_node.Id, to_node.Id, type_string_id)
	AssertNoError(err)
	AssertEqual(from_node.Id, link.FromNodeId)
	AssertEqual(to_node.Id, link.ToNodeId)
	AssertEqual(type_string_id, link.TypeStringId)

	err = gonedb.Links.RemoveFromTo(db, link.FromNodeId, link.ToNodeId, type_string_id)
	AssertNoError(err)

	_, err = gonedb.Links.Get(db, link.Id)
	AssertError(err)
}
