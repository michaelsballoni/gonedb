package test

import (
	"database/sql"
	"fmt"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func TestStrings(t *testing.T) {
	db := GetTestDb("StringsTests.db")
	defer db.Close()

	_, e := gonedb.Strings.GetVal(db, 1200)
	AssertError(e)

	str1 := "foo"
	str2 := "bazr"

	id1, _ := gonedb.Strings.GetId(db, str1)
	AssertEqual(1, id1)
	id1b, _ := gonedb.Strings.GetId(db, str1)
	AssertEqual(1, id1b)
	id1c, _ := gonedb.Strings.GetId(db, str1)
	AssertEqual(1, id1c)

	id2, _ := gonedb.Strings.GetId(db, str2)
	AssertEqual(2, id2)
	id2b, _ := gonedb.Strings.GetId(db, str2)
	AssertEqual(2, id2b)

	AssertEqual(str1, get_val(db, id1))
	AssertEqual(str1, get_val(db, id1b))
	AssertEqual(str1, get_val(db, id1c))

	AssertEqual(str2, get_val(db, id2))
	AssertEqual(str2, get_val(db, id2b))

	gonedb.Strings.FlushCaches()

	strs, err := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(err)
	AssertEqual(str1, strs[id1])
	AssertEqual(str2, strs[id2])

	strs2, err2 := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(err2)
	AssertEqual(str1, strs2[id1])
	AssertEqual(str2, strs2[id2])

	gonedb.Strings.FlushCaches()

	strs3, err3 := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(err3)
	AssertEqual(str1, strs3[id1])
	AssertEqual(str2, strs3[id2])

	_, err4 := gonedb.Strings.GetVals(db, []int64{id1, id2, -782})
	AssertError((err4))
}

func TestStringsMixup(t *testing.T) {
	db := GetTestDb("StringsTests2.db")
	defer db.Close()
	strs := []string{
		"foobar",
		"root",
		"bletmonkey",
		"child",
		"funkadelic",
		"grandkid",
		"superfly",
	}
	seen_ids := map[int64]bool{} // string ID -> dummy bool
	for _, str := range strs {
		fmt.Printf("GetId: str: %s\n", str)
		new_id, new_err := gonedb.Strings.GetId(db, str)
		AssertNoError(new_err)
		_, found_ok := seen_ids[new_id]
		fmt.Printf("seen_ids: new_id: %d - ok: %t\n", new_id, found_ok)
		AssertTrue(!found_ok)
		seen_ids[new_id] = true
	}

	for _, str := range strs {
		new_id, new_err := gonedb.Strings.GetId(db, str)
		AssertNoError(new_err)
		val := get_val(db, new_id)
		AssertEqual(str, val)
	}
}

func get_val(db *sql.DB, id int64) string {
	val, err := gonedb.Strings.GetVal(db, id)
	AssertNoError(err)
	return val
}
