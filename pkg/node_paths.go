package gonedb

import (
	"database/sql"
	"slices"
	"strings"
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

// Given a node return a /-delimited path of name string IDs up to and including the node
func (np *node_paths) GetNodesStr(db *sql.DB, node Node) (string, error) {
	path_nodes, nodes_err := np.GetNodes(db, node)
	if nodes_err != nil {
		return "", nodes_err
	}

	path_str_ids := make([]int64, 0, len(path_nodes))
	for _, node := range path_nodes {
		path_str_ids = append(path_str_ids, node.NameStringId)
	}

	strs_map, strs_err := Strings.GetVals(db, path_str_ids)
	if strs_err != nil {
		return "", strs_err
	}

	var builder strings.Builder
	for _, path_str_id := range path_str_ids {
		builder.WriteRune('/')
		builder.WriteString(strs_map[path_str_id])
	}
	return builder.String(), nil
}

// Given a /-separated path of name string IDs, return a list of nodes in the path
func (np *node_paths) GetStrNodes(db *sql.DB, path string) (*[]Node, error) {
	splits := []string{}
	var builder strings.Builder
	for _, c := range path {
		if c == '/' {
			if builder.Len() > 0 {
				splits = append(splits, builder.String())
				builder.Reset()
			}
		} else {
			builder.WriteRune(c)
		}
	}
	if builder.Len() > 0 {
		splits = append(splits, builder.String())
	}
	if len(splits) == 0 {
		return nil, nil
	}

	output := make([]Node, 0, len(splits))
	var cur_node_id int64
	for _, part := range splits {
		cur_name_string_id, err := Strings.GetId(db, part)
		if err != nil {
			return nil, nil
		}
		node_in_parent, node_err := Nodes.GetNodeInParent(db, cur_node_id, cur_name_string_id)
		if node_err != nil {
			return nil, node_err
		}
		output = append(output, node_in_parent)
		cur_node_id = node_in_parent.Id
	}

	return &output, nil
}
