package main

import (
	"fmt"
	"strings"
)

type position struct {
	toMove int

	castleWK bool
	castleWQ bool
	castleBK bool
	castleBQ bool

	enpassant int

	fiftyMoveCounter int
	fullMoveCounter  int

	moveHistory []move

	//bitboards! so many bitboards!
	colours [2]uint64 //one for each colour. OR these together to get the occupied board
	pieces  [6]uint64 //one for each kind of piece
}

func NewPosition(fen string) (pos position) {
	pos = position{}
	pos.moveHistory = make([]move, 0, 20)

	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	fenPieces := strings.Split(strings.TrimSpace(fen), " ")

	posString := strings.NewReader(fenPieces[0])
	square := 0
	for ch, _, err := posString.ReadRune(); err == nil; ch, _, err = posString.ReadRune() {
		p, ok := displayLookup[ch]
		if ok {
			pos.pieces[p.piece] = setBit(pos.pieces[p.piece], square)
			pos.colours[p.colour] = setBit(pos.colours[p.colour], square)
			square++
		} else if ch == '/' {
			continue
		} else { //number indicating empty spaces
			square += int(ch - '0')
		}
	}

	if fenPieces[1] == "w" {
		pos.toMove = WHITE
	} else if fenPieces[1] == "b" {
		pos.toMove = BLACK
	}

	castleString := strings.NewReader(fenPieces[2])
	for ch, _, err := castleString.ReadRune(); err == nil; ch, _, err = castleString.ReadRune() {
		switch ch {
		case 'K':
			pos.castleWK = true
		case 'Q':
			pos.castleWQ = true
		case 'k':
			pos.castleBK = true
		case 'q':
			pos.castleBQ = true
		}
	}

	if fenPieces[3] != "-" {
		pos.enpassant = algebraicToSquare(fenPieces[3])
	} else {
		pos.enpassant = -1
	}

	if len(fenPieces) >= 5 && fenPieces[4] != "" {
		pos.fiftyMoveCounter = int(fenPieces[4][0] - '0')
	}

	if len(fenPieces) >= 6 {
		pos.fullMoveCounter = int(fenPieces[5][0] - '0')
	}

	return
}

func (p position) Print() {
	boardString := make([]string, 64)
	for piece, board := range p.pieces {
		if board == 0 {
			continue
		}

		for i := 0; i < 64; i++ {
			if checkBit(board, i) {
				if checkBit(p.colours[WHITE], i) {
					boardString[i] = pieceNamesDisplay[WHITE][piece]
				} else if checkBit(p.colours[BLACK], i) {
					boardString[i] = pieceNamesDisplay[BLACK][piece]
				}
			}
		}
	}

	for i, p := range boardString {
		if i%8 == 0 {
			fmt.Print("+---+---+---+---+---+---+---+---+\n|")
		}
		if p == "" {
			fmt.Print("   |")
		} else {
			fmt.Print(" ", p, " |")
		}
		if i%8 == 7 {
			fmt.Print("\n")
		}
	}
	fmt.Print("+---+---+---+---+---+---+---+---+\n")

	fmt.Println("Turn ", p.fullMoveCounter)
	if p.toMove == WHITE {
		fmt.Println("White to move.")
	} else {
		fmt.Println("Black to move.")
	}

	if p.castleWK && p.castleWQ {
		fmt.Println("White can castle to both sides.")
	} else if p.castleWK {
		fmt.Println("White can castle kingside.")
	} else if p.castleWQ {
		fmt.Println("White can castle queenside.")
	}

	if p.castleBK && p.castleBQ {
		fmt.Println("Black can castle to both sides.")
	} else if p.castleBK {
		fmt.Println("Black can castle kingside.")
	} else if p.castleBQ {
		fmt.Println("Black can castle queenside.")
	}

	if p.enpassant >= 0 {
		fmt.Println("Enpassant square: ", squareToAlgebraic(p.enpassant))
	}

	fmt.Println("50 Move Counter: ", p.fiftyMoveCounter)
	if p.toMove == WHITE && p.isSquareAttacked(leftBit(p.pieces[KING]&p.colours[WHITE]), BLACK) {
		fmt.Println("White king is in check!")
	} else if p.toMove == BLACK && p.isSquareAttacked(leftBit(p.pieces[KING]&p.colours[BLACK]), WHITE) {
		fmt.Println("Black king is in check!")
	}
}

//tests whether a space is attacked by the provided colour player
func (p *position) isSquareAttacked(square, col int) bool {
	//check non-sliding pieces first
	if pawnAttacks[opponent(col)][square]&p.pieces[PAWN]&p.colours[col] != 0 || knightMoves[square]&p.pieces[KNIGHT]&p.colours[col] != 0 || kingMoves[square]&p.pieces[KING]&p.colours[col] != 0 {
		return true
	}

	//diagonals
	diagonalsSliders := p.colours[col] & (p.pieces[BISHOP] | p.pieces[QUEEN])
	for dir := UPLEFT; dir <= DOWNLEFT; dir += 2 {
		rayMoves := slidingMoves[dir][square]
		if rayMoves&diagonalsSliders != 0 { //if bishop or queen is even on diagonal
			var endSquare int
			if dir <= UPRIGHT {
				endSquare = rightBit(rayMoves & (p.colours[WHITE] | p.colours[BLACK]))
			} else {
				endSquare = leftBit(rayMoves & (p.colours[WHITE] | p.colours[BLACK]))
			}

			if checkBit(diagonalsSliders, endSquare) {
				return true
			}
		}
	}

	//rank and files
	horizontalSliders := p.colours[col] & (p.pieces[ROOK] | p.pieces[QUEEN])
	for dir := LEFT; dir <= DOWN; dir += 2 {
		rayMoves := slidingMoves[dir][square]
		if rayMoves&horizontalSliders != 0 { //if rook or queen is even on diagonal
			var endSquare int
			if dir <= UPRIGHT {
				endSquare = rightBit(rayMoves & (p.colours[WHITE] | p.colours[BLACK]))
			} else {
				endSquare = leftBit(rayMoves & (p.colours[WHITE] | p.colours[BLACK]))
			}

			if checkBit(horizontalSliders, endSquare) {
				return true
			}
		}
	}

	return false
}

func (p *position) getPieceOnSquare(square int) int {
	for i := PAWN; i <= KING; i++ {
		if checkBit(p.pieces[i], square) {
			return i
		}
	}
	return -1
}

//returns to square the king is on
func (p *position) getKingSquare(col int) int {
	return leftBit(p.pieces[KING] & p.colours[col])
}

func (p *position) doMove(m move) {
	p.fiftyMoveCounter++
	if p.toMove == BLACK {
		p.fullMoveCounter++
	}

	//update bitboards
	p.colours[p.toMove] = clearBit(p.colours[p.toMove], m.from())
	p.colours[p.toMove] = setBit(p.colours[p.toMove], m.to())

	p.pieces[m.piece()] = clearBit(p.pieces[m.piece()], m.from())

	if m.capture() {
		p.fiftyMoveCounter = 0
		p.colours[opponent(p.toMove)] = clearBit(p.colours[opponent(p.toMove)], m.to())
		p.pieces[m.capturePiece()] = clearBit(p.pieces[m.capturePiece()], m.to())
	}
	if m.promote() {
		p.pieces[m.promotedPiece()] = setBit(p.pieces[m.promotedPiece()], m.to())
	} else {
		p.pieces[m.piece()] = setBit(p.pieces[m.piece()], m.to())
	}

	if m.piece() == PAWN {
		p.fiftyMoveCounter = 0
		//capture enpassant
		if m.to() == p.enpassant {
			var captureSquare int
			if p.toMove == WHITE {
				captureSquare = p.enpassant + 8
			} else {
				captureSquare = p.enpassant - 8
			}
			p.colours[opponent(p.toMove)] = clearBit(p.colours[opponent(p.toMove)], captureSquare)
			p.pieces[PAWN] = clearBit(p.pieces[PAWN], captureSquare)
		}
	}

	//update enpassant square
	if m.pawnJump() {
		if p.toMove == WHITE {
			p.enpassant = m.from() - 8
		} else {
			p.enpassant = m.from() + 8
		}
	} else {
		p.enpassant = -1
	}

	//move rooks for castling
	if m.castleK() {
		p.pieces[ROOK] = setBit(p.pieces[ROOK], m.from()+1)
		p.colours[p.toMove] = setBit(p.colours[p.toMove], m.from()+1)
		if p.toMove == WHITE {
			p.pieces[ROOK] = clearBit(p.pieces[ROOK], 63)
			p.colours[WHITE] = clearBit(p.colours[WHITE], 63)
		} else {
			p.pieces[ROOK] = clearBit(p.pieces[ROOK], 7)
			p.colours[BLACK] = clearBit(p.colours[BLACK], 7)
		}
	} else if m.castleQ() {
		p.pieces[ROOK] = setBit(p.pieces[ROOK], m.from()-1)
		p.colours[p.toMove] = setBit(p.colours[p.toMove], m.from()-1)
		if p.toMove == WHITE {
			p.pieces[ROOK] = clearBit(p.pieces[ROOK], 56)
			p.colours[WHITE] = clearBit(p.colours[WHITE], 56)
		} else {
			p.pieces[ROOK] = clearBit(p.pieces[ROOK], 0)
			p.colours[BLACK] = clearBit(p.colours[BLACK], 0)
		}
	}

	//update castling availability
	if m.piece() == KING {
		if p.toMove == WHITE {
			p.castleWK = false
			p.castleWQ = false
		} else {
			p.castleBK = false
			p.castleBQ = false
		}
	} else if m.piece() == ROOK {
		if file(m.from()) == 1 {
			if p.toMove == WHITE {
				p.castleWQ = false
			} else {
				p.castleBQ = false
			}
		} else if file(m.from()) == 8 {
			if p.toMove == WHITE {
				p.castleWK = false
			} else {
				p.castleBK = false
			}
		}
	}
	switch m.to() { //checking for rook captures that break castling
	case 0:
		p.castleBQ = false
	case 7:
		p.castleBK = false
	case 56:
		p.castleWQ = false
	case 63:
		p.castleWK = false
	}

	if p.toMove == WHITE {
		p.toMove = BLACK
	} else {
		p.toMove = WHITE
	}

	p.moveHistory = append(p.moveHistory, m)
}

func (p position) perft(n int) (nodes int) {
	list := movegen(&p)
	if n == 1 {
		return len(list)
	} else {
		for _, m := range list {
			nextPosition := p
			nextPosition.doMove(m)
			nodes += nextPosition.perft(n - 1)
		}
	}

	return
}

func (p position) divide(n int) {
	if n == 1 {
		return
	}
	list := movegen(&p)
	total := 0
	for _, m := range list {
		nextPosition := p
		nextPosition.doMove(m)
		p := nextPosition.perft(n - 1)
		total += p
		fmt.Println(m.string()+":", p)
	}

	fmt.Println("Total: ", total)
	return
}
