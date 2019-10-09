package main

import "fmt"

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
