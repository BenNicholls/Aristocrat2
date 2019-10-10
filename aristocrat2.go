package main

func main() {
	p := NewPosition("rn1qkb1r/p1P1p1p1/b4n1p/1p1pPp2/5P2/P1P5/2P3PP/RNBQKBNR w KQkq - 1 5")
	p.Print()

	ml := movegen(p)
	ml.output()
}

type piece struct {
	colour int
	piece  int
}
