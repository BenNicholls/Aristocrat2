package main

func main() {
	p := NewPosition("")
	p.Print()
}

type piece struct {
	colour int
	piece  int
}
