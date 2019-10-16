package main

// import "fmt"

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

//directions. NOTE: the order here is important for movegen.
const (
	LEFT int = iota
	UPLEFT
	UP
	UPRIGHT
	RIGHT
	DOWNRIGHT
	DOWN
	DOWNLEFT
)

//piece move bitboards
var pawnMoves [2][64]uint64
var pawnAttacks [2][64]uint64
var knightMoves [64]uint64
var kingMoves [64]uint64
var slidingMoves [8][64]uint64 //indexed by direction

type zobristKeys struct {
	black     uint64
	pieces    [2][6][64]uint64
	enpassant [8]uint64
	castle    [4]uint64
}

var zobrist zobristKeys

type piece struct {
	colour int
	piece  int
}

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
		if rank(i) != 8 {
			pawnMoves[WHITE][i] = setBit(pawnMoves[WHITE][i], i-8)
			if rank(i) == 2 {
				pawnMoves[WHITE][i] = setBit(pawnMoves[WHITE][i], i-16)
			}
			if file(i) != 1 {
				pawnAttacks[WHITE][i] = setBit(pawnAttacks[WHITE][i], i-9)
			}
			if file(i) != 8 {
				pawnAttacks[WHITE][i] = setBit(pawnAttacks[WHITE][i], i-7)
			}
		}
		if rank(i) != 1 {
			pawnMoves[BLACK][i] = setBit(pawnMoves[BLACK][i], i+8)
			if rank(i) == 7 {
				pawnMoves[BLACK][i] = setBit(pawnMoves[BLACK][i], i+16)
			}
			if file(i) != 1 {
				pawnAttacks[BLACK][i] = setBit(pawnAttacks[BLACK][i], i+7)
			}
			if file(i) != 8 {
				pawnAttacks[BLACK][i] = setBit(pawnAttacks[BLACK][i], i+9)
			}
		}
	}

	//knight move bitboards
	for i := 0; i < 64; i++ {
		if rank(i) <= 6 {
			if file(i) != 1 {
				knightMoves[i] = setBit(knightMoves[i], i-17)
			}
			if file(i) != 8 {
				knightMoves[i] = setBit(knightMoves[i], i-15)
			}
		}
		if rank(i) <= 7 {
			if file(i) >= 3 {
				knightMoves[i] = setBit(knightMoves[i], i-10)
			}
			if file(i) <= 6 {
				knightMoves[i] = setBit(knightMoves[i], i-6)
			}
		}
		if rank(i) >= 3 {
			if file(i) != 1 {
				knightMoves[i] = setBit(knightMoves[i], i+15)
			}
			if file(i) != 8 {
				knightMoves[i] = setBit(knightMoves[i], i+17)
			}
		}
		if rank(i) >= 2 {
			if file(i) >= 3 {
				knightMoves[i] = setBit(knightMoves[i], i+6)
			}
			if file(i) <= 6 {
				knightMoves[i] = setBit(knightMoves[i], i+10)
			}
		}
	}

	//king move boards
	for i := 0; i < 64; i++ {
		if rank(i) != 1 {
			kingMoves[i] = setBit(kingMoves[i], i+8)
			if file(i) != 1 {
				kingMoves[i] = setBit(kingMoves[i], i+7)
			}
			if file(i) != 8 {
				kingMoves[i] = setBit(kingMoves[i], i+9)
			}
		}
		if rank(i) != 8 {
			kingMoves[i] = setBit(kingMoves[i], i-8)
			if file(i) != 1 {
				kingMoves[i] = setBit(kingMoves[i], i-9)
			}
			if file(i) != 8 {
				kingMoves[i] = setBit(kingMoves[i], i-7)
			}
		}
		if file(i) != 1 {
			kingMoves[i] = setBit(kingMoves[i], i-1)
		}
		if file(i) != 8 {
			kingMoves[i] = setBit(kingMoves[i], i+1)
		}
	}

	//sliding piece move boards
	for i := 0; i < 64; i++ {
		offsets := []int{-1, -9, -8, -7, 1, 9, 8, 7}
		for dir, off := range offsets {
			ray := i
			for {
				if rank(ray) == 1 && (dir == DOWNLEFT || dir == DOWN || dir == DOWNRIGHT) {
					break
				}
				if rank(ray) == 8 && (dir == UPLEFT || dir == UP || dir == UPRIGHT) {
					break
				}
				if file(ray) == 1 && (dir == UPLEFT || dir == LEFT || dir == DOWNLEFT) {
					break
				}
				if file(ray) == 8 && (dir == UPRIGHT || dir == RIGHT || dir == DOWNRIGHT) {
					break
				}
				ray += off
				slidingMoves[dir][i] = setBit(slidingMoves[dir][i], ray)
			}
		}
	}

	zobrist = zobristKeys{}
	zobrist.black = generateKey()
	for i := 0; i < 64; i++ {
		for p := PAWN; p <= KING; p++ {
			zobrist.pieces[WHITE][p][i] = generateKey()
			zobrist.pieces[BLACK][p][i] = generateKey()
		}
		if i < 8 {
			zobrist.enpassant[i] = generateKey()
			if i < 4 {
				zobrist.castle[i] = generateKey()
			}
		}
	}
}
