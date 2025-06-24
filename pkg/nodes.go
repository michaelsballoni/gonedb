package gonedb

import (
	"database/sql"
	"fmt"
	"strconv"
)

type nodes struct{}

// The API you interact with
var Nodes nodes

// The struct passed around the API with the core IDs of a node
type Node struct {
	Id           int64
	ParentId     int64
	NameStringId int64
	TypeStringId int64
}

// Get a node given an ID
func (n *nodes) Get(db *sql.DB, nodeId int64) (Node, error) {
	found_node, found := NodeCache.Get(nodeId)
	if found {
		return found_node, nil
	}

	var output Node
	output.Id = nodeId
	row := db.QueryRow("SELECT parent_id, name_string_id, type_string_id FROM nodes WHERE id = ?", nodeId)
	err := row.Scan(&output.ParentId, &output.NameStringId, &output.TypeStringId)
	if err != nil {
		return Node{}, err
	} else {
		NodeCache.Put(output)
		return output, nil
	}
}

// Create a new node
func (n *nodes) Create(db *sql.DB, parentNodeId int64, nameStringId int64, typeStringId int64) (Node, error) {
	// Add parents of the parent
	parent_node_ids, err := n.GetParentsNodeIds(db, parentNodeId)
	if err != nil {
		return Node{}, err
	} else if parentNodeId != 0 { // append the parent to its path to make our own
		parent_node_ids = append(parent_node_ids, parentNodeId)
	}

	parents_str := NodeUtils.IdsToParentsStr(parent_node_ids)

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

// Copy a node to another parent node, deep copy.
func (n *nodes) Copy(db *sql.DB, nodeId int64, newParentNodeId int64) (int64, error) {
	if nodeId == newParentNodeId {
		return 0, fmt.Errorf("Cannot copy node into itself")
	}

	// You can't copy this into a child of this
	{
		like_str, like_err := n.GetChildNodesLikeExpression(db, nodeId)
		if like_err != nil {
			return 0, nil
		}

		rows, rows_err := db.Query("SELECT id FROM nodes WHERE parents LIKE ?", like_str)
		if rows_err != nil {
			return 0, rows_err
		}
		var cur_id int64
		for rows.Next() {
			scan_err := rows.Scan(&cur_id)
			if scan_err != nil {
				return 0, scan_err
			}
			if cur_id == newParentNodeId {
				return 0, fmt.Errorf("Cannot copy node into child of source")
			}
		}
	}

	// we track all nodes we've seen to prevent
	src_node, src_err := n.Get(db, nodeId)
	if src_err != nil {
		return -1, src_err
	}

	seen_node_ids := make(map[int64]bool)
	new_node_id, copy_err := doCopy(db, src_node, newParentNodeId, seen_node_ids)
	return new_node_id, copy_err
}

// Recursive Copy workhorse routine
func doCopy(db *sql.DB, srcNode Node, newParentNodeId int64, seenNodeIds map[int64]bool) (int64, error) {
	new_node, new_err := Nodes.Create(db, newParentNodeId, srcNode.NameStringId, srcNode.TypeStringId)
	if new_err != nil {
		return -1, new_err
	}

	child_nodes, cur_err := Nodes.GetChildren(db, srcNode.Id)
	if cur_err != nil {
		return -1, cur_err
	}
	for _, child_node := range child_nodes {
		if seenNodeIds[child_node.Id] {
			return -1, fmt.Errorf("Cannot copy a node into its children")
		} else {
			seenNodeIds[child_node.Id] = true
		}
		_, copy_err := doCopy(db, child_node, new_node.Id, seenNodeIds)
		if copy_err != nil {
			return -1, copy_err
		}
	}
	return new_node.Id, nil
}

// Move a node to under a new parent node
func (n *nodes) Move(db *sql.DB, nodeId int64, newParentNodeId int64) error {
	// check inputs
	if nodeId == 0 {
		return fmt.Errorf("Cannot move null node")
	}

	// collect all children node IDs
	child_node_ids := []int64{}
	child_nodes_like, child_nodes_like_err := n.GetChildNodesLikeExpression(db, nodeId)
	if child_nodes_like_err != nil {
		return child_nodes_like_err
	}
	var cur_id int64
	query, query_err := db.Query("SELECT id FROM nodes WHERE parents LIKE ?", child_nodes_like)
	if query_err != nil {
		return query_err
	}
	for query.Next() {
		scan_err := query.Scan(&cur_id)
		if scan_err != nil {
			return scan_err
		}
		child_node_ids = append(child_node_ids, cur_id)
	}

	// compute the new parents for the node
	parents_node_ids, parents_err := n.GetParentsNodeIds(db, newParentNodeId)
	if parents_err != nil {
		return parents_err
	}
	if newParentNodeId != 0 {
		parents_node_ids = append(parents_node_ids, newParentNodeId)
	}
	new_parents_str := NodeUtils.IdsToParentsStr(parents_node_ids)

	// update the nodes parent and parents
	{
		result, err := db.Exec("UPDATE nodes SET parent_id = ?, parents = ? WHERE id = ?", newParentNodeId, new_parents_str, nodeId)
		if err != nil {
			return err
		}
		result_affected, _ := result.RowsAffected()
		if result_affected != 1 {
			return fmt.Errorf("Node not moved")
		}
	}
	NodeCache.Invalidate1(nodeId)

	// update the parents of all children nodes
	for _, child_id := range child_node_ids {
		child_parent_ids, child_parent_err := n.GetParentsNodeIds(db, child_id)
		if child_parent_err != nil {
			return child_parent_err
		}
		new_parents_str := NodeUtils.IdsToParentsStr(child_parent_ids)
		result, err := db.Exec("UPDATE nodes SET parents = ? WHERE id = ?", new_parents_str, child_id)
		if err != nil {
			return err
		}

		result_affected, _ := result.RowsAffected()
		if result_affected != 1 {
			return fmt.Errorf("Node not moved")
		}
		NodeCache.Invalidate1(child_id)
	}

	return nil
}

// Delete a node and all of its children
func (n *nodes) Remove(db *sql.DB, nodeId int64) error {
	// check inputs
	if nodeId == 0 {
		return fmt.Errorf("Cannot remove null node")
	}

	// collect all children node IDs
	child_node_ids, child_node_err := n.GetAllChildNodeIds(db, nodeId)
	if child_node_err != nil {
		return child_node_err
	}

	// delete the children nodes
	if len(child_node_ids) > 0 {
		ids, ids_err := NodeUtils.IdsToSqlIn(child_node_ids)
		if ids_err != nil {
			return ids_err
		}
		del_result, del_err := db.Exec("DELETE FROM nodes WHERE id IN (" + ids + ")")
		if del_err != nil {
			return del_err
		}
		del_count, _ := del_result.RowsAffected()
		if del_count != int64(len(child_node_ids)) {
			return fmt.Errorf("Not all child nodes removed")
		}
		NodeCache.InvalidateN(child_node_ids)
	}

	// delete the node
	{
		del_result, del_err := db.Exec("DELETE FROM nodes WHERE id = ?", nodeId)
		if del_err != nil {
			return del_err
		}
		del_count, _ := del_result.RowsAffected()
		if del_count != 1 {
			return fmt.Errorf("Node not removed")
		}
	}
	NodeCache.Invalidate1(nodeId)
	return nil
}

// Rename a node, ensuring an existing node in the same parent does not exist
func (n *nodes) Rename(db *sql.DB, nodeId int64, newNameStringId int64) error {
	if nodeId == 0 {
		return fmt.Errorf("Cannot rename null node")
	}

	parent_node, parent_err := n.GetParent(db, nodeId)
	if parent_err != nil {
		return parent_err
	}
	parent_id := parent_node.Id
	_, existing_err := n.GetNodeInParent(db, parent_id, newNameStringId)
	if existing_err == nil {
		return fmt.Errorf("Node with new name already exists")
	}

	result, err := db.Exec("UPDATE nodes SET name_string_id = ? WHERE id = ?", newNameStringId, nodeId)
	if err != nil {
		return err
	} else {
		result_count, result_err := result.RowsAffected()
		if result_err != nil {
			return result_err
		} else if result_count != 1 {
			return fmt.Errorf("Node not renamed")
		}
	}

	NodeCache.Invalidate1(nodeId)
	return nil
}

// Get the payload of a node
func (n *nodes) GetPayload(db *sql.DB, nodeId int64) (string, error) {
	var output string
	row := db.QueryRow("SELECT payload FROM nodes WHERE id = ?", nodeId)
	err := row.Scan(&output)
	return output, err
}

// Set the payload of a node
func (n *nodes) SetPayload(db *sql.DB, nodeId int64, payload string) error {
	result, err := db.Exec("UPDATE nodes SET payload = ? WHERE id = ?", payload, nodeId)
	if err != nil {
		return err
	}
	affected, affected_err := result.RowsAffected()
	if affected_err != nil {
		return affected_err
	}
	if affected != 1 {
		return fmt.Errorf("Row not affected")
	} else {
		return nil
	}
}

// Get the parent of a node
func (n *nodes) GetParent(db *sql.DB, nodeId int64) (Node, error) {
	var node Node
	row := db.QueryRow("SELECT id, parent_id, name_string_id, type_string_id FROM nodes WHERE id = (SELECT parent_id FROM nodes WHERE id = ?)", nodeId)
	err := row.Scan(&node.Id, &node.ParentId, &node.NameStringId, &node.TypeStringId)
	return node, err
}

// Get the node in a parent by anem
func (n *nodes) GetNodeInParent(db *sql.DB, parentNodeId int64, nameStringId int64) (Node, error) {
	var node Node
	node.ParentId = parentNodeId
	node.NameStringId = nameStringId
	row := db.QueryRow("SELECT id, type_string_id FROM nodes WHERE parent_id = ? AND name_string_id = ?", parentNodeId, nameStringId)
	err := row.Scan(&node.Id, &node.TypeStringId)
	return node, err
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
		return NodeUtils.StringToIds(parents_ids_str)
	}
}

// Get the IDs of all children of a given node ID, deep
func (n *nodes) GetAllChildNodeIds(db *sql.DB, nodeId int64) ([]int64, error) {
	child_node_ids := []int64{}
	child_nodes_like, child_nodes_like_err := n.GetChildNodesLikeExpression(db, nodeId)
	if child_nodes_like_err != nil {
		return []int64{}, child_nodes_like_err
	}
	query, query_err := db.Query("SELECT id FROM nodes WHERE parents LIKE ?", child_nodes_like)
	if query_err != nil {
		return []int64{}, query_err
	}
	var cur_id int64
	for query.Next() {
		scan_err := query.Scan(&cur_id)
		if scan_err != nil {
			return []int64{}, scan_err
		}
		child_node_ids = append(child_node_ids, cur_id)
	}
	return child_node_ids, nil
}

// Get the children node structs of a given node by ID, shallow
func (n *nodes) GetChildren(db *sql.DB, nodeId int64) ([]Node, error) {
	child_nodes := []Node{}
	rows, err := db.Query("SELECT id, parent_id, name_string_id, type_string_id FROM nodes WHERE id <> 0 AND parent_id = ?", nodeId)
	if err != nil {
		return []Node{}, err
	}
	var cur_output_node Node
	for rows.Next() {
		scan_err := rows.Scan(&cur_output_node.Id, &cur_output_node.ParentId, &cur_output_node.NameStringId, &cur_output_node.TypeStringId)
		if scan_err != nil {
			return []Node{}, scan_err
		} else {
			child_nodes = append(child_nodes, cur_output_node)
		}
	}
	return child_nodes, nil
}

// Get the LIKE expression for returing the child nodes in a SQL query agains the nodes table
func (n *nodes) GetChildNodesLikeExpression(db *sql.DB, nodeId int64) (string, error) {
	var original_node_parents string
	row := db.QueryRow("SELECT parents FROM nodes WHERE id = ?", nodeId)
	err := row.Scan(&original_node_parents)
	if err != nil {
		return "", err
	} else if original_node_parents == "" || original_node_parents[len(original_node_parents)-1] != '/' {
		original_node_parents = "/"
	}
	original_node_parents = original_node_parents + strconv.FormatInt(nodeId, 10) + "/%"
	return original_node_parents, nil
}
