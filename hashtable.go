package main

import "sync"

type nodeType int

const (
	EXACT nodeType = iota
	LOWER
	UPPER
)

type hashTable struct {
	sync.Mutex
	table []hashTableEntry
	size  uint64
}

//TODO: squash this down to use less space. nodeType, result, score, depth, none of these need to be 8 bytes wide
type hashTableEntry struct {
	hash     uint64
	bestMove move
	depth    int
	score    int
	result   result
	node     nodeType
}

//size of an entry in bytes
const HASHTABLEENTRYSIZE int = 8 * 6

//size is in megabytes
func newHashTable(size int) (ht *hashTable) {
	ht = new(hashTable)
	ht.size = uint64(size * 1024 * 1024 / HASHTABLEENTRYSIZE)
	ht.table = make([]hashTableEntry, ht.size)
	return
}

func initHashTable() {
	if hashSize == 0 {
		table = &hashTable{}
		usingHashtable = false
	} else {
		table = newHashTable(hashSize)
		usingHashtable = true
	}
}

func (ht *hashTable) Store(hash uint64, depth int, bestMove move, score int, result result, node nodeType) {
	if !usingHashtable {
		return
	}
	ht.Lock()
	entry := ht.table[hash%ht.size]
	if entry.hash != hash { //new position being stored. overwrite old data
		ht.table[hash%ht.size] = hashTableEntry{
			hash:     hash,
			depth:    depth,
			bestMove: bestMove,
			score:    score,
			result:   result,
			node:     node,
		}
	} else { //attempt to rewrite. rewrite if depth is higher (deeper eval)
		if entry.depth < depth {
			ht.table[hash%ht.size].depth = depth
			ht.table[hash%ht.size].bestMove = bestMove
			ht.table[hash%ht.size].score = score
			ht.table[hash%ht.size].result = result
			ht.table[hash%ht.size].node = node
		}
	}
	ht.Unlock()
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
