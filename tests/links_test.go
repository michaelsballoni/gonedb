package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestLinks(t *testing.T) {
	db := GetTestDb("TestLinks.db")
	defer db.Close()

	from_node_id := int64(5)
	to_node_id := int64(9)
	type_string_id := int64(13)

	link, create_err := gonedb.Links.Create(db, from_node_id, to_node_id, type_string_id)
	AssertNoError(create_err)
	AssertEqual(from_node_id, link.FromNodeId)
	AssertEqual(to_node_id, link.ToNodeId)
	AssertEqual(type_string_id, link.TypeStringId)
}
