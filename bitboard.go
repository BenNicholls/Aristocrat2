package main

import (
	"fmt"
	"math/bits"
)

type bitboard uint64

func (b bitboard) string() string {
	return fmt.Sprintf("%064b", b)
}

func (b bitboard) print() {
	s := b.string()
	for i := 0; i < 8; i++ {
		fmt.Println(s[i*8 : (i+1)*8])
	}
}

func (b bitboard) count() int {
	return bits.OnesCount64(uint64(b))
}

func checkBit(b bitboard, i int) bool {
	if i < 0 || i > 63 {
		return false
	}

	return ((1 << (63 - i)) & b) != 0
}

func setBit(b bitboard, i int) bitboard {
	if i < 0 || i > 63 {
		return b
	}

	return (1 << (63 - i)) | b
}

func clearBit(b bitboard, i int) bitboard {
	return (1 << (63 - i)) &^ b
}
