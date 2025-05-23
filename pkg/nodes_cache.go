package gonedb

import "sync"

type cache struct{}

var NodeCache cache

var g_cacheLock sync.RWMutex
var g_cache = make(map[int64]Node)

func (n *cache) Get(id int64) (Node, bool) {
	g_cacheLock.RLock()
	node, found := g_cache[id]
	g_cacheLock.RUnlock()
	return node, found
}

func (n *cache) Put(node Node) {
	g_cacheLock.Lock()
	g_cache[node.Id] = node
	g_cacheLock.Unlock()
}

func (n *cache) Flush() {
	g_cacheLock.Lock()
	clear(g_cache)
	g_cacheLock.Unlock()
}

func (n *cache) Invalidate1(nodeId int64) {
	g_cacheLock.Lock()
	delete(g_cache, nodeId)
	g_cacheLock.Unlock()
}

func (n *cache) InvalidateN(nodeIds []int64) {
	g_cacheLock.Lock()
	for _, nodeId := range nodeIds {
		delete(g_cache, nodeId)
	}
	g_cacheLock.Unlock()
}
