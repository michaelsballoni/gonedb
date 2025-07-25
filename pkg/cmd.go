package gonedb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

type cmd_struct struct {
	Cur Node
}

func CreateCmd() cmd_struct {
	return cmd_struct{}
}

func (c *cmd_struct) ProcessCommand(db *sql.DB, cmd string) (string, error) {
	cmds := ParseCmds(cmd)
	if len(cmds) == 0 {
		return "", fmt.Errorf("no command given")
	}

	lower_cmd := strings.ToLower(cmds[0])
	output := ""
	var err error
	switch lower_cmd {
	default:
		return "", fmt.Errorf("unrecognized command: %s", lower_cmd)
	case "help":
		output += "tell - dump all info about the current node\n"
		output += "seed - seed the current node with a file system directory\n"
		output += "cd - change the current node to another node\n"
		output += "dir - list the nodes inside the current node\n"
		output += "make - create a new node inside the current node\n"
		output += "copy - make a copy of the current node in another node\n"
		output += "move - move the current node inside another node\n"
		output += "remove - remove the current node from the system\n"
		output += "rename - rename the current node\n"
		output += "setprop - set a searchable property onto the current node\n"
		output += "setpayload - set the payload string onto the current node\n"
		output += "link - link the current node to another node\n"
		output += "unlink - remove the link from the current node to another\n"
		output += "search - use properties to search for other nodes\n"
		output += "scramblelinks - randomly link nodes throughout the system\n"
	case "seed":
		if len(cmds) != 2 {
			return "", fmt.Errorf("seed command takes one parameter, the file system directory path to seed the current node with")
		}
		err = c.Seed(db, cmds[1])
	case "cd":
		if len(cmds) != 2 {
			return "", fmt.Errorf("cd command takes one parameter, the path to change to")
		}
		err = c.Cd(db, cmds[1])
	case "dir":
		if len(cmds) != 1 {
			return "", fmt.Errorf("dir takes no parameters; it lists nodes inside the current node")
		}
		var dir_strs []string
		dir_strs, err = c.Dir(db)
		if err == nil {
			output = strings.Join(dir_strs, "\n")
			if len(output) > 0 {
				output += "\n"
			}
		}
	case "make":
		if len(cmds) != 2 {
			return "", fmt.Errorf("make command takes one parameter, name of the node to create")
		}
		err = c.MakeNode(db, cmds[1])
	case "copy":
		if len(cmds) != 2 {
			return "", fmt.Errorf("copy command takes one parameter, the node to copy the current node into")
		}
		err = c.CopyToNode(db, cmds[1])
	case "move":
		if len(cmds) != 2 {
			return "", fmt.Errorf("move command takes one parameter, the node to move the current node under")
		}
		err = c.MoveToNode(db, cmds[1])
	case "remove":
		if len(cmds) != 1 {
			return "", fmt.Errorf("remove takes no parameters; it removes the current node")
		}
		err = c.RemoveNode(db)
	case "rename":
		if len(cmds) != 2 {
			return "", fmt.Errorf("rename command takes one parameter, the name to rename the current to")
		}
		err = c.Rename(db, cmds[1])
	case "setprop":
		if len(cmds) != 3 {
			return "", fmt.Errorf("setprop command takes name and value paramters; use blank strings to erase existing properties of the matching name or all of the nodes's property")
		}
		err = c.SetProp(db, cmds[1], cmds[2])
	case "setpayload":
		if len(cmds) != 2 {
			return "", fmt.Errorf("setpayload command takes the payload to set")
		}
		err = c.SetPayload(db, cmds[1])
	case "tell":
		if len(cmds) != 1 {
			return "", fmt.Errorf("tell describes the current tab and takes not parameters")
		}
		output, err = c.Tell(db)
	case "search":
		if ((len(cmds) + 1) % 2) != 0 {
			return "", fmt.Errorf("search command takes an even number of parameters, name and value pairs to search for")
		}
		output, err = c.Search(db, cmds[1:])
	case "link":
		if len(cmds) != 2 {
			return "", fmt.Errorf("link command takes one parameter, the node to link the current node to")
		}
		err = c.Link(db, cmds[1])
	case "unlink":
		if len(cmds) != 2 {
			return "", fmt.Errorf("unlink command takes one parameter, the node to unlink the current node to")
		}
		err = c.Unlink(db, cmds[1])
	case "scramblelinks":
		if len(cmds) != 1 {
			return "", fmt.Errorf("scramblelinks command takes no parameter, it works across all nodes")
		}
		created_count, scramble_err := c.ScrambleLinks(db)
		return fmt.Sprintf("%d", created_count), scramble_err
	}
	return output, err
}

// Figure out the path to the current node
func (c *cmd_struct) GetPrompt(db *sql.DB) (string, error) {
	return NodePaths.GetNodePath(db, c.Cur, "/")
}

// Seed the current node with the given file system directory,
// adding all file system entries in the directories as children of the current gonedb node
func (c *cmd_struct) Seed(db *sql.DB, dirPath string) error {
	load_err := Loader.Load(db, dirPath, c.Cur)
	return load_err
}

// Update the current node to point to a new path
func (c *cmd_struct) Cd(db *sql.DB, newPath string) error {
	if newPath == ".." {
		c.Cur.Id = c.Cur.ParentId
		return nil
	}
	nodes_path, nodes_path_err := NodePaths.GetStrNodes(db, strings.Split(newPath, "/"))
	if nodes_path_err != nil {
		return nodes_path_err
	} else if nodes_path == nil || len(*nodes_path) == 0 {
		return fmt.Errorf("node not found at path")
	} else {
		c.Cur = (*nodes_path)[len(*nodes_path)-1]
		return nil
	}
}

// List all paths to nodes that have the current nodes as their parent
func (c *cmd_struct) Dir(db *sql.DB) ([]string, error) {
	children, child_err := Nodes.GetChildren(db, c.Cur.Id)
	if child_err != nil {
		return []string{}, child_err
	}

	output := []string{}
	for _, v := range children {
		path, path_err := NodePaths.GetNodePath(db, v, "/")
		if path_err != nil {
			return []string{}, path_err
		}
		output = append(output, path)
	}

	slices.Sort(output)
	return output, nil
}

// Create a new node with the current node as its parent
func (c *cmd_struct) MakeNode(db *sql.DB, name string) error {
	name_string_id, name_err := Strings.GetId(db, name)
	if name_err != nil {
		return name_err
	}
	_, node_err := Nodes.Create(db, c.Cur.Id, name_string_id, 0)
	if node_err != nil {
		return node_err
	}
	return nil
}

// Make a copy of the current node into another node
func (c *cmd_struct) CopyToNode(db *sql.DB, path string) error {
	dest_parent_nodes, dest_err := NodePaths.GetStrNodes(db, strings.Split(path, "/"))
	if dest_err != nil {
		return dest_err
	}

	if dest_parent_nodes == nil || len(*dest_parent_nodes) == 0 {
		return fmt.Errorf("dest path does not resolve to a new parent node")
	}
	new_parent_node := (*dest_parent_nodes)[len(*dest_parent_nodes)-1]

	_, copy_err := Nodes.Copy(db, c.Cur.Id, new_parent_node.Id)
	return copy_err
}

// Move the current node to a new parent node
func (c *cmd_struct) MoveToNode(db *sql.DB, path string) error {
	dest_parent_nodes, dest_err := NodePaths.GetStrNodes(db, strings.Split(path, "/"))
	if dest_err != nil {
		return dest_err
	}

	if dest_parent_nodes == nil || len(*dest_parent_nodes) == 0 {
		return fmt.Errorf("dest path does not resolve to a new parent node")
	}
	new_parent_node := (*dest_parent_nodes)[len(*dest_parent_nodes)-1]

	move_err := Nodes.Move(db, c.Cur.Id, new_parent_node.Id)
	return move_err
}

// Remove the current node, changing the current node to its parent
func (c *cmd_struct) RemoveNode(db *sql.DB) error {
	cur_parent, cur_err := Nodes.Get(db, c.Cur.ParentId)
	if cur_err != nil {
		return cur_err
	}

	rem_err := Nodes.Remove(db, c.Cur.Id)
	if rem_err != nil {
		return rem_err
	}

	c.Cur = cur_parent
	return nil
}

// Rename the current node to a new name
func (c *cmd_struct) Rename(db *sql.DB, newName string) error {
	new_name_string_id, name_err := Strings.GetId(db, newName)
	if name_err != nil {
		return name_err
	}
	ren_err := Nodes.Rename(db, c.Cur.Id, new_name_string_id)
	if ren_err == nil {
		c.Cur.NameStringId = new_name_string_id
	}
	return ren_err
}

// Set a name-value property onto this node
func (c *cmd_struct) SetProp(db *sql.DB, name string, value string) error {
	name_string_id, name_err := Strings.GetId(db, name)
	if name_err != nil {
		return name_err
	}
	value_string_id, val_err := Strings.GetId(db, value)
	if val_err != nil {
		return val_err
	}
	var prop_err error
	if name_string_id == 0 { // delete all values in node
		prop_err = Props.Set(db, NodeItemTypeId, c.Cur.Id, -1, -1)
	} else if value_string_id == 0 { // delete all values by name
		prop_err = Props.Set(db, NodeItemTypeId, c.Cur.Id, name_string_id, -1)
	} else { // name_string_id != && value_string_id != 0 { set prop value
		prop_err = Props.Set(db, NodeItemTypeId, c.Cur.Id, name_string_id, value_string_id)
	}
	return prop_err
}

// Set the payload onto the current node
func (c *cmd_struct) SetPayload(db *sql.DB, payload string) error {
	pay_err := Nodes.SetPayload(db, c.Cur.Id, payload)
	return pay_err
}

// Describe in exquisite detail information about the current node
func (c *cmd_struct) Tell(db *sql.DB) (string, error) {
	output := fmt.Sprintf("ID: %d\n", c.Cur.Id)

	name, name_err := Strings.GetVal(db, c.Cur.NameStringId)
	if name_err != nil {
		return "", name_err
	}
	output += fmt.Sprintf("Name: %s\n", name)

	parent_node, parent_node_err := Nodes.Get(db, c.Cur.ParentId)
	if parent_node_err != nil {
		return "", parent_node_err
	}
	parent_path, parent_err := NodePaths.GetNodePath(db, parent_node, "/")
	if parent_err != nil {
		return "", parent_err
	}
	output += fmt.Sprintf("Parent: %s\n", parent_path)

	payload, payload_err := Nodes.GetPayload(db, c.Cur.Id)
	if payload_err != nil {
		return "", payload_err
	}
	output += fmt.Sprintf("Payload: %s\n", payload)

	props_map, props_err := Props.GetAll(db, NodeItemTypeId, c.Cur.Id)
	if props_err != nil {
		return "", props_err
	}
	if len(props_map) == 0 {
		output += "Properties: (none)\n"
	} else {
		props_summary, prop_summ_err := Strings.Summarize(db, props_map)
		if prop_summ_err != nil {
			return "", prop_summ_err
		}
		output += fmt.Sprintf("Properties:\n%s\n", props_summary)
	}

	out_links, out_link_err := Links.GetOutLinks(db, c.Cur.Id)
	if out_link_err != nil {
		return "", out_link_err
	}
	if len(out_links) == 0 {
		output += "Out Links: (none)\n"
	} else {
		output += fmt.Sprintf("Out Links: (%d)\n", len(out_links))
		for _, link := range out_links {
			to_node, to_node_err := Nodes.Get(db, link.ToNodeId)
			if to_node_err != nil {
				return "", to_node_err
			}
			to_node_path, to_none_path_err := NodePaths.GetNodePath(db, to_node, "/")
			if to_none_path_err != nil {
				return "", to_none_path_err
			}
			output += to_node_path + "\n"
		}
	}

	in_links, in_link_err := Links.GetToLinks(db, c.Cur.Id)
	if in_link_err != nil {
		return "", in_link_err
	}
	if len(in_links) == 0 {
		output += "In Links: (none)\n"
	} else {
		output += fmt.Sprintf("In Links: (%d)\n", len(in_links))
		for _, link := range in_links {
			in_node, in_node_err := Nodes.Get(db, link.FromNodeId)
			if in_node_err != nil {
				return "", in_node_err
			}
			in_node_path, in_none_path_err := NodePaths.GetNodePath(db, in_node, "/")
			if in_none_path_err != nil {
				return "", in_none_path_err
			}
			output += in_node_path + "\n"
		}
	}

	return output, nil
}

// Get search results for name-value search criteria
func (c *cmd_struct) Search(db *sql.DB, name_values []string) (string, error) {
	if len(name_values) < 2 {
		return "", fmt.Errorf("pass in name / value pairs to search properties with")
	}

	if (len(name_values) % 2) != 0 {
		return "", fmt.Errorf("pass in evenly matched name / value pairs to search with")
	}
	var query SearchQuery
	for i := 0; i < len(name_values); i += 2 {
		var criteria SearchCriteria
		var name_string_err error
		criteria.NameStringId, name_string_err = Strings.GetId(db, name_values[i])
		if name_string_err != nil {
			return "", name_string_err
		}
		criteria.ValueString = name_values[i+1]
		query.Criteria = append(query.Criteria, criteria)
	}

	node_results, results_err := Search.FindNodes(db, &query)
	if results_err != nil {
		return "", results_err
	}

	var output strings.Builder
	for _, cur_node := range node_results {
		path, path_err := NodePaths.GetNodePath(db, cur_node, "/")
		if path_err != nil {
			return "", path_err
		}
		output.WriteString(path)
		output.WriteString("\n")
	}
	return output.String(), nil
}

// Create a link between the current node and another node
func (c *cmd_struct) Link(db *sql.DB, toPath string) error {
	to_nodes, to_node_err := NodePaths.GetStrNodes(db, strings.Split(toPath, "/"))
	if to_node_err != nil {
		return to_node_err
	}
	to_node := (*to_nodes)[len(*to_nodes)-1]
	Links.Create(db, c.Cur.Id, to_node.Id, NodeItemTypeId)
	return nil
}

// Remove a link between this and another node
func (c *cmd_struct) Unlink(db *sql.DB, toPath string) error {
	to_nodes, to_node_err := NodePaths.GetStrNodes(db, strings.Split(toPath, "/"))
	if to_node_err != nil {
		return to_node_err
	}
	to_node := (*to_nodes)[len(*to_nodes)-1]
	Links.RemoveFromTo(db, c.Cur.Id, to_node.Id, NodeItemTypeId)
	return nil
}

// Scramble links among all nodes
func (c *cmd_struct) ScrambleLinks(db *sql.DB) (int64, error) {
	// get all node ids
	rows, err := db.Query("SELECT id FROM nodes")
	if err != nil {
		return -1, err
	}
	ids := []int64{}
	var cur_id int64
	for rows.Next() {
		err := rows.Scan(&cur_id)
		if err != nil {
			return -1, err
		}
		ids = append(ids, cur_id)
	}

	// link half
	ids_len := len(ids)
	half := ids_len / 2
	if half < 2 {
		return -1, fmt.Errorf("need at least two nodes to scramble links")
	}

	// walk them in random order linking them to another at a random order
	seen_ids := map[int64]bool{}
	created := int64(0)
	for i := 1; i <= half; i++ {
		from := ids[rand.Int()%ids_len]
		if seen_ids[from] {
			continue
		}
		seen_ids[from] = true

		to := ids[rand.Int()%ids_len]
		if from == to {
			continue
		}

		_, err = Links.Create(db, from, to, 0)
		if err != nil {
			return -1, err
		}

		created += 1
	}
	return created, nil
}

// Given a command-line, handle quoted or unquoted strings as separate parameters
func ParseCmds(cmd string) []string {
	output := []string{}
	collector := ""
	in_quote := false
	for _, c := range cmd {
		if c == '"' {
			if !in_quote {
				collector = strings.TrimSpace(collector)
			}

			if in_quote || len(collector) > 0 {
				output = append(output, collector)
				collector = ""
			}

			in_quote = !in_quote
			continue
		}

		if !in_quote && c == ' ' {
			collector = strings.TrimSpace(collector)
			if len(collector) > 0 {
				output = append(output, collector)
				collector = ""
			}
			continue
		}

		collector += string(c)
	}

	collector = strings.TrimSpace(collector)
	if len(collector) > 0 {
		output = append(output, collector)
		collector = ""
	}

	return output
}
