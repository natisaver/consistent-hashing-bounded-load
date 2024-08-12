package main

import (
	"fmt"
	"natisaver/consistenthashbound/models"
)

func main() {
	servers := []*models.Server{}
	for i := 0; i < 8; i++ {
		servers = append(servers, &models.Server{Name: fmt.Sprintf("node%d", i)})
	}
	cfg := models.Config{
		PartitionCount:   271,
		VirtualNodeCount: 40,
		LoadFactor:       1.2,
	}
	r := models.NewRing(servers, cfg)

	// add 2 servers
	servers2 := []*models.Server{}
	for i := 8; i < 10; i++ {
		servers2 = append(servers2, &models.Server{Name: fmt.Sprintf("node%d", i)})
	}
	r.AddServers(servers2)

}
