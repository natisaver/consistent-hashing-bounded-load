package models

import (
	"fmt"
	"sort"

	"natisaver/consistenthashbound/utils"
	"github.com/cespare/xxhash"
)

type Ring struct {
	config Config

	// ring
	sortedRing     []uint64           	// sorted ring contains all the server virtual node hashes, can improve using self-balancing trees, O(lg(N)) vs O(N) for add/delete

	// virtual nodes
	virtualNodeMap map[uint64]*Server 	// Maps virtual node hashed key -> server instance

	// servers
	serverList  map[string]*Server  	// Maps server name -> server instance
	serverLoads map[string]int      	// Maps server name -> number of partitions its holding

	// partitions
	partitions map[int]*Server 			// Maps partition ID -> server instance its in
}

// initialises a new ring with a starting set of servers
func NewRing(servers []*Server, config Config) *Ring {
	// check if config is nil, use default params
	if config.PartitionCount == 0 {
		config.PartitionCount = DefaultPartitionCount
	}
	if config.VirtualNodeCount == 0 {
		config.VirtualNodeCount = DefaultVirtualNodeCount
	}
	if config.LoadFactor == 0 {
		config.LoadFactor = DefaultLoadFactor
	}

	r := &Ring{
		config:         config,
		serverList:     make(map[string]*Server),
		virtualNodeMap: make(map[uint64]*Server),
		sortedRing:     []uint64{},
		serverLoads:    make(map[string]int),
		partitions:     make(map[int]*Server),
	}

	r.AddServers(servers)
	return r
}

// adds a list of servers to ring
func (r *Ring) AddServers(servers []*Server) {
	for _, s := range servers {
		// add server to server list
		r.serverList[s.Name] = s
		// create virtual nodes on ring for server
		r.addServerVirtualNodes(s)
	}

	// since servers added
	// redistribute partitions
	movedPartitions := r.redistributePartitions()

	r.printMetrics(movedPartitions)
}

// adds a single server to the ring
func (r *Ring) AddServer(s *Server) {
	r.AddServers([]*Server{s})
}

// removes a server by server name in ring
func (r *Ring) RemoveServer(serverName string) {
	if _, ok := r.serverList[serverName]; !ok {
		// server does not exist in keys
		return
	}

	// remove all virtual nodes
	r.deleteServerVirtualNodes(serverName)

	// delete server from the serverlist
	delete(r.serverList, serverName)

	// since server removed
	// redistribute partitions
	movedPartitions := r.redistributePartitions()

	r.printMetrics(movedPartitions)
}

func (r *Ring) addServerVirtualNodes(server *Server) {
	for i := 0; i < r.config.VirtualNodeCount; i++ {
		hashedKey := hashVirtualNode(server.Name, i)
		r.insertSorted(hashedKey)
		r.virtualNodeMap[hashedKey] = server
		
	}
}

// deleteServerVirtualNodes deletes all virtual nodes of a server from the ring
func (r *Ring) deleteServerVirtualNodes(serverName string) {
	for i := 0; i < r.config.VirtualNodeCount; i++ {
		hashedKey := hashVirtualNode(serverName, i)

		// delete from sortedRing
		r.removeSorted(hashedKey)

		// delete from virtual node map
		delete(r.virtualNodeMap, hashedKey)
	}
}

// insert into sorted Ring
func (r *Ring) insertSorted(hash uint64) {
	// find idx to insert O(log(N))
	// slice = [10, 20, 30, 40], hash = 25, returns idx=2
	idx := sort.Search(len(r.sortedRing), func(i int) bool {
		return r.sortedRing[i] >= hash
	})
	// insert at idx O(N)
	r.sortedRing = append(r.sortedRing, 0)          // add dummy to extend slice
	copy(r.sortedRing[idx+1:], r.sortedRing[idx:])  // shift elements frm idx onward by 1
	r.sortedRing[idx] = hash						// insert hash at idx
}

// remove from sorted ring
func (r *Ring) removeSorted(hash uint64) {
	// O(log(N))
	idx := sort.Search(len(r.sortedRing), func(i int) bool {
		return r.sortedRing[i] >= hash
	})
	if idx == len(r.sortedRing) {
		return // not found
	}
	// O(N)
	if r.sortedRing[idx] == hash {
		r.sortedRing = append(r.sortedRing[:idx], r.sortedRing[idx+1:]...)
	}
}