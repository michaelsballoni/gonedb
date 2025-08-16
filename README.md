# gonedb

## So What Is gonedb for?
gonedb gives you the power of a graph DB in your Go program

As far as gonedb goes, being a graph database is nothing fancier than...
+ nodes with optional payloads
+ links between nodes with optional payloads
+ and name-value properties on nodes and links

Then you can...
- Add graph nodes with parent relationships, names, custom payloads, and properties
- Add graph links with a from and to node relationships, custom payloads, and properties
- Traverse graphs via parent or link relationships
- Search nodes and links by name, payload, or properties
- Generate "clouds" of connected links growing out from a seed node

## Background & Installation
- No server required, just storage for a single database file (SQLite)
- More background on graph databases can be found on [Wikipedia](https://en.wikipedia.org/wiki/Graph_database)
- Install Go [here](https://go.dev/doc/install)
- Install SQLite on Go [here](https://medium.com/@yaravind/go-sqlite-on-windows-f91ef2dacfe)

## Getting Started - POC cmd program
The best way to dive into gonedb in is to check out the cmd application and its associated pkg/cmd.go class below

NOTE: VERY IMPORTANT: 
The cmd application and its cmd.go class are NOT how you use gonedb.  
You would use classes like pkg/nodes.go and others described after this POC sample to add gonedb to your own application.\
So enjoy this cmd.go stuff, just know that it's just a POC, not the gonedb for you to use.

    C:\Users\...\source\gonedb\cmd> go run cmd.go cmd-test.db
    db_file: cmd-test.db file_exists: false
    Opening database...
    SQLite version: 3.46.1
    Setting up gonedb schema...done!
    > rem rem is just for rem-arkihng on the other commands
    > rem That db_file, Opening database, that's a header as part of running the cmd POC
    >
    > rem To create a node you run the *make* command, specifying just the name of the node
    > make root
    > rem To list the name paths of the direct descendants of a node use the *dir* command
    > dir
    root

    > rem Let's make a new node inside root, and try the commands again
    > rem First we change the current node to root
    > cd /root
    root> rem Notice that the command prompt includes root now
    root> 
    root> rem Now we make a child node
    root> make child1
    root> rem And let's see the full paths of nodes directly inside root with the *dir* command
    root> dir
    root/child1

    root> rem Let's create a couple grandchildren; full paths are needed in most situations with most commands
    root> cd /root/child1
    root/child1> make grandchild1
    root/child1> make grandchild2
    root/child1> dir
    root/child1/grandchild1
    root/child1/grandchild2

    root/child1> rem Now let's create new a child of root, and make a *copy* of child1 in it
    root/child1> cd /root
    root> make child2
    root> dir
    root/child1
    root/child2
    root> cd /root/child1
    root/child1> tell
        ID: 2
        Name: child1
        Parent: root
        Payload:
        Properties: (none)
        Out Links: (none)
        In Links: (none)

    root/child1>
    root/child1> copy /root/child2
    root/child1> cd /root/child2
    root/child2> dir
    root/child2/child1

    root/child2> cd /root/child2/child1
    root/child2/child1> dir
    root/child2/child1/grandchild1
    root/child2/child1/grandchild2

    root/child2/child1> cd /root
    root> rem You see that child1 was copied under child2,
    root> rem and child1's children were copied over as well

    root> rem Let's make a new child3 under root and *move* child1's second grandchild under it
    root> make child3
    root> dir
    root/child1
    root/child2
    root/child3

    root> cd /root/child2/child1/grandchild2
    root/child2/child1/grandchild2> move /root/child3
    root/child3/grandchild2> cd /root/child3
    root/child3> dir
    root/child3/grandchild2

    root/child3> cd /root
    root> rem You see that grandchild2 was copied under child3

    root> rem Let's *remove* the child1 copy from child2
    root> cd root/child2/child1
    root/child2/child1> remove
    root/child2> dir
    root/child2> cd /root
    root> rem You see that nothing was output by child2's dir; child1's copy is no longer in child2

    root> rem Let's *rename* child3 to child0
    root> cd /root/child3
    root/child3> rename child0
    root/child0> tell
        ID: 9
        Name: child0
        Parent: root
        ...
    root/child0> cd /root

    root> rem Let's put a property on child0 and use *search* to find it by name and by property
    root> cd child0
    root/child0> setprop property-name property-value
    root/child0> search property-name not-the-right-value
    root/child0> search property-name property-value
    root/child0

    root/child0> search name child-wrong-name
    root/child0> search name child0
    root/child0

    root/child0> cd /root

    root> rem Let's set the payload on child0 and view it using tell
    root> cd child0
    root/child0> setpayload "this is the payload"
    root/child0> tell
        ID: 9
        Name: child0
        Parent: root
        Payload: this is the payload
        Properties:
        property-name property-value
        ...
    root/child0> cd /root

    root> rem Let's *link* child0 to child1, tell child0, then *unlink* the two and tell again
    root> cd child0
    root/child0> link root/child1
    root/child0> tell
        ID: 9
        Name: child0
        ...
        Out Links: (1)
        root/child1
        In Links: (none)
    root/child0> unlink root/child1
        root/child0> tell
        ID: 9
        Name: child0
        ...
        Out Links: (none)
        In Links: (none)
    root/child0> cd /root

    root> rem You can use the gonedb you've glimpsed above, it's wonderful.  But maybe you need a little magic.

    root> rem First we issue the scramblelinks command which creates random links between all the nodes
    root> rem Again, this is totally POC, you would probably never want to do this in real code.  
    root> rem And it can't be undone.  
    root> rem But if you need a lot of links, this does the job.
    root> scramblelinks
    links created: 6

    root> rem Then we cast our magic wand, bloomcloud, which takes the current node as the seed of blowing up generations of links
    root> rem It manages a database table for each cloud, and each cloud table is a list of links
    root> rem Each expansion grows the cloud to include links that either...
    root> rem A) link from outside the cloud into a node in the cloud
    root> rem B) link from inside the cloud out to a node outside the cloud
    root> rem bloomcloud does out to three expansions

    root> bloomcloud
    Gen: 1
    Gen: 1 - Added: 3
    3 4
    3 9
    8 3

    Gen: 2
    Gen: 2 - Added: 2
    9 4
    9 2

    Gen: 3
    Gen: 3 - Added: 0
    root> rem After the "Added: " lines, each line is a link, from node ID, space, to node ID

And that's it.  

Remember, this has been a POC for the purpose of this walkthrough.  Do not use it for any production work.  

The documentation of the gonedb API is where the road really hits the road, I hope you've stuck around for that.

# Reference

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

Here are the structs you pass into the functions:

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

And that's all of search: you give it sensible criteria, it ANDs it all together and gives you sensible results

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
- each expansion grows the cloud to include links that either:
- link from outside the cloud into a node in the cloud
- link from inside the cloud out to a node outside the cloud
