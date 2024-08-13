package models

import (
	"fmt"
	"sort"

	"github.com/cespare/xxhash"
)

// xxhash for speed, deterministic output, uniform distribution so replicas are spread out evenly

const (
	DefaultPartitionCount   int = 271
	DefaultVirtualNodeCount int = 20
	// If there are 100 partitions and 10 servers, ideally, each server would manage 10 partitions.
	// Load Factor = 1.0, each server is expected to handle the ideal no. of 10 partitions.
	// Load Factor = 1.25, each server is allowed to handle up to 25% more partitions than the ideal number.
	DefaultLoadFactor float64 = 1.25 // margin of imbalance relative to ideal load

)

type Server struct {
	Name string
}

type Ring struct {
	config Config

	// virtual nodes
	sortedRing     []uint64           // can improve using self-balancing trees, O(lg(N)) vs O(N) for add/delete
	virtualNodeMap map[uint64]*Server // maps virtual node hash to server

	// servers
	serverList  map[string]*Server // maps server name to server
	serverLoads map[string]float64 // maps server name to no. of partitions on it

	// partitions
	partitions map[int]*Server // maps partition id to server

}

type Config struct {
	// Keys are distributed among partitions. Prime numbers are good to
	// distribute keys uniformly. Select a big PartitionCount if you have
	// too many keys.
	PartitionCount int
	// each server is represented multiple times on the ring to distribute load
	VirtualNodeCount int
	// Load is used to calculate average load. See the code, the paper and Google's blog post to learn about it.
	LoadFactor float64
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
	}

	r.AddServers(servers)
	return r
}

// adds a new server to ring
func (r *Ring) Add(server *Server) {
	// add server to server list
	r.serverList[server.Name] = server
	// create virtual nodes on ring for server
	for i := range r.config.VirtualNodeCount {
		virutalNodeKey := fmt.Sprintf("%s%d", server.Name, i)
		hashedKey := xxhash.Sum64String(virutalNodeKey)
		// map virtual nodes to the server
		r.virtualNodeMap[hashedKey] = server
		// add virtual node to ring
		r.sortedRing = append(r.sortedRing, hashedKey)
	}
	// sort the ring
	// hashkeys in ascending order
	sort.Slice(r.sortedRing, func(i, j int) bool {
		return r.sortedRing[i] < r.sortedRing[j]
	})

	// since server added
	// redistribute partitions
	movedPartitions := r.distributePartitionsAndLoad()

	r.printMetrics(movedPartitions)
}

// adds a list of servers to ring
func (r *Ring) AddServers(servers []*Server) {
	for _, s := range servers {
		// add server to server list
		r.serverList[s.Name] = s
		// create virtual nodes on ring for server
		for i := range r.config.VirtualNodeCount {
			virutalNodeKey := fmt.Sprintf("%s%d", s.Name, i)
			hashedKey := xxhash.Sum64String(virutalNodeKey)
			// map virtual nodes to the server
			r.virtualNodeMap[hashedKey] = s
			// add virtual node to ring
			r.sortedRing = append(r.sortedRing, hashedKey)
		}
	}

	// sort the ring
	// hashkeys in ascending order
	sort.Slice(r.sortedRing, func(i, j int) bool {
		return r.sortedRing[i] < r.sortedRing[j]
	})

	// since servers added
	// redistribute partitions
	movedPartitions := r.distributePartitionsAndLoad()

	r.printMetrics(movedPartitions)
}


// removes a new server to ring
func (r *Ring) Remove(serverName string) {
	if val, ok := r.serverList[server]; !ok {
		// server does not exist
		return
	}
	// remove the virtual nodes in the ring
	for i := range r.config.VirtualNodeCount {
		virutalNodeKey := fmt.Sprintf("%s%d", serverName, i)
		hashedKey := xxhash.Sum64String(virutalNodeKey)
		// delete virtual node from ring
		r.deleteVirtualNode(hashedKey)
		// delete from virtual node map
		delete(r.virtualNodeMap, hashedKey)
	}
	// delete server from the serverlist
	delete(r.serverList, name)
	
	// since server removed
	// redistribute partitions
	movedPartitions := r.distributePartitionsAndLoad()

	r.printMetrics(movedPartitions)
}

func (r *Ring) deleteVirtualNode(val uint64) {
	for i := 0; i < len(r.sortedRing); i++ {
		if r.sortedRing[i] == val {
			r.sortedRing = append(r.sortedRing[:i], r.sortedRing[i+1:]...)
			break
		}
	}
}
