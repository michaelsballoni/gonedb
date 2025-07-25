package gonedb

import (
	"database/sql"
	"fmt"
)

type clouds struct{}

// The API you interact with
var Clouds clouds

// The struct passed around the API with the core IDs of a node
type Cloud struct {
	m_seedNodeId int64
	m_tableName  string
}

// Create a null clould given just the seed node ID
// The internal table for this clould is seeded using this ID
func (c *clouds) GetCloud(seedNodeId int64) Cloud {
	return Cloud{m_seedNodeId: seedNodeId, m_tableName: fmt.Sprintf("cloudlinks_%d", seedNodeId)}
}

// Initialize the database table for this cloud
// Call again to drop and recreate the table
func (c *Cloud) Init(db *sql.DB) error {
	_, db_err := db.Exec("DROP TABLE IF EXISTS " + c.m_tableName)
	if db_err != nil {
		return nil
	}

	_, db_err = db.Exec(
		"CREATE TABLE " + c.m_tableName +
			`(
		gen INTEGER NOT NULL,
		id INTEGER UNIQUE NOT NULL,
		from_node_id INTEGER NOT NULL,
		to_node_id INTEGER NOT NULL,
		type_string_id INTEGER NOT NULL)`)
	if db_err != nil {
		return nil
	}

	_, db_err = db.Exec(fmt.Sprintf("CREATE INDEX link_from_%d ON %s (from_node_id, to_node_id)", c.m_seedNodeId, c.m_tableName))
	if db_err != nil {
		return nil
	}

	_, db_err = db.Exec(fmt.Sprintf("CREATE INDEX link_to_%d ON %s (to_node_id, from_node_id)", c.m_seedNodeId, c.m_tableName))
	if db_err != nil {
		return nil
	}

	_, db_err = db.Exec(fmt.Sprintf("CREATE INDEX gen_linkid_%d ON %s (gen, id)", c.m_seedNodeId, c.m_tableName))
	if db_err != nil {
		return nil
	}

	return nil
}

// See the table with the links to and from the need node ID
func (c *Cloud) Seed(db *sql.DB) (int64, error) {
	sql :=
		"INSERT INTO " + c.m_tableName + " (gen, id, from_node_id, to_node_id, type_string_id)" +
			"SELECT 0, id, from_node_id, to_node_id, type_string_id " + // gen 0
			"FROM links " +
			"WHERE from_node_id = ? OR to_node_id = ?"
	res, err := db.Exec(sql, c.m_seedNodeId, c.m_seedNodeId)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	return affected, err
}

/*
int64_t expand();
int max_generation() const;
func (c *Cloud) GetLinks(minGeneration int64, maxGeneration int64): ([]Link, error) {
}
*/
