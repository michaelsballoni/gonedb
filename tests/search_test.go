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

	node1, err = gonedb.Nodes.Create(db, item_id0, GetTestStringId(db, "show"), 0)
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

	/* FORNOW
	//
	// SEARCH BY PAYLOAD
	//
	{
		nodes::set_payload(db, item_id1, L"some payload");
		search_query search1({ search_criteria(strings::get_id(db, L"payload"), L"not that payload") });
		auto no_payload_results = search::find_nodes(db, search1);
		Assert::IsTrue(no_payload_results.empty());

		search_query search2({ search_criteria(strings::get_id(db, L"payload"), L"some payload") });
		auto with_payload_results = search::find_nodes(db, search2);
		Assert::AreEqual(size_t(1), with_payload_results.size());
		Assert::IsTrue(!with_payload_results[0].payload.has_value());

		search_query search3({ search_criteria(strings::get_id(db, L"payload"), L"some payload") });
		search3.m_includePayload = true;
		auto with_payload_results2 = search::find_nodes(db, search3);
		Assert::AreEqual(size_t(1), with_payload_results2.size());
		Assert::AreEqual(std::wstring(L"some payload"), with_payload_results2[0].payload.value());
	}

	//
	// SEARCH BY NAME
	//
	{
		search_query search1({ search_criteria(strings::get_id(db, L"name"), L"slow poke") });
		auto no_results = search::find_nodes(db, search1);
		Assert::IsTrue(no_results.empty());

		search_query search2({ search_criteria(strings::get_id(db, L"name"), L"show") });
		auto with_results = search::find_nodes(db, search2);
		Assert::AreEqual(size_t(1), with_results.size());
		Assert::IsTrue(with_results[0] == node1);
	}

	//
	// SEARCH BY PATH
	//
	auto node2 = nodes::create(db, node1.id, strings::get_id(db, L"leafy"), 0);
	{
		search_query search1({ search_criteria(strings::get_id(db, L"path"), L"/fred/nothing/ha ha") });
		auto no_results = search::find_nodes(db, search1);
		Assert::IsTrue(no_results.empty());

		search_query search2({ search_criteria(strings::get_id(db, L"path"), L"/show") });
		auto with_results = search::find_nodes(db, search2);
		Assert::AreEqual(size_t(1), with_results.size());
		Assert::IsTrue(with_results[0] == node2);
	}

	//
	// SEARCH BY PARENT
	//
	{
		auto node3 = nodes::create(db, node1.id, strings::get_id(db, L"leaf"), 0);
		auto node4 = nodes::create(db, node3.id, strings::get_id(db, L"leafier"), 0);

		search_query search2({ search_criteria(strings::get_id(db, L"path"), L"/show") });
		auto with_results = search::find_nodes(db, search2);
		Assert::AreEqual(size_t(3), with_results.size());
		Assert::IsTrue(hasNode(with_results, node2.Id));
		Assert::IsTrue(hasNode(with_results, node3.Id));
		Assert::IsTrue(hasNode(with_results, node4.Id));

		search_query search1({ search_criteria(strings::get_id(db, L"parent"), L"/fred/nothing/ha ha") });
		auto no_results = search::find_nodes(db, search1);
		Assert::IsTrue(no_results.empty());

		search_query search3({ search_criteria(strings::get_id(db, L"parent"), L"/show") });
		auto with_results2 = search::find_nodes(db, search3);
		Assert::AreEqual(size_t(2), with_results2.size());
		Assert::IsTrue(hasNode(with_results2, node2.Id));
		Assert::IsTrue(hasNode(with_results2, node3.Id));
	}

	//
	// SEARCH BY TYPE
	//
	{
		auto node3 = nodes::create(db, node1.id, strings::get_id(db, L"leaf2"), strings::get_id(db, L"type1"));
		auto node4 = nodes::create(db, node3.id, strings::get_id(db, L"leafier2"), strings::get_id(db, L"type2"));

		search_query search2({ search_criteria(strings::get_id(db, L"type"), L"type1") });
		auto with_results = search::find_nodes(db, search2);
		Assert::AreEqual(size_t(1), with_results.size());
		Assert::IsTrue(hasNode(with_results, node3.Id));

		search_query search3({ search_criteria(strings::get_id(db, L"type"), L"type2") });
		auto with_results2 = search::find_nodes(db, search3);
		Assert::AreEqual(size_t(1), with_results2.size());
		Assert::IsTrue(hasNode(with_results2, node4.Id));

		search_query search4({ search_criteria(strings::get_id(db, L"type"), L"fred") });
		auto with_results3 = search::find_nodes(db, search4);
		Assert::AreEqual(size_t(0), with_results3.size());
	}
	*/
}
