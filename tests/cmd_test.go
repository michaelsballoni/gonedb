package test

import (
	"os"
	"path/filepath"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestCmd(t *testing.T) {
	db := GetTestDb("TestCmd.db")
	defer db.Close()

	playground_dir, err := os.MkdirTemp("", "gonedb-playground-*")
	AssertNoError(err)
	defer os.RemoveAll(playground_dir)

	err = os.Mkdir(filepath.Join(playground_dir, "dir1"), 0700)
	AssertNoError(err)

	err = os.Mkdir(filepath.Join(playground_dir, "dir2"), 0700)
	AssertNoError(err)

	cmd := gonedb.CreateCmd()
	cmd.ProcessCommand(db, "make root")
	cmd.ProcessCommand(db, "cd root")
	cmd.ProcessCommand(db, "seed \""+playground_dir+"\"")

	var output string
	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/dir1\nroot/dir2\n", output)
}
