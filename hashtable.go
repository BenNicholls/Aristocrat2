package main

import "sync"

type hashTable struct {
	sync.RWMutex
	table []hashTableEntry
	size  uint64
}

type hashTableEntry struct {
	hash     uint64
	depth    int
	bestMove move
	score    int
}

//size of an entry in bytes
const HASHTABLEENTRYSIZE int = 8 * 4

//size is in megabytes
func NewHashTable(size int) *hashTable {
	ht := new(hashTable)
	size = size / 2 //NOTE: this is here because for some reason i can't figure out, aristocrat
	//uses twice as much memory as requested. TODO: figure this out.
	ht.size = uint64(size * 1000 * 1000 / HASHTABLEENTRYSIZE)
	ht.table = make([]hashTableEntry, ht.size)

	return ht
}

func (ht *hashTable) Store(hash uint64, depth int, bestMove move, score int) {
	ht.Lock()
	ht.table[hash%ht.size] = hashTableEntry{
		hash:     hash,
		depth:    depth,
		bestMove: bestMove,
		score:    score,
	}
	ht.Unlock()
}

func (ht *hashTable) Load(hash uint64) (hashTableEntry, bool) {
	ht.RLock()
	entry := ht.table[hash%ht.size]
	ht.RUnlock()
	if entry.hash == 0 || entry.hash != hash {
		return entry, false
	}
	return entry, true
}
