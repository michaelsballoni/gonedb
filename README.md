# gonedb
gonedb is a file-based graph DB library that gives you the power of a graph DB in your Go program
- No server required, just storage for a single database file
- More background on graph databases can be found on [Wikipedia](https://en.wikipedia.org/wiki/Graph_database)

As far as gonedb goes, being a graph database is nothing fancier than nodes with optional payloads, links between nodes with optional payloads, and name-value properties on nodes and links

The best way to dive in is to check out the pkg/cmd.go class and run the associated cmd application so you can play around

Run these commands to flex things a bit:

    > make root
    > cd root
    root> seed ..\..
    root> scramblelinks
    root> bloomcloud

seed walks the file system at the path adding nodes to the current node, recursively.
Use whatever directory path relative to the cwd, or absolute directory path, that will give you 1000's of file system entries, without blowing up your system
I use ..\.. as it blows up nicely on my system

scramblelinks walks the count of all nodes adding links between a random from node and a random to node.

bloomcloud builds a cloud using the current node as seed, then extends out three generations, outputting the new links at each go.  This can really blow up your console!

## So What Is GoneDB for?
You add your nodes and links with payloads and properties to the database, then you can search for nodes and links\
You can create "clouds" of links growing out from a seed node and enumerate the generations of nodes and links that bloomed

## Nodes
    type Node struct {
        Id           int64
        ParentId     int64
        NameStringId int64
        TypeStringId int64
    }

You create your nodes with Nodes.Create(db, parentNodeId, nameStringId, typeStringId)\
You pass around db *sql.DB as the global data processing engine, the SQLite database connection\
You get the string IDs from Strings.GetId(db *sql.DB, str string)\
Nodes.Create returns a Node, a struct that is passed around

Once you have a created a node and have its ID, you get get a fresh struct from the global pool with Nodes.Get(db, nodeId)\
The Nodes class keeps the pool up-to-date among its functions; anything you do out from under the pool requires calls into the NodeCache class,\
which has Get(), Put(), Flush(), Invalidate1(nodeId), and InvalidateN(nodeIds)

Nodes implement the graph as a tree; the Link / Links classes implement freeform connections; we'll go over them next 

You don't have to use gonedb as a tree or a graph if it doesn't fit you...you can use both!  

Same with payloads and properties, you only pay for what you use

### Tree stuff
- Copy(db, nodeId, newParentNodeId) (newNodeId, error)
    * Properties are copied but links are not (desired?)
- Move(db, nodeId, newParentNodeId) error
- Remove(db, nodeId) error
    * Properties removed, but (bug) links are not
- Rename(db, nodeId, newNameStringId) error

### Payloads associate nodes with strings, perhaps JSON or HTML or XML?
- GetPayload(db, nodeId) (string, error)
- SetPayload(db, nodeId, payload) error

## Links
type Link struct {
	Id           int64
	FromNodeId   int64
	ToNodeId     int64
	TypeStringId int64
}

Links.Create(db, fromNodeId, toNodeId, typeStringId) (Link, error)
 - you create links from one node to another (or the same) node

Links.RemoveFromTo(db, fromNodeId, toNodeId, typeStringId) error
- you remote links from one node to another (or the same) node

Links.Get(db, linkId) (Link, error)
- Get a link by ID

### Payloads associate links with strings, perhaps JSON or HTML or XML?
- GetPayload(db, linkId) (string, error)
- SetPayload(db, linkId, payload) error

### Get all links out from a node or in to a node
- GetOutLinks(db, nodeId) ([]Link, error)
- GetToLinks(db, nodeId) ([]Link, error)

## Props
You can add name/value pairs onto nodes and links:

Props.Set(db, itemTypeId, itemId, nameStringId, valueStringId) error
- the itemTypeId is either of the special values NodeItemTypeId or LinkItemTypeId
- this item business is so that the node and link props can be handled equivalently by search
- it's a bit of future-proofing, too

Props.Get(db, itemTypeId, itemId, nameStringId) (int64, error)
- returns the matching value string ID for a given item and name ID, or error

Props.GetAll(db *sql.DB, itemTypeId int64, itemId int64) (map[int64]int64, error)
- remove a mapping of string IDs, name-to-value, for a given item

## Search
You can search on
- Payload contents
- Node names
- Property names and values
- Node parent
- Node descendant

Searching is performed using:

- Search.FindNodes(db, searchQuery *SearchQuery) ([]Node, error)
or
- Search.FindLinks(db, searchQuery *SearchQuery) ([]Link, error)

type SearchQuery struct {
	Criteria       []SearchCriteria
	OrderBy        string
	OrderAscensing bool
	Limit          int64
}

type SearchCriteria struct {
	NameStringId int64
	ValueString  string
	UseLike      bool
}

And that's all of search: you give it criteria, it ANDs it all together and gives you results

## Cloud
You start with a node and follow links to expand a group of nodes, links really, a cloud.  Cloud computing, no?

Clouds.GetCloud(cloudName string, seedNodeId int64) (Cloud, error)
- creates the internal worktable name and captures the seed node ID into the struct returned
- the worktables is a lot like the lists table

Cloud.Drop(db) error
- drops the cloud's worktable

Cloud.Init(db) error
- create's the cloud's worktable

Cloud.Seed(db) (int64, error)
- initialize the worktable with the seed node's links; returns the number of links added

Cloud.GetLinks(db, minGeneration, maxGeneration) ([]Link, error)
- get the links in the cloud in a generation range

Cloud.GaxMaxGeneration(db) (int64, error)
- how far out have we gone
- NOTE: If a cloud can no longer be expanded, all links have been enumerated, the return value of this function will stop changing

Cloud.Expand(db *sql.DB) (int64, error)
- expand the cloud out one generation; returns the number of links added
