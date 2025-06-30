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
	root_node := gonedb.Node{}

	// get test input dir
	exe_path, exe_err := os.Executable()
	AssertNoError((exe_err))
	test_path := filepath.Dir(exe_path) + string(filepath.Separator) + ".." + string(filepath.Separator)

	// do the load
	load_err := gonedb.Loader.Load(db, test_path, root_node)
	AssertNoError(load_err)

	// ensure it all worked out
	cmp_err := OsCmp(db, test_path, root_node)
	AssertNoError(cmp_err)
}

func OsCmp(db *sql.DB, curPath string, curNode gonedb.Node) error {
	fmt.Printf("OsCmp: %s\n", curPath)
	entries, entry_err := os.ReadDir(curPath)
	if len(entries) == 0 || entry_err != nil { // ignore empty / sharing errors, etc.
		return nil
	}
	fmt.Printf("OsCmp: %s -> %d\n", curPath, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		name_string_id, str_err := gonedb.Strings.GetId(db, name)
		if str_err != nil {
			return str_err
		}

		new_node, node_err := gonedb.Nodes.GetNodeInParent(db, curNode.Id, name_string_id)
		if new_node.Id < 0 {
			return fmt.Errorf("not found found with name")
		}
		if node_err != nil {
			return node_err
		}

		if entry.IsDir() {
			abs_path := filepath.Join(curPath, name)
			load_err := OsCmp(db, abs_path, new_node)
			if load_err != nil {
				return load_err
			}
		}
	}

	return nil
}
