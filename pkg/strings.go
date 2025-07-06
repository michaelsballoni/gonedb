package gonedb

import (
	"database/sql"
	"fmt"
	"sort"
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

	_, exec_err := db.Exec("INSERT INTO strings (val) VALUES (?) ON CONFLICT(val) DO NOTHING", str)
	if exec_err != nil {
		return -1, exec_err
	}

	row = db.QueryRow("SELECT id FROM strings WHERE val = ?", str)
	err = row.Scan(&string_id)
	if err == nil {
		put_id_in_cache(str, string_id)
		return string_id, nil
	} else {
		return -1, exec_err
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
		for rows.Next() {
			rows.Scan(&id, &val)
			output[id] = val
			found_vals[id] = val
		}
		rows.Close()

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

// Take a stringID-to-stringID map and return a stringValue-to-stringValue map
func (s *nlstrings) Fill(db *sql.DB, nameValueStrIds map[int64]int64) (map[string]string, error) {
	output := map[string]string{}
	var name_str, val_str string
	var cur_err error
	for name_string_id, value_string_id := range nameValueStrIds {
		name_str, cur_err = s.GetVal(db, name_string_id)
		if cur_err != nil {
			return map[string]string{}, cur_err
		}
		val_str, cur_err = s.GetVal(db, value_string_id)
		if cur_err != nil {
			return map[string]string{}, cur_err
		}
		output[name_str] = val_str
	}
	return output, nil
}

// Take a stringID-to-stringID map and return a line-delimited string summarizing the names and values
func (s *nlstrings) Summarize(db *sql.DB, nameValueStrIds map[int64]int64) (string, error) {
	string_map, fill_err := s.Fill(db, nameValueStrIds)
	if fill_err != nil {
		return "", fill_err
	}

	names := make([]string, 0, len(string_map))
	for k := range string_map {
		names = append(names, k)
	}
	sort.Strings(names)

	var builder strings.Builder
	for _, name := range names {
		if builder.Len() > 0 {
			builder.WriteRune('\n')
		}
		builder.WriteString(name + " " + string_map[name])
	}
	return builder.String(), nil
}

// Flush the string -> ID cache
func (s *nlstrings) FlushValToIdCache() {
	g_strToIdCacheLock.Lock()
	clear(g_strToIdCache)
	g_strToIdCacheLock.Unlock()
}

// Flush the ID -> string cache
func (s *nlstrings) FlushFromValToIdCache() {
	g_fromIdToStrCacheLock.Lock()
	clear(g_fromIdToStrCache)
	g_fromIdToStrCacheLock.Unlock()
}

// Flush all string-related caches
func (s *nlstrings) FlushCaches() {
	s.FlushValToIdCache()
	s.FlushFromValToIdCache()
}

var g_strToIdCacheLock sync.RWMutex
var g_strToIdCache = make(map[string]int64)

func get_id_from_cache(str string) int64 {
	g_strToIdCacheLock.RLock()
	string_id, ok := g_strToIdCache[str]
	g_strToIdCacheLock.RUnlock()
	if ok {
		return string_id
	} else {
		return -1
	}
}

func put_id_in_cache(str string, id int64) {
	g_strToIdCacheLock.Lock()
	g_strToIdCache[str] = id
	g_strToIdCacheLock.Unlock()
}

func get_ids_to_vals_cache(ids []int64, output map[int64]string) string {
	g_fromIdToStrCacheLock.RLock()
	defer g_fromIdToStrCacheLock.RUnlock()
	var ids_output strings.Builder
	for _, id := range ids {
		str, ok := g_fromIdToStrCache[id]
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

var g_fromIdToStrCacheLock sync.RWMutex
var g_fromIdToStrCache = make(map[int64]string)

func get_val_from_cache(id int64) (string, bool) {
	g_fromIdToStrCacheLock.RLock()
	str, ok := g_fromIdToStrCache[id]
	g_fromIdToStrCacheLock.RUnlock()
	if ok {
		return str, true
	} else {
		return "", false
	}
}

func put_val_in_cache(id int64, str string) {
	g_fromIdToStrCacheLock.Lock()
	g_fromIdToStrCache[id] = str
	g_fromIdToStrCacheLock.Unlock()
}

func put_vals_in_cache(vals *map[int64]string) {
	g_fromIdToStrCacheLock.Lock()
	defer g_fromIdToStrCacheLock.Unlock()
	for id, str := range *vals {
		g_fromIdToStrCache[id] = str
	}
}
