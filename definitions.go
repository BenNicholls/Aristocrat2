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

//piece move bitboards
var whitePawnMoves [64]uint64
var whitePawnAttacks [64]uint64
var blackPawnMoves [64]uint64
var blackPawnAttacks [64]uint64

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

	//generate pawn move bitboards
	for i := 0; i < 64; i++ {
		//no pawns on back rows
		if rank(i) == 1 || rank(i) == 8 {
			continue
		}
		whitePawnMoves[i] = setBit(whitePawnMoves[i], i-8)
		blackPawnMoves[i] = setBit(blackPawnMoves[i], i+8)

		if rank(i) == 2 {
			whitePawnMoves[i] = setBit(whitePawnMoves[i], i-16)
		} else if rank(i) == 7 {
			blackPawnMoves[i] = setBit(blackPawnMoves[i], i+16)
		}

		if file(i) != 1 {
			whitePawnAttacks[i] = setBit(whitePawnAttacks[i], i-9)
			blackPawnAttacks[i] = setBit(blackPawnAttacks[i], i+7)
		}
		if file(i) != 8 {
			whitePawnAttacks[i] = setBit(whitePawnAttacks[i], i-7)
			blackPawnAttacks[i] = setBit(blackPawnAttacks[i], i+9)
		}
	}
}
