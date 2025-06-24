package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestLoader(t *testing.T) {
	db := GetTestDb("TestLoader.db")
	defer db.Close()

	// get root node
	root_node, node_err := gonedb.Nodes.Create(db, 0, 0, 0)
	AssertNoError(node_err)

	// get test exe path, as good of a file location to work with as any
	exe_path_str, exe_err := os.Executable()
	AssertNoError(exe_err)
	exe_abs_path, path_err := filepath.Abs(exe_path_str)
	AssertNoError(path_err)
	exe_path := filepath.Base(exe_abs_path)

	// do the load
	load_err := gonedb.Loader.Load(db, exe_path, root_node)
	AssertNoError(load_err)

	// ensure it all worked out
	cmp_err := OsCmp(db, exe_path, root_node)
	AssertNoError(cmp_err)
}

func OsCmp(db *sql.DB, curPath string, curNode gonedb.Node) error {
	entries, entry_err := os.ReadDir(curPath)
	if entry_err != nil {
		return entry_err
	}

	for _, entry := range entries {
		abs_path, abs_err := filepath.Abs(entry.Name())
		if abs_err != nil {
			return abs_err
		}
		name := entry.Name()

		name_string_id, str_err := gonedb.Strings.GetId(db, name)
		if str_err != nil {
			return str_err
		}

		new_node, node_err := gonedb.Nodes.GetNodeInParent(db, curNode.Id, name_string_id)
		if node_err != nil {
			return node_err
		}
		if new_node.Id <= 0 {
			return fmt.Errorf("not found found with name")
		}
		if node_err != nil {
			return node_err
		}

		if entry.IsDir() {
			OsCmp(db, abs_path, new_node)
		}
	}

	return nil
}
