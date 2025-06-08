package gonedb

import (
	"database/sql"
	"strconv"
	"strings"
)

type search struct{}

var Search search

func (s *search) FindNodes(db *sql.DB, searchQuery *SearchQuery) ([]Node, error) {
	if len(searchQuery.Criteria) == 0 {
		return []Node{}, nil
	}

	var find_params_obj find_params
	find_params_obj.Query = searchQuery
	find_params_obj.ItemType = NodeItemTypeId
	find_params_obj.ItemTable = "nodes"
	find_params_obj.Columns = "Items.id, parent_id, name_string_id, type_string_id"

	sql_params := map[string]variant{}
	sql, err := get_find_sql(db, &find_params_obj, sql_params)
	//fmt.Printf("PRE: sql: %s - err: %v\n", sql, err)
	if err != nil {
		return []Node{}, err
	}

	for k, v := range sql_params {
		sql = strings.ReplaceAll(sql, k, varToSql(v))
	}
	//fmt.Printf("PRM: sql: %s - err: %v\n", sql, err)

	output := []Node{}
	var cur_node Node
	query, query_err := db.Query(sql)
	if query_err != nil {
		return []Node{}, query_err
	}
	for query.Next() {
		scan_err := query.Scan(&cur_node.Id, &cur_node.ParentId, &cur_node.NameStringId, &cur_node.TypeStringId)
		if scan_err != nil {
			return []Node{}, scan_err
		}
		output = append(output, cur_node)
	}
	return output, nil
}

func (s *search) FindLinks(db *sql.DB, searchQuery *SearchQuery) ([]Link, error) {
	if len(searchQuery.Criteria) == 0 {
		return []Link{}, nil
	}

	var find_params_obj find_params
	find_params_obj.Query = searchQuery
	find_params_obj.ItemType = LinkItemTypeId
	find_params_obj.ItemTable = "links"
	find_params_obj.Columns = "Items.id, from_node_id, to_node_id, type_string_id"

	sql_params := map[string]variant{}
	sql, err := get_find_sql(db, &find_params_obj, sql_params)
	if err != nil {
		return []Link{}, err
	}

	for k, v := range sql_params {
		sql = strings.ReplaceAll(sql, k, varToSql(v))
	}

	output := []Link{}
	var cur_link Link
	query, query_err := db.Query(sql)
	if query_err != nil {
		return []Link{}, query_err
	}
	for query.Next() {
		scan_err := query.Scan(&cur_link.Id, &cur_link.FromNodeId, &cur_link.ToNodeId, &cur_link.TypeStringId)
		if scan_err != nil {
			return []Link{}, scan_err
		}
		output = append(output, cur_link)
	}
	return output, nil
}

type SearchQuery struct {
	Criteria       []SearchCriteria
	OrderBy        string
	OrderAscensing bool
	Limit          int64
}

func get_asc_str(sq *SearchQuery) string {
	if sq.OrderAscensing {
		return "ASC"
	} else {
		return "DESC"
	}
}

type SearchCriteria struct {
	NameStringId int64
	ValueString  string
	UseLike      bool
}

type find_params struct {
	Query     *SearchQuery
	ItemTable string
	ItemType  int64
	Columns   string
}

type variant struct {
	IsNum  bool
	NumVal int64
	StrVal string
}

func createNumVar(num int64) variant {
	var output variant
	output.IsNum = true
	output.NumVal = num
	return output
}
func createStrVar(str string) variant {
	var output variant
	output.IsNum = false
	output.StrVal = str
	return output
}
func varToSql(v variant) string {
	if v.IsNum {
		return strconv.FormatInt(v.NumVal, 10)
	} else {
		return "'" + strings.ReplaceAll(v.StrVal, "'", "''") + "'"
	}
}

func get_find_sql(db *sql.DB, findParams *find_params, sqlParams map[string]variant) (string, error) {
	var name_string_id, payload_string_id, type_string_id, order_by_string_id int64
	var err error
	// FORNOW parent_string_id, path_string_id,
	/*
		parent_string_id, err = Strings.GetId(db, "parent")
		if err != nil {
			return "", err
		}

		path_string_id, err = Strings.GetId(db, "path")
		if err != nil {
			return "", err
		}
	*/
	name_string_id, err = Strings.GetId(db, "name")
	if err != nil {
		return "", err
	}

	payload_string_id, err = Strings.GetId(db, "payload")
	if err != nil {
		return "", err
	}

	type_string_id, err = Strings.GetId(db, "type")
	if err != nil {
		return "", err
	}

	sqlParams["@node_item_type_id"] = createNumVar(findParams.ItemType)

	if len(findParams.Query.OrderBy) > 0 {
		order_by_string_id, err = Strings.GetId(db, findParams.Query.OrderBy)
		if err != nil {
			return "", err
		}
		sqlParams["@order_by_string_id"] = createNumVar(order_by_string_id)
	}

	sql := "SELECT " + findParams.Columns + " FROM " + findParams.ItemTable + " AS Items "
	if len(findParams.Query.OrderBy) > 0 {
		sql += "JOIN props AS ItemProps ON ItemProps.itemid = Items.id JOIN strings AS ItemStrings ON ItemStrings.id = ItemProps.valstrid "
	}

	where := ""
	for param_num, crit := range findParams.Query.Criteria {
		if len(where) > 0 {
			where += "\nAND "
		}

		param_num_str := strconv.Itoa(param_num + 1)

		if crit.NameStringId == name_string_id { // searching by name
			new_sql := "Items.id IN (SELECT InnerNodes.id FROM nodes InnerNodes JOIN strings NameStrings ON NameStrings.id = InnerNodes.name_string_id WHERE "
			if crit.UseLike {
				sqlParams["@valstr"+param_num_str] = createStrVar(crit.ValueString)
				new_sql += "NameStrings.val LIKE @valstr" + param_num_str
			} else {
				val_string_id, val_string_err := Strings.GetId(db, crit.ValueString)
				if val_string_err != nil {
					return "", val_string_err
				}
				sqlParams["@valstrid"+param_num_str] = createNumVar(val_string_id)
				new_sql += "NameStrings.id = @valstrid" + param_num_str
			}
			new_sql += ")"

			where += new_sql
		} else if crit.NameStringId == type_string_id { // searching by type
			val_string_id, val_string_err := Strings.GetId(db, crit.ValueString)
			if val_string_err != nil {
				return "", val_string_err
			}
			sqlParams["@valstrid"+param_num_str] = createNumVar(val_string_id)
			where += "type_string_id = @valstrid" + param_num_str
		} else if crit.NameStringId == payload_string_id { // search by payload
			new_sql := ""
			sqlParams["@valstr"+param_num_str] = createStrVar(crit.ValueString)
			if crit.UseLike {
				new_sql += "payload LIKE @valstr" + param_num_str
			} else {
				new_sql += "payload = @valstr" + param_num_str
			}
			where += new_sql
			/*
				FORNOW - Add back path lookup stuff to Nodes to support search like this
				} else if (crit.m_nameStringId == parent_string_id) { // search directly within a parent node
					Nodes.GetParentsNodeIds()
					Nodes.GET
					auto parent_path_opt = nodes::get_path_nodes(db, crit.m_valueString);
					if (parent_path_opt.has_value())
					{
						int64_t parent_id = parent_path_opt.value().back().id;
						where += "Items.parent_id = " + std::to_wstring(parent_id);
					}
					else
						where += "1 = 0"; // no path, no results
				} else if (crit.m_nameStringId == path_string_id) { // search deeply within a parent node
					auto child_like_opt = nodes::get_path_to_parent_like(db, crit.m_valueString);
					if (child_like_opt.has_value())
					{
						sqlParams["@valstr" + param_num_str] = child_like_opt.value();
						where += "Items.parents LIKE @valstr" + param_num_str;
					}
					else
						where += "1 = 0"; // no path, no results
			*/
		} else {
			sqlParams["@namestrid"+param_num_str] = createNumVar(crit.NameStringId)

			new_sql := "Items.id IN (SELECT itemid FROM props WHERE itemtypeid = @node_item_type_id AND namestrid = @namestrid" + param_num_str
			if crit.UseLike {
				sqlParams["@valstr"+param_num_str] = createStrVar(crit.ValueString)
				new_sql += " AND valstrid IN (SELECT id FROM strings WHERE val LIKE @valstr" + param_num_str + ")"
			} else {
				val_string_id, val_string_err := Strings.GetId(db, crit.ValueString)
				if val_string_err != nil {
					return "", val_string_err
				}
				sqlParams["@valstrid"+param_num_str] = createNumVar(val_string_id)
				new_sql += " AND valstrid = @valstrid" + param_num_str
			}
			new_sql += ")"

			where += new_sql
		}
	}
	sql += "WHERE " + where

	if len(findParams.Query.OrderBy) > 0 {
		sql += "\nAND ItemProps.itemtypeid = @node_item_type_id AND ItemProps.namestrid = @order_by_string_id"
		sql += "\nORDER BY ItemStrings.val " + get_asc_str(findParams.Query)
	}

	if findParams.Query.Limit > 0 {
		sql += "\nLIMIT " + strconv.FormatInt(findParams.Query.Limit, 10)
	}

	return sql, nil
}
