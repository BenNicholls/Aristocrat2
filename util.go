package main

import "fmt"

var seed uint64 = 80

func algebraicToSquare(s string) int {
	if len(s) != 2 {
		return 0
	}

	square := 0
	switch s[0] {
	case 'a':
		square = 0
	case 'b':
		square = 1
	case 'c':
		square = 2
	case 'd':
		square = 3
	case 'e':
		square = 4
	case 'f':
		square = 5
	case 'g':
		square = 6
	case 'h':
		square = 7
	default:
		return 0
	}

	return square + 8*(8-int(s[1]-'0'))
}

func squareToAlgebraic(s int) string {

	alg := ""
	switch s % 8 {
	case 0:
		alg += "a"
	case 1:
		alg += "b"
	case 2:
		alg += "c"
	case 3:
		alg += "d"
	case 4:
		alg += "e"
	case 5:
		alg += "f"
	case 6:
		alg += "g"
	case 7:
		alg += "h"
	}

	return alg + fmt.Sprint(8-(s/8))
}

//return the rank of the square
func rank(square int) int {
	return 8 - square/8
}

//returns the number of the file of the square
func file(square int) int {
	return square%8 + 1
}

func opponent(col int) int {
	if col == WHITE {
		return BLACK
	}
	return WHITE
}

func generateKey() uint64 {
	seed = (6364136223846793005*seed + 1442695040888963407) % 0x1000000000000000
	return seed
}
