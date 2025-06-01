package gonedb

import (
	"database/sql"
	"sort"
	"strings"
)

type props struct{}

// The API you interact with
var Props props

var NodeItemTypeId int64 = 1
var LinkItemTypeId int64 = 2

// Get a named property onto an item
func (p *props) Set(db *sql.DB, itemTypeId int64, itemId int64, nameStringId int64, valueStringId int64) error { // use < 0 to delete
	if valueStringId >= 0 {
		_, err :=
			db.Exec(`INSERT INTO props (itemtypeid, itemid, namestrid, valstrid) VALUES(?, ?, ?, ?) ON CONFLICT(itemtypeid, itemid, namestrid) DO UPDATE SET valstrid = ?`,
				itemTypeId,
				itemId,
				nameStringId,
				valueStringId,
				valueStringId)
		return err
	} else { // delete the value for the name
		if nameStringId >= 0 {
			_, err := db.Exec("DELETE FROM props WHERE itemtypeid = ? AND itemid = ? AND namestrid = ?", itemTypeId, itemId, nameStringId)
			return err
		} else { // delete all values for the node
			_, err := db.Exec("DELETE FROM props WHERE itemtypeid = ? AND itemid = ?", itemTypeId, itemId)
			return err
		}
	}
}

// Get the value of one named property from an item
func (p *props) Get(db *sql.DB, itemTypeId int64, itemId int64, nameStringId int64) (int64, error) {
	result := db.QueryRow("SELECT valstrid FROM props WHERE itemtypeid = ? AND itemid = ? AND namestrid = ?", itemTypeId, itemId, nameStringId)
	var value_string_id int64
	scan_err := result.Scan(&value_string_id)
	return value_string_id, scan_err
}

// Get all name-value property strings from an item
func (p *props) GetAll(db *sql.DB, itemTypeId int64, itemId int64) (map[int64]int64, error) {
	result, query_err := db.Query("SELECT namestrid, valstrid FROM props WHERE itemtypeid = ? AND itemid = ?", itemTypeId, itemId)
	if query_err != nil {
		return map[int64]int64{}, query_err
	}

	output := map[int64]int64{}
	var name_string_id, value_string_id int64
	for result.Next() {
		scan_err := result.Scan(&name_string_id, &value_string_id)
		if scan_err != nil {
			return map[int64]int64{}, scan_err
		}
		output[name_string_id] = value_string_id
	}
	return output, nil
}

// Take a stringID-to-stringID map and return a stringValue-to-stringValue map
func (p *props) Fill(db *sql.DB, nameValueStrIds map[int64]int64) (map[string]string, error) {
	output := map[string]string{}
	var name_str, val_str string
	var cur_err error
	for name_string_id, value_string_id := range nameValueStrIds {
		name_str, cur_err = Strings.GetVal(db, name_string_id)
		if cur_err != nil {
			return map[string]string{}, cur_err
		}
		val_str, cur_err = Strings.GetVal(db, value_string_id)
		if cur_err != nil {
			return map[string]string{}, cur_err
		}
		output[name_str] = val_str
	}
	return output, nil
}

// Take a stringID-to-stringID map and return a line-delimited string summarizing the names and values
func (p *props) Summarize(db *sql.DB, nameValueStrIds map[int64]int64) (string, error) {
	string_map, fill_err := p.Fill(db, nameValueStrIds)
	if fill_err != nil {
		return "", fill_err
	}

	names := make([]string, 0, len(string_map))
	for k, _ := range string_map {
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
