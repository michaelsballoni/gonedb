package test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unicode"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestCmd(t *testing.T) {
	// set up shop
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

	// make the root node, seed with temp file system directory
	output, err = cmd.ProcessCommand(db, "make root")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "seed \""+playground_dir+"\"")
	AssertNoError(err)
	AssertEqual("", output)

	// go into dir1 and create deeper
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

	// go back to root and create new parent for copy of dir1
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

	// make dir3 in new_dir1_parent next to dir1 copy
	output, err = cmd.ProcessCommand(db, "make dir3")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/dir1\nroot/new_dir1_parent/dir3\n", output)

	// ensure dir1 copy has copy of deeper
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/dir1/deeper\n", output)

	// move new_dir1_parent's deeper node directly into new_dir1_parent
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/dir1/deeper")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "move root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	// ensure deeper copy is now in new_dir1_parent
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/deeper\nroot/new_dir1_parent/dir1\nroot/new_dir1_parent/dir3\n", output)

	// ensure copy of dir1 in new_dir1_parent has no children
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("", output)

	// rename the copy of deeper
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/deeper")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "rename about_to_remove_deeper")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/about_to_remove_deeper\nroot/new_dir1_parent/dir1\nroot/new_dir1_parent/dir3\n", output)

	// remove the copy of deeper
	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent/about_to_remove_deeper")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "remove")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "cd root/new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "dir")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent/dir1\nroot/new_dir1_parent/dir3\n", output)

	// rename new_dir1_parent to new_dir1_parent2, then search for it
	output, err = cmd.ProcessCommand(db, "rename new_dir1_parent2")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "search name new_dir1_parent")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "search name new_dir1_parent2")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent2\n", output)

	output, err = cmd.ProcessCommand(db, "search prop1 value1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "setprop prop1 value1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "search prop1 not-value1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "search prop1 value1")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent2\n", output)

	output, err = cmd.ProcessCommand(db, "search payload payload1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "setpayload payload1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "search payload payload1")
	AssertNoError(err)
	AssertEqual("root/new_dir1_parent2\n", output)

	output, err = cmd.ProcessCommand(db, "link root/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "unlink root/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "link root/dir1")
	AssertNoError(err)
	AssertEqual("", output)

	output, err = cmd.ProcessCommand(db, "tell")
	AssertNoError(err)
	AssertTrue(strings.Contains(output, "\nName: new_dir1_parent2\n"))
	AssertTrue(strings.Contains(output, "\nParent: root\n"))
	AssertTrue(strings.Contains(output, "\nPayload: payload1\n"))
	AssertTrue(strings.Contains(output, "\nProperties:\nprop1 value1\n"))
	AssertTrue(strings.Contains(output, "\nOut Links: (1)\nroot/dir1\n"))
	AssertTrue(strings.Contains(output, "\nIn Links: (none)\n"))

	output, err = cmd.ProcessCommand(db, "scramblelinks")
	AssertNoError(err)
	count_str := strings.Replace(output, "links created: ", "", 1)
	count, atoi_err := strconv.Atoi(count_str)
	AssertNoError(atoi_err)
	AssertTrue(int(count) > 0)

	output, err = cmd.BloomCloud(db)
	AssertNoError(err)
	fmt.Printf("BLOOM CLOUD OUTPUT:\n%s\n", output)
	bloom_scanner := bufio.NewScanner(strings.NewReader(output))
	for gen := 1; gen <= 3; gen += 1 {
		var cur_line string

		for {
			bloom_scanner.Scan()
			cur_line = bloom_scanner.Text()
			if len(cur_line) > 0 && !unicode.IsDigit(rune(cur_line[0])) {
				break
			}
		}
		AssertEqual(fmt.Sprintf("Gen: %d", gen), cur_line)

		bloom_scanner.Scan()
		cur_line = bloom_scanner.Text()
		line_prefix := fmt.Sprintf("Gen: %d - Added: ", gen)
		add_prefix := strings.Index(cur_line, line_prefix)
		AssertEqual(0, add_prefix)
		added_str := cur_line[len(line_prefix):]
		added_count, err := strconv.Atoi(added_str)
		AssertNoError(err)
		if added_count <= 0 {
			break
		}

		bloom_scanner.Scan()
		cur_line = bloom_scanner.Text()
		AssertEqual(0, strings.Index(cur_line, fmt.Sprintf("Gen: %d - Count: ", gen)))
	}

	err = cmd.DropCloud(db)
	AssertNoError(err)
}
