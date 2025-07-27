package gonedb

import (
	"database/sql"
	"fmt"
	"regexp"
)

type clouds struct{}

// The API you interact with
var Clouds clouds

// The struct passed around the API with the core IDs of a node
type Cloud struct {
	m_cloudName  string
	m_seedNodeId int64
	m_tableName  string
}

// Create a null clould given just the seed node ID
// The internal table for this clould is seeded using this ID and name
// The name much match the RE [\w]+$
func (c *clouds) GetCloud(cloudName string, seedNodeId int64) (Cloud, error) {
	name_ok, name_err := regexp.MatchString(`^[\w]+$`, cloudName)
	if name_err != nil {
		return Cloud{}, name_err
	}
	if !name_ok {
		return Cloud{}, fmt.Errorf("invalid cloud name, must match ^[/w]+$")
	}
	return Cloud{m_seedNodeId: seedNodeId, m_tableName: fmt.Sprintf("cloudlinks_%d_%s", seedNodeId, cloudName)}, nil
}

// Drop the database out of the database
func (c *Cloud) Drop(db *sql.DB) error {
	_, db_err := db.Exec("DROP TABLE IF EXISTS " + c.m_tableName)
	if db_err != nil {
		return nil
	}
	return nil
}

// Initialize the database table for this cloud
// Call again to drop and recreate the table
func (c *Cloud) Init(db *sql.DB) error {
	_, db_err := db.Exec("DROP TABLE IF EXISTS " + c.m_tableName)
	if db_err != nil {
		return nil
	}

	_, db_err = db.Exec(
		"CREATE TABLE " + c.m_tableName + `(
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

	_, db_err = db.Exec(fmt.Sprintf("CREATE INDEX link_id_%d ON %s (id)", c.m_seedNodeId, c.m_tableName))
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

// Get the links in a particular generation range
func (c *Cloud) GetLinks(db *sql.DB, minGeneration int64, maxGeneration int64) ([]Link, error) {
	output := []Link{}
	sql := "SELECT id, from_node_id, to_node_id, type_string_id FROM " + c.m_tableName + " WHERE gen >= ? AND gen <= ?"
	rows, err := db.Query(sql, minGeneration, maxGeneration)
	if err != nil {
		return []Link{}, err
	}
	var cur_link Link
	for rows.Next() {
		err = rows.Scan(&cur_link.Id, &cur_link.FromNodeId, &cur_link.ToNodeId, &cur_link.TypeStringId)
		if err != nil {
			return []Link{}, err
		}
		output = append(output, cur_link)
	}
	return output, nil
}

// Get the max generation for this cloud
func (c *Cloud) GaxMaxGeneration(db *sql.DB) (int64, error) {
	row := db.QueryRow("SELECT MAX(gen) FROM " + c.m_tableName)
	var max_gen int64
	err := row.Scan(&max_gen)
	return max_gen, err
}

// Expand this cloud out when generation
// Returns number of links added, or error
func (c *Cloud) Expand(db *sql.DB) (int64, error) {
	// settle on the current gen
	max_gen, err := c.GaxMaxGeneration(db)
	if err != nil {
		return -1, err
	}
	max_gen += 1

	// add new stuff
	sql :=
		"INSERT INTO " + c.m_tableName + " (gen, id, from_node_id, to_node_id, type_string_id)" +
			"SELECT " + fmt.Sprintf("%d", max_gen) + ", id, from_node_id, to_node_id, type_string_id " +
			"FROM links " +
			"WHERE " +
			"(" +
			"from_node_id IN (SELECT from_node_id FROM " + c.m_tableName + ") " +
			"OR " +
			"from_node_id IN (SELECT to_node_id FROM " + c.m_tableName + ") " +
			"OR " +
			"to_node_id IN (SELECT to_node_id FROM " + c.m_tableName + ") " +
			"OR " +
			"to_node_id IN (SELECT from_node_id FROM " + c.m_tableName + ") " +
			") " +
			"AND id NOT IN (SELECT id FROM " + c.m_tableName + ")"
	exec_res, exec_err := db.Exec(sql)
	if exec_err != nil {
		return -1, exec_err
	}
	affected, affected_err := exec_res.RowsAffected()
	return affected, affected_err
}
