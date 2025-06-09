# Finish Port from C++

### Search
Add back Nodes machinery for looking up nodes by name string ID paths, deep or not
    static std::optional<std::wstring> get_path_to_parent_like(db& db, const std::wstring& path);

### Command Processor
class cmd
{
public:
    cmd(db& db);

    void mount(const std::wstring& dirPath);

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

private:
    db& m_db;
    node m_cur;
    int64_t m_nodeItemTypeId;
};

### Command Line POC utility
...

## Do the porting

# Move Forward w/Reforms and New Features
## Rethink cloud library
One table with just node IDs and generations, then branch from and to the latest generation
