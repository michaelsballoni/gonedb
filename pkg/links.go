package gonedb

import (
	"database/sql"
	"fmt"
)

type links struct{}

// The API you interact with
var Links links

type Link struct {
	Id           int64
	FromNodeId   int64
	ToNodeId     int64
	TypeStringId int64
}

// Create a link
func (l *links) Create(db *sql.DB, fromNodeId int64, toNodeId int64, typeStringId int64) (Link, error) {
	var new_id int64
	row := db.QueryRow("INSERT INTO links (from_node_id, to_node_id, type_string_id) VALUES (?, ?, ?) RETURNING id",
		fromNodeId,
		toNodeId,
		typeStringId)
	err := row.Scan(&new_id)
	if err != nil {
		row := db.QueryRow("SELECT id FROM links where from_node_id = ?, to_node_id = ?, type_string_id = ?",
			fromNodeId, toNodeId, typeStringId)
		err := row.Scan(&new_id)
		if err != nil {
			return Link{}, err
		}
	}
	return Link{Id: new_id, FromNodeId: fromNodeId, ToNodeId: toNodeId, TypeStringId: typeStringId}, nil
}

// Remove a link by ID
func (l *links) Remove(db *sql.DB, linkId int64) error {
	result, err := db.Exec("DELETE FROM links WHERE id = ?", linkId)
	if err != nil {
		return err
	} else {
		_, aff_err := result.RowsAffected()
		if aff_err != nil {
			return aff_err
		} else {
			return nil
		}
	}
}

// Remove a link by to/from IDs
func (l *links) RemoveFromTo(db *sql.DB, fromNodeId int64, toNodeId int64, typeStringId int64) error {
	result, err := db.Exec("DELETE FROM links WHERE from_node_id = ? AND to_node_id = ? AND type_string_id = ?", fromNodeId, toNodeId, typeStringId)
	if err != nil {
		return err
	} else {
		_, aff_err := result.RowsAffected()
		if aff_err != nil {
			return aff_err
		} else {
			return nil
		}
	}
}

// Get a link by ID
func (l *links) Get(db *sql.DB, linkId int64) (Link, error) {
	var output Link
	output.Id = linkId
	row := db.QueryRow("SELECT from_node_id, to_node_id, type_string_id FROM links WHERE id = ?", linkId)
	err := row.Scan(&output.FromNodeId, &output.ToNodeId, &output.TypeStringId)
	if err != nil {
		return Link{}, err
	} else {
		return output, nil
	}
}

// Get the payload of a link
func (l *links) GetPayload(db *sql.DB, linkId int64) (string, error) {
	var output string
	row := db.QueryRow("SELECT payload FROM links WHERE id = ?", linkId)
	err := row.Scan(&output)
	return output, err
}

// Set the payload of a link
func (l *links) SetPayload(db *sql.DB, linkId int64, payload string) error {
	if linkId == 0 {
		return fmt.Errorf("cannot set payload on null link")
	}

	result, err := db.Exec("UPDATE links SET payload = ? WHERE id = ?", payload, linkId)
	if err != nil {
		return err
	}
	affected, affected_err := result.RowsAffected()
	if affected_err != nil {
		return affected_err
	}
	if affected != 1 {
		return fmt.Errorf("row not affected")
	} else {
		return nil
	}
}

func (l *links) GetOutLinks(db *sql.DB, nodeId int64) ([]Link, error) {
	query, query_err := db.Query("SELECT id, from_node_id, to_node_id, type_string_id FROM links WHERE from_node_id = ?", nodeId)
	if query_err != nil {
		return []Link{}, query_err
	}
	output := []Link{}
	var cur_link Link
	for query.Next() {
		scan_err := query.Scan(&cur_link.Id, &cur_link.FromNodeId, &cur_link.ToNodeId, &cur_link.TypeStringId)
		if scan_err != nil {
			return []Link{}, scan_err
		}
		output = append(output, cur_link)
	}
	return output, nil
}

func (l *links) GetToLinks(db *sql.DB, nodeId int64) ([]Link, error) {
	query, query_err := db.Query("SELECT id, from_node_id, to_node_id, type_string_id FROM links WHERE to_node_id = ?", nodeId)
	if query_err != nil {
		return []Link{}, query_err
	}
	output := []Link{}
	var cur_link Link
	for query.Next() {
		scan_err := query.Scan(&cur_link.Id, &cur_link.FromNodeId, &cur_link.ToNodeId, &cur_link.TypeStringId)
		if scan_err != nil {
			return []Link{}, scan_err
		}
		output = append(output, cur_link)
	}
	return output, nil
}
