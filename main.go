package main

import (
	"fmt"
	"natisaver/consistenthashbound/models"
)

func main() {
	// main.go -> server.go -> ring.go -> distribute.go

	initservers := []*models.Server{}
	// start with 8 servers
	for i := 0; i < 8; i++ {
		initservers = append(initservers, &models.Server{Name: fmt.Sprintf("node%d", i)})
	}
	
	cfg := models.Config{
		// partitions => like buckets that store the data keys
		// partitionID = hash(key) % prime number, so that keys are spread out without repeated patterns
		PartitionCount:   271, 
		// virtual nodes => representations of the server on the ring
		// With 10 servers × 50 virtual nodes = 500 vnodes,
		// each vnode owns 271/500 ≈ 0.5 partitions (either 0 or 1), in this way no server controls meaningful load
		// VirtualNodeCount ≥ 2 × PartitionCount / ServerCount
		VirtualNodeCount: 50, 
		// load factor => max allowable % load of server, 1.2 allows 20% more load
		// it should never be < 1, as that means a partition will never be assigned, skipping forever in a loop
		// cant be 1 as that expects each server to hold exactly the average no. of paritions, but hashing is random
		LoadFactor:       1.2,
	}
	// init ring with 8 servers
	r := models.NewRing(initservers, cfg)

	// add 2 more servers, leading to partition movement
	addservers := []*models.Server{}
	for i := 8; i < 10; i++ {
		addservers = append(addservers, &models.Server{Name: fmt.Sprintf("node%d", i)})
	}
	r.AddServers(addservers)

}
