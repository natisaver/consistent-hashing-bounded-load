package models

import (
	"fmt"
	"math"
	"sort"

	"github.com/cespare/xxhash"
)

func (r *Ring) averageLoad() float64 {
	if len(r.serverList) == 0 {
		return 0
	}

	numPartitionsPerServer := float64(r.config.PartitionCount / (len(r.serverList)))
	avgLoad := numPartitionsPerServer * r.config.LoadFactor
	return math.Ceil(avgLoad)
}

// Main consistent hashing with bounded load algorithm
func (r *Ring) distributePartitionsAndLoad() int {
	// Store the previous partition assignments
	previousPartitions := make(map[int]*Server)
	for partID, server := range r.partitions {
		previousPartitions[partID] = server
	}

	// recalculate each server's loads and the mapping of partition to server
	serverLoads := make(map[string]float64)
	partitions := make(map[int]*Server)
	movedPartitions := 0

	// Iterate over each partition ID
	// and distribute the partitions into each server
	for partID := 0; partID < r.config.PartitionCount; partID++ {
		// hash partition id
		hashedPtnKey := xxhash.Sum64([]byte{byte(partID)})

		// With the hashed parition key
		// find the closest virtual node on the ring to assign the partition to
		virtualNodeIndex := sort.Search(len(r.sortedRing), func(i int) bool {
			return r.sortedRing[i] >= hashedPtnKey
		})

		// Apply Consistent Hashing with Bounds
		// attempt to store partition into server
		// if server is full, move clockwise to next virtual node
		avgLoad := r.averageLoad()
		success := false

		for attempts := 0; attempts < len(r.sortedRing); attempts++ {
			virtualNode := r.sortedRing[virtualNodeIndex]
			server := *r.virtualNodeMap[virtualNode]

			// bound is not full
			// server has capacity for this partition
			if serverLoads[server.Name]+1 <= avgLoad {
				// Check if partition is moved
				if prevServer := previousPartitions[partID]; prevServer != nil && prevServer.Name != server.Name {
					movedPartitions++
				}
				// Assign partition to the server and update the load
				serverLoads[server.Name]++
				partitions[partID] = &server
				// fmt.Printf("Success, assigned partition %d to server %s\n", partID, server.Name)
				success = true
				break
			}

			// bound is full, check next clockwise virtual node
			virtualNodeIndex = (virtualNodeIndex + 1)
		}

		// If after all attempts no server could take the partition, print the failure message
		if !success {
			fmt.Printf("Failed to redistribute partition %d, all servers are full\n", partID)
			fmt.Printf("Consider increasing averageLoad or increasing number of servers\n")
			break
		}

	}

	// Update the partition and load maps in the object
	r.partitions = partitions
	r.serverLoads = serverLoads

	return movedPartitions
}

func (r *Ring) printMetrics(movedPartitions int) {
	avgLoad := r.averageLoad()
	totalServers := len(r.serverLoads)
	min, max, avg := r.getMinMaxAvgLoadMetrics(totalServers)

	fmt.Println("====RESULTS=====")
	fmt.Printf("Ideal Average Load: %.2f\n", avgLoad)
	fmt.Printf("Actual Average Load: %.2f\n", avg)
	fmt.Printf("Min Load: %.2f, Max Load: %.2f\n", min, max)
	fmt.Println("--------------")
	fmt.Printf("Total Number of Servers: %d\n", totalServers)

	fmt.Printf("Partitions Redistributed: %d/%d (%.2f%%)\n", movedPartitions, r.config.PartitionCount, (float64(movedPartitions)/float64(r.config.PartitionCount))*100)
	fmt.Println("--------------")
	for serverName, load := range r.serverLoads {
		fmt.Printf("Server: %s, Load (No. Partitions): %.2f/%.2f\n",
			serverName, load, avgLoad)
	}
	fmt.Println("")
}

func (r *Ring) getMinMaxAvgLoadMetrics(totalServers int) (float64, float64, float64) {
	// Calculate the total, min, and max load
	var totalLoad, minLoad, maxLoad float64
	if totalServers > 0 {
		minLoad = r.serverLoads[fmt.Sprintf("node%d", 0)]
		maxLoad = r.serverLoads[fmt.Sprintf("node%d", 0)]
	}

	for _, load := range r.serverLoads {
		totalLoad += load
		if load < minLoad {
			minLoad = load
		}
		if load > maxLoad {
			maxLoad = load
		}
	}

	actualAvgLoad := totalLoad / float64(totalServers)
	return minLoad, maxLoad, actualAvgLoad
}
