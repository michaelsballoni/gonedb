package gonedb

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type nlstrings struct{}

var Strings nlstrings

// Given a string value, ensure it's in the DB and return the ID for it
func (s *nlstrings) GetId(db *sql.DB, str string) (int64, error) {
	if len(str) == 0 {
		return 0, nil
	}

	string_id := get_id_from_cache(str)
	if string_id >= 0 {
		return string_id, nil
	}

	row := db.QueryRow("SELECT id FROM strings WHERE val = ?", str)
	err := row.Scan(&string_id)
	if err == nil {
		put_id_in_cache(str, string_id)
		return string_id, nil
	}

	row = db.QueryRow("INSERT INTO strings (val) VALUES (?) on conflict(val) DO UPDATE SET val = val RETURNING id", str)
	err = row.Scan(&string_id)
	if err == nil {
		put_id_in_cache(str, string_id)
		return string_id, nil
	} else {
		return -1, err
	}
}

// Given a database ID, try to find the correspdonding string value
func (s *nlstrings) GetVal(db *sql.DB, id int64) (string, error) {
	if id == 0 {
		return "", nil
	}

	str, found := get_val_from_cache(id)
	if found {
		return str, nil
	}

	row := db.QueryRow("SELECT val FROM strings WHERE id = ?", id)
	err := row.Scan(&str)
	if err != nil {
		return "", fmt.Errorf("Strings.GetVal: Value not found: %d: %v", id, err)
	}

	put_val_in_cache(id, str)
	return str, nil
}

// Given a slice of database IDs, return a map from ID to value
func (s *nlstrings) GetVals(db *sql.DB, ids []int64) (map[int64]string, error) {
	output := map[int64]string{}
	if len(ids) == 0 {
		return output, nil
	}

	ids_str := get_ids_to_vals_cache(ids, output)

	if len(ids_str) > 0 {
		var id int64
		var val string
		var found_vals = map[int64]string{}

		rows, err := db.Query("SELECT id, val FROM strings WHERE id IN (" + ids_str + ")")
		if err != nil {
			return map[int64]string{}, err
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&id, &val)
			output[id] = val
			found_vals[id] = val
		}

		if len(found_vals) > 0 {
			put_vals_in_cache(&found_vals)
		}
	}

	for _, id := range ids {
		_, ok := output[id]
		if !ok {
			return nil, fmt.Errorf("Strings.GetVals: String not found: %d", id)
		}
	}

	return output, nil
}

// Flush the string -> ID cache
func (s *nlstrings) FlushValToIdCache() {
	g_toIdCacheLock.Lock()
	clear(g_toIdCache)
	g_toIdCacheLock.Unlock()
}

// Flush the ID -> string cache
func (s *nlstrings) FlushFromValToIdCache() {
	g_fromIdCacheLock.Lock()
	clear(g_fromIdCache)
	g_fromIdCacheLock.Unlock()
}

// Flush all string-related caches
func (s *nlstrings) FlushCaches() {
	s.FlushValToIdCache()
	s.FlushFromValToIdCache()
}

var g_toIdCacheLock sync.RWMutex
var g_toIdCache = make(map[string]int64)

var g_fromIdCacheLock sync.RWMutex
var g_fromIdCache = make(map[int64]string)

func get_id_from_cache(str string) int64 {
	g_toIdCacheLock.RLock()
	defer g_toIdCacheLock.RUnlock()
	string_id, ok := g_toIdCache[str]
	if ok {
		return string_id
	} else {
		return -1
	}
}

func put_id_in_cache(str string, id int64) {
	g_toIdCacheLock.Lock()
	defer g_toIdCacheLock.Unlock()
	g_toIdCache[str] = id
}

func get_val_from_cache(id int64) (string, bool) {
	g_fromIdCacheLock.RLock()
	defer g_fromIdCacheLock.RUnlock()
	str, ok := g_fromIdCache[id]
	if ok {
		return str, true
	} else {
		return "", false
	}
}

func put_val_in_cache(id int64, str string) {
	g_fromIdCacheLock.Lock()
	defer g_fromIdCacheLock.Unlock()
	g_fromIdCache[id] = str
}

func put_vals_in_cache(vals *map[int64]string) {
	g_fromIdCacheLock.Lock()
	defer g_fromIdCacheLock.Unlock()
	for id, str := range *vals {
		g_fromIdCache[id] = str
	}
}

func get_ids_to_vals_cache(ids []int64, output map[int64]string) string {
	g_fromIdCacheLock.RLock()
	defer g_fromIdCacheLock.RUnlock()

	var ids_output strings.Builder
	for _, id := range ids {
		str, ok := g_fromIdCache[id]
		if ok {
			output[id] = str
		} else {
			if ids_output.Len() > 0 {
				ids_output.WriteRune(',')
			}
			ids_output.WriteString(strconv.FormatInt(id, 10))
		}
	}
	return ids_output.String()
}
