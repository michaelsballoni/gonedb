package gonedb

import (
	"database/sql"
	"os"
	"path/filepath"
)

type loader struct{}

// The API you interact with
var Loader loader

// Load a file system directory into a node
func (loader *loader) Load(db *sql.DB, curPath string, curNode Node) error {
	// Get the file system entries at the current level
	entries, entry_err := os.ReadDir(curPath)
	if entry_err != nil {
		return entry_err
	}

	// walk the entries
	for _, entry := range entries {
		abs_path, abs_err := filepath.Abs(entry.Name())
		if abs_err != nil {
			return abs_err
		}
		name := entry.Name()

		// get the node's name
		name_string_id, str_err := Strings.GetId(db, name)
		if str_err != nil {
			return str_err
		}

		// get or create the child node
		new_node, node_err := Nodes.GetNodeInParent(db, curNode.Id, name_string_id)
		if node_err != nil {
			return node_err
		}
		if new_node.Id <= 0 {
			var new_node_type_id int64 = 0
			new_node, node_err = Nodes.Create(db, curNode.Id, name_string_id, new_node_type_id)
		}
		if node_err != nil {
			return node_err
		}

		// recurse on sub-child-dirs
		if entry.IsDir() {
			loader.Load(db, abs_path, new_node)
		}
	}

	// all done
	return nil
}
