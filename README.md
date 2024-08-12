## Adding Servers
- dynamically added 2 more servers from 8 -> 10
- only about 16% of partitions had to be redistributed
- load is balanced between all servers
```go
	cfg := models.Config{
		PartitionCount:   271,
		VirtualNodeCount: 40,
		LoadFactor:       1.2,
	}
```
- ![image](https://github.com/user-attachments/assets/d0e1ccfa-60a7-43b8-9da0-5f1c02f813fa)



For future improvements,use self-balancing tree structure, such as an AVL tree or a Red-Black tree. Self-balancing trees provide logarithmic time complexity for insertion, deletion, and lookup operations, which would significantly enhance the efficiency of the ring management. This change would allow for more scalable and performant handling of virtual nodes, especially in dynamic environments where virtual nodes are frequently added or removed.
