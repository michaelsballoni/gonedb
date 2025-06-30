package gonedb

import (
	"database/sql"
	"fmt"
	"strings"
)

type cmd_struct struct {
	Cur Node
}

func CreateCmd() cmd_struct {
	return cmd_struct{}
}

// Mount the given file system directory into the current node
func (c *cmd_struct) Mount(db *sql.DB, dirPath string) error {
	load_err := Loader.Load(db, dirPath, c.Cur)
	return load_err
	/* FORNOW
	// normalize the input path
	dirPath = strings.TrimSpace(dirPath)
	if len(dirPath) == 0 {
		dirPath = "."
	}
	if dirPath == "." {
		return nil
	}
	if dirPath == "/" {
		c.Cur = Node{}
		return nil
	}

	// get the path to the new node
	var new_node_path string = ""
	if dirPath[0] == '/' {
		node_strs, strs_err := NodePaths.GetStrs(db, c.Cur)
		if strs_err != nil {
			return strs_err
		}

		new_node_path = strings.Join(node_strs, "/")
		if len(new_node_path) == 0 || new_node_path[len(new_node_path)-1] != '/' {
			new_node_path += "/"
		}
		new_node_path += dirPath
	} else {
		new_node_path = dirPath
	}

	// split the new node path
	path_parts := strings.Split(new_node_path, "/")
	path_parts_out := make([]string, 0, len(path_parts))
	for _, v := range path_parts {
		v_trim := strings.TrimSpace(v)
		if v_trim != "" {
			path_parts_out = append(path_parts_out, v_trim)
		}
	}

	// get the nodes for the new path
	nodes_path, nodes_path_err := NodePaths.GetStrNodes(db, path_parts_out)
	if nodes_path_err != nil {
		return nodes_path_err
	} else if nodes_path == nil || len(*nodes_path) == 0 {
		return fmt.Errorf("node not found at path")
	}

	c.Cur = (*nodes_path)[len(*nodes_path)-1]
	*/
}

func (c *cmd_struct) Cd(db *sql.DB, newPath string) error {
	nodes_path, nodes_path_err := NodePaths.GetStrNodes(db, strings.Split(newPath, "/"))
	if nodes_path_err != nil {
		return nodes_path_err
	} else if nodes_path == nil || len(*nodes_path) == 0 {
		return fmt.Errorf("node not found at path")
	}

	c.Cur = (*nodes_path)[len(*nodes_path)-1]
	return nil
}

/* FORNOW
std::vector<std::wstring> cmd::dir()
{
	std::vector<std::wstring> paths;

	for (auto child : nodes::get_children(m_db, m_cur.id))
		paths.emplace_back(nodes::get_path_str(m_db, child));
	std::sort(paths.begin(), paths.end());
	return paths;
}

void cmd::mknode(const std::wstring& newNodeName)
{
	nodes::create(m_db, m_cur.id, strings::get_id(m_db, newNodeName));
}

void cmd::copy(const std::wstring& newParentNode)
{
	nodes::copy(m_db, m_cur.id, get_node_from_path(newParentNode).id);
}

void cmd::move(const std::wstring& newParentNode)
{
	nodes::move(m_db, m_cur.id, get_node_from_path(newParentNode).id);
}

void cmd::remove()
{
	int64_t orig_id = m_cur.id;
	auto parent = nodes::get(m_db, m_cur.parentId);
	nodes::remove(m_db, orig_id);
	m_cur = parent;
}

void cmd::rename(const std::wstring& newName)
{
	int64_t new_name_string_id = strings::get_id(m_db, newName);
	nodes::rename(m_db, m_cur.id, new_name_string_id);
	m_cur.nameStringId = new_name_string_id;
}

void cmd::set_prop(const std::vector<std::wstring>& cmds)
{
	if (cmds.size() < 2)
		throw nldberr("Specify the name of the property to set");
	else if (cmds.size() > 3)
		throw nldberr("Specify the name and value of the property to set");
	else if (cmds.size() == 2)
		props::set(m_db, m_nodeItemTypeId, m_cur.id, strings::get_id(m_db, cmds[1]), -1);
	else
		props::set(m_db, m_nodeItemTypeId, m_cur.id, strings::get_id(m_db, cmds[1]), strings::get_id(m_db, cmds[2]));
}

void cmd::set_payload(const std::wstring& payload)
{
	nodes::set_payload(m_db, m_cur.id, payload);
}

std::wstring cmd::tell()
{
	std::wstringstream stream;

	stream << L"ID:      " << m_cur.id << L"\n";
	stream << L"Name:    " << strings::get_val(m_db, m_cur.nameStringId) << L"\n";
	stream << L"Parent:  " << nodes::get_path_str(m_db, m_cur) << L"\n";
	stream << L"Payload: " << nodes::get_payload(m_db, m_cur.id) << L"\n";

	auto prop_string_ids = props::get(m_db, m_nodeItemTypeId, m_cur.id);
	if (!prop_string_ids.empty())
	{
		stream << L"Properties:\n" << props::summarize(m_db, prop_string_ids) << L"\n";
	}
	else
		stream << L"Properties: (none)" << L"\n";

	auto out_links = links::get_out_links(m_db, m_cur.id);
	if (!out_links.empty())
	{
		stream << L"Out Links:" << L"\n";
		for (const auto& out_link : out_links)
			stream << nodes::get_path_str(m_db, nodes::get(m_db, out_link.toNodeId)) << L"\n";
	}
	else
		stream << L"Out Links: (none)" << L"\n";

	auto in_links = links::get_in_links(m_db, m_cur.id);
	if (!in_links.empty())
	{
		stream << L"In Links:" << L"\n";
		for (const auto& in_link : in_links)
			stream << nodes::get_path_str(m_db, nodes::get(m_db, in_link.fromNodeId)) << L"\n";
	}
	else
		stream << L"In Links:  (none)" << L"\n";

	return stream.str();
}

std::wstring cmd::search(const std::vector<std::wstring>& cmd)
{
	if (cmd.size() < 3)
		throw nldberr("Pass in name / value pairs to search properties with");

	if (((int)cmd.size() - 1) % 2)
		throw nldberr("Pass in evenly matched name / value pairs to search with");

	search_query query;
	for (size_t s = 1; s + 1 < cmd.size(); s += 2)
	{
		query.m_criteria.push_back
		(
			search_criteria(strings::get_id(m_db, cmd[s]), cmd[s + 1])
		);
	}
	std::wstring output;
	for (const auto& node : search::find_nodes(m_db, query))
	{
		if (!output.empty())
			output += '\n';
		output += nodes::get_path_str(m_db, node);
	}
	return output;
}

void cmd::link(const std::wstring& toPath)
{
	auto to_node = get_node_from_path(toPath);
	links::create(m_db, m_cur.id, to_node.id);
}

void cmd::unlink(const std::wstring& toPath)
{
	auto to_node = get_node_from_path(toPath);
	links::remove(m_db, m_cur.id, to_node.id);
}

std::vector<std::wstring> cmd::parse_cmds(const std::wstring& cmd)
{
	std::vector<std::wstring> output;
	std::wstring collector;

	bool in_quote = false;

	for (size_t s = 0; s < cmd.length(); ++s)
	{
		wchar_t c = cmd[s];
		if (c == '\"')
		{
			if (!in_quote)
				collector = trim(collector);

			if (in_quote || !collector.empty())
			{
				output.emplace_back(collector);
				collector.clear();
			}

			in_quote = !in_quote;
			continue;
		}

		if (!in_quote && c == ' ')
		{
			collector = trim(collector);
			if (!collector.empty())
			{
				output.emplace_back(collector);
				collector.clear();
			}
			continue;
		}

		collector += c;
	}

	collector = trim(collector);
	if (!collector.empty())
	{
		output.emplace_back(collector);
		collector.clear();
	}

	return output;
}

std::wstring get_cur_path();
node get_node_from_path(const std::wstring& path);

void cd(const std::wstring& newPath);
std::vector<std::wstring> dir();

void mknode(const std::wstring& newNodeName);
void copy(const std::wstring& newParent);
void move(const std::wstring& newParent);
void remove();
void rename(const std::wstring& newName);

void set_prop(const std::vector<std::wstring>& cmds);
void set_payload(const std::wstring& payload);

std::wstring tell();

std::wstring search(const std::vector<std::wstring>& cmd);

void link(const std::wstring& toPath);
void unlink(const std::wstring& toPath);

static std::vector<std::wstring> parse_cmds(const std::wstring& cmd);
*/
