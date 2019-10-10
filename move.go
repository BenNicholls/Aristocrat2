package main

import "fmt"

//these magical numbers define the spec for the move data structure
const (
	//masks
	M_SPACEMASK = 0b111111
	M_PIECEMASK = 0b111

	//sizes (bits)
	M_SPACESIZE = 6
	M_PIECESIZE = 3

	//offsets
	M_FROMOFFSET         = 0
	M_TOOFFSET           = 6
	M_PIECEOFFSET        = 12
	M_PROMOTEPIECEOFFSET = 15

	//flags
	M_PROMOTEFLAG = (1 << 61)
	M_CAPTUREFLAG = (1 << 62)
	M_TURNFLAG    = (1 << 63)
)

type move uint64

//from space
func (m move) from() int {
	return int(M_SPACEMASK & (m << M_FROMOFFSET))
}

//to space
func (m move) to() int {
	return int(M_SPACEMASK & (m >> M_TOOFFSET))
}

//piece moving
func (m move) piece() int {
	return int(M_PIECEMASK & (m >> M_PIECEOFFSET))
}

//piece to promote to
func (m move) promotedPiece() int {
	return int(M_PIECEMASK & (m >> M_PROMOTEPIECEOFFSET))
}

//bit 61: promotion flag
func (m move) promote() bool {
	return M_PROMOTEFLAG&m != 0
}

//bit 62: capture flag.
func (m move) capture() bool {
	return M_CAPTUREFLAG&m != 0
}

//bit 63: turn. 0 = WHITE, 1 = BLACK
func (m move) turn() int {
	if M_TURNFLAG&m == 0 {
		return WHITE
	}
	return BLACK
}

func (m move) string() (s string) {
	s += pieceNamesShort[m.piece()] + squareToAlgebraic(m.from())
	if m.capture() {
		s += "x"
	}
	s += squareToAlgebraic(m.to())
	if m.promote() {
		s += pieceNamesShort[m.promotedPiece()]
	}
	return
}

func (m move) output() {
	fmt.Println(m.string())
	if m.turn() == WHITE {
		fmt.Println("Piece: White", pieceNames[m.piece()])
	} else {
		fmt.Println("Piece: Black", pieceNames[m.piece()])
	}
	fmt.Println("From", squareToAlgebraic(m.from()), "to", squareToAlgebraic(m.to()))
	if m.capture() {
		fmt.Println("It is a capture.")
	}
	if m.promote() {
		fmt.Println("Promoting to a ", pieceNames[m.promotedPiece()])
	}
}

func packMove(from, to, piece, promotePiece, turn int, capture bool) move {
	var m uint64

	//values (right side)
	m = uint64(promotePiece)
	m = (m << M_PIECESIZE) | uint64(piece)
	m = (m << M_SPACESIZE) | uint64(to)
	m = (m << M_SPACESIZE) | uint64(from)

	//flags (left side)
	if promotePiece != 0 {
		m = m | M_PROMOTEFLAG
	}
	if capture {
		m = m | M_CAPTUREFLAG
	}
	if turn == BLACK {
		m = m | M_TURNFLAG
	}

	return move(m)
}
