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
	//fmt.Printf("Load: %s\n", curPath)
	entries, entry_err := os.ReadDir(curPath)
	if len(entries) == 0 || entry_err != nil { // ignore empty / sharing errors, etc.
		return nil
	}
	//fmt.Printf("Load: %s -> %d\n", curPath, len(entries))

	// walk the entries
	for _, entry := range entries {
		// get the node's name
		name := entry.Name()
		name_string_id, str_err := Strings.GetId(db, name)
		if str_err != nil {
			return str_err
		}

		// get or create the child node
		new_node, node_err := Nodes.GetNodeInParent(db, curNode.Id, name_string_id)
		if new_node.Id < 0 {
			var new_node_type_id int64 = 0
			new_node, node_err = Nodes.Create(db, curNode.Id, name_string_id, new_node_type_id)
		}
		if node_err != nil {
			return node_err
		}

		// recurse on sub-child-dirs
		if entry.IsDir() {
			abs_path := filepath.Join(curPath, name)
			load_err := loader.Load(db, abs_path, new_node)
			if load_err != nil {
				return load_err
			}
		}
	}

	// all done
	return nil
}
