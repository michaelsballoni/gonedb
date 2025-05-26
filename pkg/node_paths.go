package gonedb

import (
	"database/sql"
	"slices"
)

type node_paths struct{}

// The API you interact with
var NodePaths node_paths

func (np *node_paths) GetNodes(db *sql.DB, cur Node) ([]Node, error) {
	output := []Node{}
	seen_node_ids := map[int64]bool{}
	cur_node := cur
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

/* FORNOW - Needed?
func (np *node_paths) GetStr(db *sql.DB, nodeId int64) (string, error) {
	// FORNOW
	return "", nil
}

func (np *node_paths) GetStrNodes(db *sql.DB, nodeId int64) ([]Node, error) {
	// FORNOW
	return []Node{}, nil
}
*/
