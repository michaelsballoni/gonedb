package gonedb

import (
	"database/sql"
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
