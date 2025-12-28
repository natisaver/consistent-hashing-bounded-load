package models

import (
	"fmt"
	"math"
	"sort"

	"github.com/cespare/xxhash"
)

// main function => redistributePartitions is called whenever servers are added/removed
// It redistributes partitions among servers while respecting the bounded load constraint
func (r *Ring) redistributePartitions() int {
	prevPartitionState := r.copyPreviousPartitions()

	newServerLoad := make(map[string]float64)        // Tracks how many partitions each server currently has
	newPartitionState := make(map[int]*Server)       // New assignment of partitions to servers
	partitionsMoved := 0                             // Counter for partitions that changed servers
	maxAllowedLoad := r.maxAllowedLoadPerServer()    // Maximum partitions per server allowed

	for partitionID := 0; partitionID < r.config.PartitionCount; partitionID++ {
		// for each partition, find the next closest server that can take it without exceeding max load
		server := r.findNextClosestServerForPartition(partitionID, newServerLoad, maxAllowedLoad)
		// if all servers full, cannot assign partition
		if server == nil {
			fmt.Printf("Failed to assign partition %d, all servers are full\n", partitionID)
			fmt.Printf("Consider increasing maxAllowedLoad or adding more servers\n")
			break
		}
		
		// if server found for partition
		// compare with previous assigned server, to see if partition moved
		if prevServer := prevPartitionState[partitionID]; prevServer != nil && prevServer.Name != server.Name {
			partitionsMoved++
		}

		// assign partition to the server and update load of server
		newPartitionState[partitionID] = server
		newServerLoad[server.Name]++
	}

	// update ring states
	r.partitions = newPartitionState
	r.serverLoads = newServerLoad

	return partitionsMoved
}

// creates a new map, to store current partition id -> server assignments
func (r *Ring) copyPreviousPartitions() map[int]*Server {
	prev := make(map[int]*Server)
	for k, v := range r.partitions {
		prev[k] = v
	}
	return prev
}

//applies consistent hashing + bounded load
func (r *Ring) findNextClosestServerForPartition(
	partitionID int,
	newServerLoad map[string]float64,
	maxAllowedLoad float64,
) *Server {
	// ring is circular, so if the hash is bigger than the largest virtual server hash
	// should be assigned to first virtual server hash in ring
	hash := utils.HashPartition(partitionID)
	serverIdxToAssign := sort.Search(len(r.sortedRing), func(i int) bool {
		return r.sortedRing[i] >= hash
	})
	if serverIdxToAssign == len(r.sortedRing) {
		serverIdxToAssign = 0
	}

	// check in clockwise order if server can take partition without exceeding max load
	for i := 0; i < len(r.sortedRing); i++ {
		virtualNode := r.sortedRing[serverIdxToAssign]
		server := r.virtualNodeMap[virtualNode]
		// test if +1 to add this partition, is it within bounds
		if newServerLoad[server.Name]+1 <= maxAllowedLoad {
			return server
		}
		// check next server in ring if cur server full, % as ring is circular
		serverIdxToAssign = (serverIdxToAssign + 1) % len(r.sortedRing)
	}

	// all servers full, cannot assign partition
	return nil
}