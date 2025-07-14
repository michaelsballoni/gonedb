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

	var output string

	output, err = cmd.ProcessCommand(db, "make root")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "seed \""+playground_dir+"\"")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/dir1\nroot/dir2\n", output)

	output, err = cmd.ProcessCommand(db, "cd root/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "make deeper")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/dir1/deeper\n", output)

	output, err = cmd.ProcessCommand(db, "cd root")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "make new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "copy root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "make dir3")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/dir1\nroot/new_dir1_parent/dir3\n", output)

	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/dir1/deeper\n", output)
}
