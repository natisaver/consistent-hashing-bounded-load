package models

import (
	"fmt"
	"sort"

)

// default configuration
const (
	DefaultPartitionCount   int = 271
	DefaultVirtualNodeCount int = 20
	// If there are 100 partitions and 10 servers:
	// Load Factor = 1.0, each server is then expected to handle the ideal no. of 10 partitions.
	// Load Factor = 1.25, each server is allowed to handle up to 25% more partitions than the ideal number.
	DefaultLoadFactor float64 = 1.25 // margin of imbalance relative to ideal load

)

// simulate server with a servername
type Server struct {
	Name string
}

// configurations for each server
type Config struct {
	// Keys are distributed among partitions. Prime numbers are good to
	// distribute keys uniformly. Select a big PartitionCount if you have
	// too many keys.
	PartitionCount int
	// each server is represented multiple times on the ring to distribute load
	VirtualNodeCount int
	// Bounded load for each server in ring, can refer to Google's blog post to learn about it.
	LoadFactor float64
}