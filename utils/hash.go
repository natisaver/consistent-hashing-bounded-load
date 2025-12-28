package utils

import (
	"encoding/binary"

	"github.com/cespare/xxhash"
)

// xxhash for speed, deterministic output, uniform distribution so replicas are spread out evenly

// hashVirtualNode hashes "serverName + virtual replica index"
func HashVirtualNode(serverName string, replica int) uint64 {
	key := serverName + strconv.Itoa(replica)
	return xxhash.Sum64String(key)
}

// hashPartition hashes partition ID to uint64
// uint64, 64bits => 8 bytes
// make a fixed 8-byte representation for each partition (identified by id) to hash
func HashPartition(partitionId int) uint64 {
	var buf [8]byte // hash function cant work on int direct, only on array of bytes, 8 bytes buffer to represent uint64
	binary.LittleEndian.PutUint64(buf[:], uint64(partitionId)) // convert partitionId to uint64 and store in buffer, MSB, LSB doesnt matter
	return xxhash.Sum64(buf[:]) // hash the byte array
}