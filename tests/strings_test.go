package test

import (
	"database/sql"
	"testing"

	gonedb "github.com/michaelsballoni/gonedb/pkg"
)

func get_val(t *testing.T, db *sql.DB, id int64) string {
	val, err := gonedb.Strings.GetVal(db, id)
	AssertNoError(t, err)
	return val
}

func TestStrings(t *testing.T) {
	db := GetTestDb("TestStrings.db")
	defer db.Close()

	_, e := gonedb.Strings.GetVal(db, 1200)
	AssertError(t, e)

	str1 := "foo"
	str2 := "bazr"

	id1, _ := gonedb.Strings.GetId(db, str1)
	if id1 != 1 {
		t.Errorf("GetId.id1: %d", id1)
		return
	}
	id1b, _ := gonedb.Strings.GetId(db, str1)
	if id1b != 1 {
		t.Errorf("GetId.id1b: %d", id1b)
		return
	}
	id1c, _ := gonedb.Strings.GetId(db, str1)
	if id1c != 1 {
		t.Errorf("GetId.id1c: %d", id1c)
		return
	}

	id2, _ := gonedb.Strings.GetId(db, str2)
	if id2 != 2 {
		t.Errorf("GetId.id2: %d", id2)
		return
	}
	id2b, _ := gonedb.Strings.GetId(db, str2)
	if id2b != 2 {
		t.Errorf("GetId.id2b: %d", id2b)
		return
	}

	AssertEqual(t, str1, get_val(t, db, id1))
	AssertEqual(t, str1, get_val(t, db, id1b))
	AssertEqual(t, str1, get_val(t, db, id1c))

	AssertEqual(t, str2, get_val(t, db, id2))
	AssertEqual(t, str2, get_val(t, db, id2b))

	gonedb.Strings.FlushCaches()

	strs, err := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(t, err)
	AssertEqual(t, str1, strs[id1])
	AssertEqual(t, str2, strs[id2])

	strs2, err2 := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(t, err2)
	AssertEqual(t, str1, strs2[id1])
	AssertEqual(t, str2, strs2[id2])

	gonedb.Strings.FlushCaches()

	strs3, err3 := gonedb.Strings.GetVals(db, []int64{id1, id2})
	AssertNoError(t, err3)
	AssertEqual(t, str1, strs3[id1])
	AssertEqual(t, str2, strs3[id2])

	_, err4 := gonedb.Strings.GetVals(db, []int64{id1, id2, -782})
	if err4 == nil {
		AssertTrue(t, false)
	}
}
