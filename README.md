## Adding Servers
- dynamically added 2 more servers from 8 -> 10
- only about 16% of partitions had to be redistributed
- load is balanced between all servers

Config
```go
	cfg := models.Config{
		PartitionCount:   271,
		VirtualNodeCount: 40,
		LoadFactor:       1.2,
	}
```

Output
```bash
go run "d:\coding projects\consistenthashbound\main.go"
====RESULTS=====
Ideal Average Load: 40.00
Actual Average Load: 33.88
Min Load: 27.00, Max Load: 39.00
--------------
Total Number of Servers: 8
Partitions Redistributed: 0/271 (0.00%)
--------------
Server: node2, Load (No. Partitions): 36.00/40.00
Server: node4, Load (No. Partitions): 27.00/40.00
Server: node3, Load (No. Partitions): 30.00/40.00
Server: node0, Load (No. Partitions): 34.00/40.00
Server: node7, Load (No. Partitions): 37.00/40.00
Server: node5, Load (No. Partitions): 39.00/40.00
Server: node1, Load (No. Partitions): 37.00/40.00
Server: node6, Load (No. Partitions): 31.00/40.00

====RESULTS=====
Ideal Average Load: 33.00
Actual Average Load: 27.10
Min Load: 22.00, Max Load: 32.00
--------------
Total Number of Servers: 10
Partitions Redistributed: 45/271 (16.61%)
--------------
Server: node2, Load (No. Partitions): 31.00/33.00
Server: node4, Load (No. Partitions): 24.00/33.00
Server: node3, Load (No. Partitions): 22.00/33.00
Server: node9, Load (No. Partitions): 23.00/33.00
Server: node0, Load (No. Partitions): 27.00/33.00
Server: node8, Load (No. Partitions): 22.00/33.00
Server: node7, Load (No. Partitions): 30.00/33.00
Server: node5, Load (No. Partitions): 32.00/33.00
Server: node1, Load (No. Partitions): 32.00/33.00
Server: node6, Load (No. Partitions): 28.00/33.00

```
- ![image](https://github.com/user-attachments/assets/d0e1ccfa-60a7-43b8-9da0-5f1c02f813fa)



For future improvements,use self-balancing tree structure, such as an AVL tree or a Red-Black tree. Self-balancing trees provide logarithmic time complexity for insertion, deletion, and lookup operations, which would significantly enhance the efficiency of the ring management. This change would allow for more scalable and performant handling of virtual nodes, especially in dynamic environments where virtual nodes are frequently added or removed.
