package gonedb

import "database/sql"

func Setup(db *sql.DB) {
	// strings
	db.Exec("CREATE TABLE strings (id INTEGER PRIMARY KEY, val STRING UNIQUE NOT NULL)")
	db.Exec("INSERT INTO strings (id, val) VALUES (0, '')")

	// nodes
	db.Exec(`
		CREATE TABLE nodes
		(
		id INTEGER PRIMARY KEY,
		parent_id INTEGER NOT NULL,
		type_string_id INTEGER NOT NULL,
		name_string_id INTEGER NOT NULL,
		parents STRING NOT NULL DEFAULT '',
		payload STRING NOT NULL DEFAULT ''
		)`)
	db.Exec("CREATE UNIQUE INDEX node_parents ON nodes (parent_id, id)")
	db.Exec("CREATE UNIQUE INDEX node_names ON nodes (parent_id, name_string_id)")
	db.Exec("CREATE UNIQUE INDEX node_parent_strs ON nodes (parents, id)")
	db.Exec("CREATE INDEX node_payloads ON nodes (payload, id)")
	db.Exec("INSERT INTO nodes (id, parent_id, type_string_id, name_string_id) VALUES (0, 0, 0, 0)")

	// links
	db.Exec(`
		CREATE TABLE links
		(
		id INTEGER PRIMARY KEY,
		from_node_id INTEGER NOT NULL,
		to_node_id INTEGER NOT NULL,
		type_string_id INTEGER NOT NULL,
		payload STRING NOT NULL DEFAULT ''
		)`)
	db.Exec("CREATE UNIQUE INDEX link_from ON links (from_node_id, to_node_id, type_string_id)")
	db.Exec("CREATE UNIQUE INDEX link_to ON links (to_node_id, from_node_id, type_string_id)")
	db.Exec("CREATE INDEX link_payloads ON links (payload, id)")
	db.Exec("INSERT INTO links (id, from_node_id, to_node_id, type_string_id) VALUES (0, 0, 0, 0)")

	// props
	db.Exec("CREATE TABLE props (id INTEGER PRIMARY KEY, itemtypeid INTEGER, itemid INTEGER, namestrid INTEGER, valstrid INTEGER)")
	db.Exec("CREATE UNIQUE INDEX item_props ON props (itemtypeid, itemid, namestrid)")
	db.Exec("CREATE INDEX prop_vals ON props (valstrid, namestrid, itemtypeid, itemid)")
	db.Exec("INSERT INTO props (id, itemtypeid, itemid, namestrid, valstrid) VALUES (0, 0, 0, 0, 0)")
}
