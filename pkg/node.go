package gonedb

import (
	"errors"
	"strconv"
	"strings"
)

type Node struct {
	Id           int64
	ParentId     int64
	NameStringId int64
	TypeStringId int64
}

/*
func CreateNode(db *sql.DB, parentNodeId, nameStringId, typeStringId int64) Node {
	db.Exec(
}
*/

// Turn IDs into a string used in a SELECT IN (here)
// Errors if an empty input is provided
func IdsToSqlIn(ids []int64) (string, error) {
	if len(ids) == 0 {
		return "", errors.New("ids_to_sql_in: called with no IDs")
	}

	var output strings.Builder
	for _, id := range ids {
		if output.Len() > 0 {
			output.WriteRune(',')
		}
		output.WriteString(strconv.FormatInt(id, 10))
	}
	return output.String(), nil
}

// only IDs found between separators are returned
// nothing in, nothing out
func StringToIds(str string, separator rune) ([]int64, error) {
	var ids []int64
	var collector strings.Builder

	for _, c := range str {
		if c == separator {
			if collector.Len() > 0 {
				v, e := strconv.ParseInt(collector.String(), 10, 64)
				if e != nil {
					return ids, e
				}
				ids = append(ids, v)
				collector.Reset()
			}
		} else {
			collector.WriteRune(c)
		}
	}
	if collector.Len() > 0 {
		v, e := strconv.ParseInt(collector.String(), 10, 64)
		if e != nil {
			return ids, e
		}
		ids = append(ids, v)
		collector.Reset()
	}
	return ids, nil
}

// Convert node IDs into an ID path string
func IdsToParentsStr(ids []int64) string {
	var output strings.Builder
	for _, id := range ids {
		if id != 0 {
			output.WriteString(strconv.FormatInt(id, 10))
			output.WriteRune('/')
		}
	}
	return output.String()
}

/*
static std::vector<int64_t> get_parents_node_ids(db& db, int64_t nodeId)
{
	if (nodeId == 0)
		return std::vector<int64_t>();

	auto parents_ids_str_opt =
		db.execScalarString(L"SELECT parents FROM nodes WHERE id = @nodeId", { {L"@nodeId", nodeId} });

	if (!parents_ids_str_opt.has_value())
		throw nldberr("get_parents_node_ids: Node not found: " + std::to_string(nodeId));
	else
		return str_to_ids(parents_ids_str_opt.value(), '/');
}

void checkName(const std::wstring& name)
{
	if (name.find('/') != std::wstring::npos)
		throw nldberr("Invalid node name, cannot contain /");
}

node nodes::create(db& db, int64_t parentNodeId, int64_t nameStringId, int64_t typeStringId, const std::optional<std::wstring>& payload)
{
	checkName(strings::get_val(db, nameStringId));

	auto parent_node_ids = get_parents_node_ids(db, parentNodeId);
	if (parentNodeId != 0)
		parent_node_ids.push_back(parentNodeId);

	std::wstring parents_str = ids_to_parents_str(parent_node_ids);

	int64_t new_id = -1;
	if (parents_str.empty() && !payload.has_value())
	{
		new_id =
			db.execInsert
			(
				L"INSERT INTO nodes (parent_id, name_string_id, type_string_id) "
				L"VALUES (@parentNodeId, @nameStringId, @typeStringId)",
				{
					{ L"@parentNodeId", parentNodeId },
					{ L"@nameStringId", nameStringId },
					{ L"@typeStringId", typeStringId },
				}
			);
	}
	else if (!payload.has_value())
	{
		new_id =
			db.execInsert
			(
				L"INSERT INTO nodes (parent_id, name_string_id, type_string_id, parents) "
				L"VALUES (@parentNodeId, @nameStringId, @typeStringId, @parents)",
				{
					{ L"@parentNodeId", parentNodeId },
					{ L"@nameStringId", nameStringId },
					{ L"@typeStringId", typeStringId },
					{ L"@parents", parents_str },
				}
			);
	}
	else if (parents_str.empty())
	{
		new_id =
			db.execInsert
			(
				L"INSERT INTO nodes (parent_id, name_string_id, type_string_id, payload) "
				L"VALUES (@parentNodeId, @nameStringId, @typeStringId, @payload)",
				{
					{ L"@parentNodeId", parentNodeId },
					{ L"@nameStringId", nameStringId },
					{ L"@typeStringId", typeStringId },
					{ L"@payload", payload.value() }
				}
			);
	}
	else
	{
		new_id =
			db.execInsert
			(
				L"INSERT INTO nodes (parent_id, name_string_id, type_string_id, parents, payload) "
				L"VALUES (@parentNodeId, @nameStringId, @typeStringId, @parents, @payload)",
				{
					{ L"@parentNodeId", parentNodeId },
					{ L"@nameStringId", nameStringId },
					{ L"@typeStringId", typeStringId },
					{ L"@parents", parents_str },
					{ L"@payload", payload.value() }
				}
			);
	}
	return node(new_id, parentNodeId, nameStringId, typeStringId, payload);
}


/*
void copy(db& db, int64_t nodeId, int64_t newParentNodeId);
void move(db& db, int64_t nodeId, int64_t newParentNodeId);
void remove(db& db, int64_t nodeId);
void rename(db& db, int64_t nodeId, int64_t newNameStringId);

std::wstring get_payload(db& db, int64_t nodeId);
void set_payload(db& db, int64_t nodeId, const std::wstring& payload);

node get(db& db, int64_t nodeId);
void invalidate_cache(int64_t nodeId);
void invalidate_cache(const std::vector<int64_t>& nodeIds);
void flush_cache();
std::optional<node> get_node_in_parent(db& db, int64_t parentNodeId, int64_t nameStringId);

node get_parent(db& db, int64_t nodeId);
std::vector<node> get_parents(db& db, int64_t nodeId);
std::vector<node> get_children(db& db, int64_t nodeId);
std::vector<node> get_all_children(db& db, int64_t nodeId);

std::vector<node> get_path(db& db, const node& cur);
std::wstring get_path_str(db& db, const node& cur);
std::optional<std::vector<node>> get_path_nodes(db& db, const std::wstring& path);

std::optional<std::wstring> get_path_to_parent_like(db& db, const std::wstring& path);
*/
