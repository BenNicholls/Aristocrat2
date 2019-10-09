package main

const (
	WHITE int = 0
	BLACK int = 1
)

const (
	PAWN int = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

var pieceNames [6]string           //long form piece names
var pieceNamesShort [6]string      //piece names for chess notation
var pieceNamesDisplay [2][6]string //piece names for display. also the FEN shortform. indexed by colour first.
var displayLookup map[rune]piece

func init() {
	pieceNames = [6]string{"Pawn", "Knight", "Bishop", "Rook", "Queen", "King"}
	pieceNamesShort = [6]string{"", "N", "B", "R", "Q", "K"}
	pieceNamesDisplay = [2][6]string{{"P", "N", "B", "R", "Q", "K"}, {"p", "n", "b", "r", "q", "k"}}

	displayLookup = make(map[rune]piece)
	displayLookup['P'] = piece{WHITE, PAWN}
	displayLookup['N'] = piece{WHITE, KNIGHT}
	displayLookup['B'] = piece{WHITE, BISHOP}
	displayLookup['R'] = piece{WHITE, ROOK}
	displayLookup['Q'] = piece{WHITE, QUEEN}
	displayLookup['K'] = piece{WHITE, KING}
	displayLookup['p'] = piece{BLACK, PAWN}
	displayLookup['n'] = piece{BLACK, KNIGHT}
	displayLookup['b'] = piece{BLACK, BISHOP}
	displayLookup['r'] = piece{BLACK, ROOK}
	displayLookup['q'] = piece{BLACK, QUEEN}
	displayLookup['k'] = piece{BLACK, KING}
}
