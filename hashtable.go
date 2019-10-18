package main

import "sync"

type hashTable struct {
	sync.Mutex
	table []hashTableEntry
	size  uint64
}

type hashTableEntry struct {
	hash     uint64
	bestMove move
	depth    int
	score    int
	result   result
}

//size of an entry in bytes
const HASHTABLEENTRYSIZE int = 8 * 5

//size is in megabytes
func NewHashTable(size int) (ht *hashTable) {
	ht = new(hashTable)
	ht.size = uint64(size * 1024 * 1024 / HASHTABLEENTRYSIZE)
	ht.table = make([]hashTableEntry, ht.size)
	usingHashtable = true

	return
}

func (ht *hashTable) Store(hash uint64, depth int, bestMove move, score int, result result) {
	if !usingHashtable {
		return
	}
	ht.Lock()
	entry := ht.table[hash%ht.size]
	ht.Unlock()
	if entry.hash != hash { //new position beting stored. overwrite old data
		ht.Lock()
		ht.table[hash%ht.size] = hashTableEntry{
			hash:     hash,
			depth:    depth,
			bestMove: bestMove,
			score:    score,
			result:   result,
		}
		ht.Unlock()
	} else { //attempt to rewrite. rewrite if depth is higher (deeper eval)
		if entry.depth < depth {
			ht.Lock()
			ht.table[hash%ht.size].depth = depth
			ht.table[hash%ht.size].bestMove = bestMove
			ht.table[hash%ht.size].score = score
			ht.table[hash%ht.size].result = result
			ht.Unlock()
		}
	}
}

func (ht *hashTable) Load(hash uint64) (entry hashTableEntry, ok bool) {
	if !usingHashtable {
		return
	}
	ht.Lock()
	entry = ht.table[hash%ht.size]
	ht.Unlock()
	if entry.hash != 0 && entry.hash == hash {
		ok = true
	}
	return
}
