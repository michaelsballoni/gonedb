package gonedb

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
)

// The API you interact with
type nodes struct{}

var Nodes nodes

// The struct passed around the API with the core IDs of a node
type Node struct {
	Id           int64
	ParentId     int64
	NameStringId int64
	TypeStringId int64
}

// Turn IDs into a string used in a SELECT IN (here)
// Errors if an empty input is provided
func (n *nodes) IdsToSqlIn(ids []int64) (string, error) {
	if len(ids) == 0 {
		return "", errors.New("ids is empty")
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

// Turn a separator-delimited string into a slice of IDs
// only IDs found between separators are returned
// nothing in, nothing out
func (n *nodes) StringToIds(str string) ([]int64, error) {
	sep := "/"
	if str == "" || str == sep {
		return []int64{}, nil
	}

	strs := strings.Split(str, sep)
	ids := make([]int64, len(strs))
	for i, v := range strs {
		n, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return []int64{}, e
		}
		ids[i] = n
	}
	return ids, nil
}

// Convert IDs into an ID path string
func (n *nodes) IdsToParentsStr(ids []int64) string {
	if len(ids) == 0 || (len(ids) == 1 && ids[0] == 0) {
		return ""
	}

	strs := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != 0 {
			strs = append(strs, strconv.FormatInt(id, 10))
		}
	}
	return strings.Join(strs, string('/')) + "/"
}

// Create a new node
func (n *nodes) Create(db *sql.DB, parentNodeId int64, nameStringId int64, typeStringId int64) (Node, error) {
	parent_node_ids, err := n.GetParentsNodeIds(db, parentNodeId)
	if err != nil {
		return Node{}, err
	}
	if parentNodeId != 0 {
		parent_node_ids = append(parent_node_ids, parentNodeId)
	}

	parents_str := n.IdsToParentsStr(parent_node_ids)

	var new_id int64
	row := db.QueryRow("INSERT INTO nodes (parent_id, name_string_id, type_string_id, parents) VALUES (?, ?, ?, ?) RETURNING id",
		parentNodeId,
		nameStringId,
		typeStringId,
		parents_str)
	err = row.Scan(&new_id)
	if err != nil {
		return Node{}, err
	} else {
		return Node{Id: new_id, ParentId: parentNodeId, NameStringId: nameStringId, TypeStringId: typeStringId}, nil
	}
}

func (n *nodes) Get(db *sql.DB, nodeId int64) (Node, error) {
	found_node, found := n.GetFromCache(nodeId)
	if found {
		return found_node, nil
	}

	var ret_node Node
	ret_node.Id = nodeId
	row := db.QueryRow("SELECT parent_id, name_string_id, type_string_id FROM nodes WHERE id = ?", nodeId)
	err := row.Scan(&ret_node.ParentId, &ret_node.NameStringId, &ret_node.TypeStringId)
	if err != nil {
		return Node{}, err
	} else {
		n.PutIntoCache(ret_node)
		return ret_node, nil
	}
}

// Get the node ID parents of a node
func (n *nodes) GetParentsNodeIds(db *sql.DB, nodeId int64) ([]int64, error) {
	if nodeId == 0 {
		return []int64{}, nil
	}

	var parents_ids_str string
	row := db.QueryRow("SELECT parents FROM nodes WHERE id = ?", nodeId)
	err := row.Scan(&parents_ids_str)
	if err != nil {
		return []int64{}, err
	} else {
		return n.StringToIds(parents_ids_str)
	}
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

var g_cacheLock sync.RWMutex
var g_cache = make(map[int64]Node)

func (n *nodes) GetFromCache(id int64) (Node, bool) {
	g_cacheLock.RLock()
	defer g_cacheLock.RUnlock()
	node, found := g_cache[id]
	return node, found
}

func (n *nodes) PutIntoCache(node Node) {
	g_cacheLock.Lock()
	g_cache[node.Id] = node
	g_cacheLock.Unlock()
}

func (n *nodes) FlushCache() {
	g_cacheLock.Lock()
	clear(g_cache)
	g_cacheLock.Unlock()
}

func (n *nodes) InvalidateCache1(nodeId int64) {
	g_cacheLock.Lock()
	delete(g_cache, nodeId)
	g_cacheLock.Unlock()
}

func (n *nodes) InvalidateCacheN(nodeIds []int64) {
	g_cacheLock.Lock()
	for _, nodeId := range nodeIds {
		delete(g_cache, nodeId)
	}
	g_cacheLock.Unlock()
}
