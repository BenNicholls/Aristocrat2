package main

import (
	"fmt"
	"math/bits"
)

func bitboardToString(b uint64) string {
	return fmt.Sprintf("%064b", b)
}

func outputBitboard(b uint64) {
	s := bitboardToString(b)
	for i := 0; i < 8; i++ {
		fmt.Println(s[i*8 : (i+1)*8])
	}
}

func countBits(b uint64) int {
	return bits.OnesCount64(b)
}

func checkBit(b uint64, i int) bool {
	return ((1 << (63 - i)) & b) != 0
}

func setBit(b uint64, i int) uint64 {
	return (1 << (63 - i)) | b
}

func clearBit(b uint64, i int) uint64 {
	return b &^ (1 << (63 - i))
}

func forEachBit(b uint64, f func(square int)) {
	for i := 0; b != 0; {
		i = bits.LeadingZeros64(b)
		f(i)
		b = clearBit(b, i)
	}
}

//reports position of leftmost bit
func leftBit(b uint64) int {
	return bits.LeadingZeros64(b)
}

//reports position of rightmost bit
func rightBit(b uint64) int {
	return 63 - bits.TrailingZeros64(b)
}
