package models

import (
	"fmt"
	"math"
	"sort"
)

// given the current servers and total partitions, what is the maximum load per server, adjusted by LoadFactor?
func (r *Ring) maxAllowedLoadPerServer() float64 {
	if len(r.serverList) == 0 {
		return 0 // if no servers, no load
	}

	// Divide partitions by servers → gives ideal load per server
	numPartitionsPerServer := float64(r.config.PartitionCount) / float64(len(r.serverList))
	// LoadFactor = extra “headroom” for each server, round up to whole number
	// 100 partitions / 4 servers = 25
	// LoadFactor = 1.25 → avgLoad = 31.25
	// Means each server can safely hold up to ~32 partitions before we start moving to next server
	maxLoadPerServer := numPartitionsPerServer * r.config.LoadFactor
	return math.Ceil(maxLoadPerServer)
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
