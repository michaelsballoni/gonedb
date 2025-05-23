package gonedb

import (
	"database/sql"
)

type node_paths struct{}

// The API you interact with
var NodePaths node_paths

func (np *node_paths) GetNodes(db *sql.DB, nodeId int64) ([]Node, error) {
	// FORNOW
	return []Node{}, nil
}

func (np *node_paths) GetStr(db *sql.DB, nodeId int64) (string, error) {
	// FORNOW
	return "", nil
}

func (np *node_paths) GetStrNodes(db *sql.DB, nodeId int64) ([]Node, error) {
	// FORNOW
	return []Node{}, nil
}
