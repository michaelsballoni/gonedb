package gonedb

import (
	"database/sql"
	"slices"
)

type node_paths struct{}

// The API you interact with
var NodePaths node_paths

// Given a node, return the list of anscestors up to and including it
func (np *node_paths) GetNodes(db *sql.DB, node Node) ([]Node, error) {
	output := []Node{}
	seen_node_ids := map[int64]bool{}
	cur_node := node
	for {
		if cur_node.Id == 0 {
			break
		} else {
			output = append(output, cur_node)
		}

		var parent_err error
		cur_node, parent_err = Nodes.GetParent(db, cur_node.Id)
		if parent_err != nil {
			return []Node{}, parent_err
		} else if seen_node_ids[cur_node.Id] {
			break
		} else {
			seen_node_ids[cur_node.Id] = true
		}
	}
	slices.Reverse(output)
	return output, nil
}

// Given a node return a list of path names up to and including the node
func (np *node_paths) GetStrs(db *sql.DB, node Node) ([]string, error) {
	path_nodes, nodes_err := np.GetNodes(db, node)
	if nodes_err != nil {
		return []string{}, nodes_err
	}

	path_str_ids := make([]int64, 0, len(path_nodes))
	for _, node := range path_nodes {
		path_str_ids = append(path_str_ids, node.NameStringId)
	}

	strs_map, strs_err := Strings.GetVals(db, path_str_ids)
	if strs_err != nil {
		return []string{}, strs_err
	}

	output := make([]string, 0, len(path_str_ids))
	for _, path_str_id := range path_str_ids {
		output = append(output, strs_map[path_str_id])
	}
	return output, nil
}

// Given a list of name strings, return a list of nodes in the path, or nil of the path does not resolve to a node
func (np *node_paths) GetStrNodes(db *sql.DB, pathParts []string) (*[]Node, error) {
	output := []Node{}
	var cur_node_id int64
	for _, part := range pathParts {
		if part == "" {
			continue
		}

		cur_name_string_id, str_err := Strings.GetId(db, part)
		if str_err != nil {
			return nil, str_err
		}

		node_in_parent, node_err := Nodes.GetNodeInParent(db, cur_node_id, cur_name_string_id)
		if node_err != nil {
			if node_err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, node_err
			}
		}

		output = append(output, node_in_parent)
		cur_node_id = node_in_parent.Id
	}
	return &output, nil
}
