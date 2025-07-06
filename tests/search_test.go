package test

import (
	"database/sql"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestBasicSearch(t *testing.T) {
	db := GetTestDb("TestBasicSearch.db")
	defer db.Close()
	ctxt := createCtxt(db)

	{
		var search_query gonedb.SearchQuery
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	ctxt.err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, ctxt.item_id0, GetTestStringId(db, "foo"), GetTestStringId(db, "bar"))
	AssertNoError(ctxt.err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "foo"), ValueString: "not it"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node0, ctxt.node_results[0])
	}

	ctxt.err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, ctxt.item_id0, GetTestStringId(db, "blet"), GetTestStringId(db, "monkey"))
	AssertNoError(ctxt.err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
				{NameStringId: GetTestStringId(db, "foo"), ValueString: "not it"},
			}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
				{NameStringId: GetTestStringId(db, "blet"), ValueString: "monkey"},
			}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node0, ctxt.node_results[0])
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
				{NameStringId: GetTestStringId(db, "blet"), ValueString: "monk%", UseLike: true},
			}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node0, ctxt.node_results[0])
	}

	ctxt.err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, ctxt.item_id1, GetTestStringId(db, "flint"), GetTestStringId(db, "stone"))
	AssertNoError(ctxt.err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "flint"), ValueString: "not it"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "flint"), ValueString: "stone"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
	}
}

func TestSearchLinks(t *testing.T) {
	db := GetTestDb("TestSearchLinks.db")
	defer db.Close()
	ctxt := createCtxt(db)

	var link1 gonedb.Link
	link1, ctxt.err = gonedb.Links.Create(db, ctxt.item_id0, ctxt.item_id1, 0)
	AssertNoError(ctxt.err)
	item_id2 := link1.Id

	ctxt.err = gonedb.Props.Set(db, gonedb.LinkItemTypeId, item_id2, GetTestStringId(db, "link"), GetTestStringId(db, "sink"))
	AssertNoError(ctxt.err)

	var link_results []gonedb.Link

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "link"), ValueString: "not it"}}
		link_results, ctxt.err = gonedb.Search.FindLinks(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(link_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "link"), ValueString: "sink"}}
		link_results, ctxt.err = gonedb.Search.FindLinks(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(link_results))
		AssertEqual(link1, link_results[0])
	}
}

func TestSearchOrderLimits(t *testing.T) {
	db := GetTestDb("TestSearchOrderLimits.db")
	defer db.Close()
	ctxt := createCtxt(db)

	ctxt.err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, ctxt.item_id0, GetTestStringId(db, "some"), GetTestStringId(db, "one"))
	AssertNoError(ctxt.err)
	ctxt.err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, ctxt.item_id1, GetTestStringId(db, "some"), GetTestStringId(db, "two"))
	AssertNoError(ctxt.err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "some"), ValueString: "%", UseLike: true},
			}
		search_query.OrderBy = "some"
		search_query.OrderAscensing = true
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(2, len(ctxt.node_results))
		AssertEqual(ctxt.node0, ctxt.node_results[0])
		AssertEqual(ctxt.node1, ctxt.node_results[1])
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "some"), ValueString: "%", UseLike: true},
			}
		search_query.OrderBy = "some"
		search_query.OrderAscensing = false
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(2, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
		AssertEqual(ctxt.node0, ctxt.node_results[1])
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria =
			[]gonedb.SearchCriteria{
				{NameStringId: GetTestStringId(db, "some"), ValueString: "%", UseLike: true},
			}
		search_query.OrderBy = "some"
		search_query.OrderAscensing = false
		search_query.Limit = 1
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
	}
}

func TestSearchPayload(t *testing.T) {
	db := GetTestDb("TestSearchPayload.db")
	defer db.Close()
	ctxt := createCtxt(db)

	ctxt.err = gonedb.Nodes.SetPayload(db, ctxt.item_id1, "some payload")
	AssertNoError(ctxt.err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "payload"), ValueString: "not that payload"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "payload"), ValueString: "some payload"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
	}
}

// SEARCH BY NAME
func TestSearchByName(t *testing.T) {
	db := GetTestDb("TestSearchName.db")
	defer db.Close()
	ctxt := createCtxt(db)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "name"), ValueString: "slow poke"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "name"), ValueString: "show"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
	}
}

func TestSearchByType(t *testing.T) {
	db := GetTestDb("TestSearchType.db")
	defer db.Close()
	ctxt := createCtxt(db)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "type"), ValueString: "not my type"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "type"), ValueString: "type1"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(ctxt.node1, ctxt.node_results[0])
	}
}

func TestSearchByParent(t *testing.T) {
	db := GetTestDb("TestSearchParent.db")
	defer db.Close()
	ctxt := createCtxt(db)

	root_parent, err := gonedb.Nodes.Create(db, ctxt.item_id0, GetTestStringId(db, "trunk"), GetTestStringId(db, "tree"))
	AssertNoError(err)
	leaf_node, err := gonedb.Nodes.Create(db, root_parent.Id, GetTestStringId(db, "leaf"), GetTestStringId(db, "plant"))
	AssertNoError(err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "parent"), ValueString: "/not the root"}}
		ctxt.node_results, err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "parent"), ValueString: "/trunk"}}
		ctxt.node_results, err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(err)
		AssertEqual(1, len(ctxt.node_results))
		AssertEqual(leaf_node, ctxt.node_results[0])
	}
}

func TestSearchByPath(t *testing.T) {
	db := GetTestDb("TestSearchPath.db")
	defer db.Close()
	ctxt := createCtxt(db)

	root_parent, err := gonedb.Nodes.Create(db, ctxt.item_id0, GetTestStringId(db, "trunk"), GetTestStringId(db, "tree"))
	AssertNoError(err)
	leaf_node, err := gonedb.Nodes.Create(db, root_parent.Id, GetTestStringId(db, "leaf"), GetTestStringId(db, "plant"))
	AssertNoError(err)
	leafy_node, err := gonedb.Nodes.Create(db, leaf_node.Id, GetTestStringId(db, "leafy"), GetTestStringId(db, "plant"))
	AssertNoError(err)
	leafier_node, err := gonedb.Nodes.Create(db, leaf_node.Id, GetTestStringId(db, "leafier"), GetTestStringId(db, "plant"))
	AssertNoError(err)

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/sprunklediodl7y/leaf"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		AssertEqual(0, len(ctxt.node_results))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/trunk/leaf"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		// DEBUG
		//fmt.Printf(">>>>>> /trunk/leaf node_results: %v <<<<<<<\n", ctxt.node_results)
		AssertEqual(2, len(ctxt.node_results))
		AssertTrue(
			(ctxt.node_results[0] == leafy_node && ctxt.node_results[1] == leafier_node) ||
				(ctxt.node_results[1] == leafy_node && ctxt.node_results[0] == leafier_node))
	}

	{
		var search_query gonedb.SearchQuery
		search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/trunk"}}
		ctxt.node_results, ctxt.err = gonedb.Search.FindNodes(db, &search_query)
		AssertNoError(ctxt.err)
		// DEBUG
		//fmt.Printf(">>>>>> /trunk node_results: %v <<<<<<<\n", ctxt.node_results)
		AssertEqual(3, len(ctxt.node_results))
		ids := map[int64]bool{}
		for _, n := range ctxt.node_results {
			ids[n.Id] = true
		}
		AssertTrue(ids[leaf_node.Id])
		AssertTrue(ids[leafy_node.Id])
		AssertTrue(ids[leafier_node.Id])
	}
}

type searchTestCtxt struct {
	err          error
	node0        gonedb.Node
	node1        gonedb.Node
	item_id0     int64
	item_id1     int64
	node_results []gonedb.Node
}

func createCtxt(db *sql.DB) searchTestCtxt {
	var output searchTestCtxt

	output.node0, output.err = gonedb.Nodes.Get(db, 0)
	AssertNoError(output.err)
	item_id0 := output.node0.Id

	output.node1, output.err = gonedb.Nodes.Create(db, item_id0, GetTestStringId(db, "show"), GetTestStringId(db, "type1"))
	AssertNoError(output.err)
	output.item_id1 = output.node1.Id

	return output
}
