package test

import (
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestSearch(t *testing.T) {
	db := GetTestDb("TestSearch.db")
	defer db.Close()

	var err error
	var node0, node1 gonedb.Node
	node0, err = gonedb.Nodes.Get(db, 0)
	AssertNoError(err)
	item_id0 := node0.Id

	node1, err = gonedb.Nodes.Create(db, item_id0, GetTestStringId(db, "show"), GetTestStringId(db, "type1"))
	AssertNoError(err)
	item_id1 := node1.Id

	node_results := []gonedb.Node{}

	//
	// NODES
	//
	{
		{
			var search_query gonedb.SearchQuery
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, item_id0, GetTestStringId(db, "foo"), GetTestStringId(db, "bar"))
		AssertNoError(err)

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "foo"), ValueString: "not it"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id0, node_results[0].Id)
		}

		err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, item_id0, GetTestStringId(db, "blet"), GetTestStringId(db, "monkey"))
		AssertNoError(err)

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria =
				[]gonedb.SearchCriteria{
					{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
					{NameStringId: GetTestStringId(db, "foo"), ValueString: "not it"},
				}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria =
				[]gonedb.SearchCriteria{
					{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
					{NameStringId: GetTestStringId(db, "blet"), ValueString: "monkey"},
				}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id0, node_results[0].Id)
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria =
				[]gonedb.SearchCriteria{
					{NameStringId: GetTestStringId(db, "foo"), ValueString: "bar"},
					{NameStringId: GetTestStringId(db, "blet"), ValueString: "monk%", UseLike: true},
				}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id0, node_results[0].Id)
		}

		err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, item_id1, GetTestStringId(db, "flint"), GetTestStringId(db, "stone"))
		AssertNoError(err)

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "flint"), ValueString: "not it"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "flint"), ValueString: "stone"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
		}
	}

	//
	// LINKS
	//
	var link1 gonedb.Link
	link1, err = gonedb.Links.Create(db, item_id0, item_id1, 0)
	AssertNoError(err)
	item_id2 := link1.Id

	{
		err = gonedb.Props.Set(db, gonedb.LinkItemTypeId, item_id2, GetTestStringId(db, "link"), GetTestStringId(db, "sink"))
		AssertNoError(err)

		link_results := []gonedb.Link{}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "link"), ValueString: "not it"}}
			link_results, err = gonedb.Search.FindLinks(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(link_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "link"), ValueString: "sink"}}
			link_results, err = gonedb.Search.FindLinks(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(link_results))
			AssertEqual(item_id2, link_results[0].Id)
		}
	}

	//
	// ORDER BY / LIMIT
	//
	{
		err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, item_id0, GetTestStringId(db, "some"), GetTestStringId(db, "one"))
		AssertNoError(err)
		err = gonedb.Props.Set(db, gonedb.NodeItemTypeId, item_id1, GetTestStringId(db, "some"), GetTestStringId(db, "two"))
		AssertNoError(err)

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria =
				[]gonedb.SearchCriteria{
					{NameStringId: GetTestStringId(db, "some"), ValueString: "%", UseLike: true},
				}
			search_query.OrderBy = "some"
			search_query.OrderAscensing = true
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(2, len(node_results))
			AssertEqual(item_id0, node_results[0].Id)
			AssertEqual(item_id1, node_results[1].Id)
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria =
				[]gonedb.SearchCriteria{
					{NameStringId: GetTestStringId(db, "some"), ValueString: "%", UseLike: true},
				}
			search_query.OrderBy = "some"
			search_query.OrderAscensing = false
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(2, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
			AssertEqual(item_id0, node_results[1].Id)
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
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
		}
	}

	//
	// SEARCH BY PAYLOAD
	//
	{
		err = gonedb.Nodes.SetPayload(db, item_id1, "some payload")
		AssertNoError(err)

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "payload"), ValueString: "not that payload"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "payload"), ValueString: "some payload"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
		}
	}

	//
	// SEARCH BY NAME
	//
	{
		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "name"), ValueString: "slow poke"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "name"), ValueString: "show"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
		}
	}

	//
	// SEARCH BY TYPE
	//
	{
		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "type"), ValueString: "not my type"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(0, len(node_results))
		}

		{
			var search_query gonedb.SearchQuery
			search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "type"), ValueString: "type1"}}
			node_results, err = gonedb.Search.FindNodes(db, &search_query)
			AssertNoError(err)
			AssertEqual(1, len(node_results))
			AssertEqual(item_id1, node_results[0].Id)
		}
	}

	//
	// SEARCH BY PARENT OR PATH
	//
	{
		{
			root_parent, err := gonedb.Nodes.Create(db, item_id0, GetTestStringId(db, "trunk"), GetTestStringId(db, "tree"))
			AssertNoError(err)
			leaf_node, err := gonedb.Nodes.Create(db, root_parent.Id, GetTestStringId(db, "leaf"), GetTestStringId(db, "plant"))
			AssertNoError(err)
			leafy_node, err := gonedb.Nodes.Create(db, leaf_node.Id, GetTestStringId(db, "leafy"), GetTestStringId(db, "plant"))
			AssertNoError(err)
			leafier_node, err := gonedb.Nodes.Create(db, leaf_node.Id, GetTestStringId(db, "leafier"), GetTestStringId(db, "plant"))
			AssertNoError(err)

			{
				var search_query gonedb.SearchQuery
				search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "parent"), ValueString: "/not the root"}}
				node_results, err = gonedb.Search.FindNodes(db, &search_query)
				AssertNoError(err)
				AssertEqual(0, len(node_results))
			}

			{
				var search_query gonedb.SearchQuery
				search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "parent"), ValueString: "/trunk"}}
				node_results, err = gonedb.Search.FindNodes(db, &search_query)
				AssertNoError(err)
				AssertEqual(1, len(node_results))
				AssertEqual(leaf_node.Id, node_results[0].Id)
			}

			{
				var search_query gonedb.SearchQuery
				search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/sprunklediodl7y/leaf"}}
				node_results, err = gonedb.Search.FindNodes(db, &search_query)
				AssertNoError(err)
				AssertEqual(0, len(node_results))
			}

			{
				var search_query gonedb.SearchQuery
				search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/trunk/leaf"}}
				node_results, err = gonedb.Search.FindNodes(db, &search_query)
				AssertNoError(err)
				AssertEqual(2, len(node_results))
				AssertTrue(
					(node_results[0].Id == leafy_node.Id && node_results[1].Id == leafier_node.Id) ||
						(node_results[1].Id == leafy_node.Id && node_results[0].Id == leafier_node.Id))
			}

			{
				var search_query gonedb.SearchQuery
				search_query.Criteria = []gonedb.SearchCriteria{{NameStringId: GetTestStringId(db, "path"), ValueString: "/trunk"}}
				node_results, err = gonedb.Search.FindNodes(db, &search_query)
				AssertNoError(err)
				AssertEqual(3, len(node_results))
				ids := map[int64]bool{}
				for _, n := range node_results {
					ids[n.Id] = true
				}
				AssertTrue(ids[leaf_node.Id])
				AssertTrue(ids[leafy_node.Id])
				AssertTrue(ids[leafier_node.Id])
			}
		}
	}
}
