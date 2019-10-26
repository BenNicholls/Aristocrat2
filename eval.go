package main

//piece values
var PIECEVALUES [2][5]int = [2][5]int{
	[5]int{100, 300, 300, 500, 900},
	[5]int{-100, -300, -300, -500, -900},
}

//returns the evaluation of the position in centipawns
func eval(p *position) (score int) {
	//material evaluation
	for colour := WHITE; colour <= BLACK; colour++ {
		for piece := PAWN; piece <= QUEEN; piece++ {
			score += countBits(p.colours[colour]&p.pieces[piece]) * PIECEVALUES[colour][piece]
		}
	}

	return scoreModifier[p.toMove] * score
}
