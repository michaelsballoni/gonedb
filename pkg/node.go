package gonedb

import (
	_ "github.com/mattn/go-sqlite3"
)

type Node struct {
	Id           int64
	ParentId     int64
	NameStringId int64
	TypeStringId int64
}

/*
func node CreateNode(db& db, int64_t parentNodeId, int64_t nameStringId, int64_t typeStringId = 0) node {

}

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
