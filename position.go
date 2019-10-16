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

	hash uint64 //zobrist hash. generated at start, then incrementally updated.
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

	pos.hash = pos.generateZobristHash()

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

	fmt.Println("HASH:", p.hash)
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

//removes a piece from the position, updating bitboards and stuff as appropriate
func (p *position) removePiece(colour, piece, square int) {
	p.colours[colour] = clearBit(p.colours[colour], square)
	p.pieces[piece] = clearBit(p.pieces[piece], square)
	p.hash ^= zobrist.pieces[colour][piece][square]
}

func (p *position) addPiece(colour, piece, square int) {
	p.colours[colour] = setBit(p.colours[colour], square)
	p.pieces[piece] = setBit(p.pieces[piece], square)
	p.hash ^= zobrist.pieces[colour][piece][square]
}

func (p *position) doMove(m move) {
	p.fiftyMoveCounter++
	if p.toMove == BLACK {
		p.fullMoveCounter++
	}

	//update bitboards
	p.removePiece(p.toMove, m.piece(), m.from())

	if m.capture() {
		p.fiftyMoveCounter = 0
		if m.to() == p.enpassant {
			var captureSquare int
			if p.toMove == WHITE {
				captureSquare = p.enpassant + 8
			} else {
				captureSquare = p.enpassant - 8
			}
			p.removePiece(opponent(p.toMove), PAWN, captureSquare)
		} else {
			p.removePiece(opponent(p.toMove), m.capturePiece(), m.to())
		}
	}

	if m.promote() {
		p.addPiece(p.toMove, m.promotedPiece(), m.to())
	} else {
		p.pieces[m.piece()] = setBit(p.pieces[m.piece()], m.to())
		p.addPiece(p.toMove, m.piece(), m.to())
	}

	if m.piece() == PAWN {
		p.fiftyMoveCounter = 0
	}

	//update enpassant square
	if p.enpassant != -1 {
		p.hash ^= zobrist.enpassant[file(p.enpassant)-1]
	}
	if m.pawnJump() {
		if p.toMove == WHITE {
			p.enpassant = m.from() - 8
		} else {
			p.enpassant = m.from() + 8
		}
		p.hash ^= zobrist.enpassant[file(p.enpassant)-1]
	} else {
		p.enpassant = -1
	}

	//move rooks for castling
	if m.castleK() {
		p.addPiece(p.toMove, ROOK, m.from()+1)
		if p.toMove == WHITE {
			p.removePiece(WHITE, ROOK, 63)
		} else {
			p.removePiece(BLACK, ROOK, 7)
		}
	} else if m.castleQ() {
		p.addPiece(p.toMove, ROOK, m.from()-1)
		if p.toMove == WHITE {
			p.removePiece(WHITE, ROOK, 56)
		} else {
			p.removePiece(BLACK, ROOK, 0)
		}
	}

	//update castling availability
	if m.piece() == KING {
		if p.toMove == WHITE {
			if p.castleWK {
				p.castleWK = false
				p.hash ^= zobrist.castle[0]
			}
			if p.castleWQ {
				p.castleWQ = false
				p.hash ^= zobrist.castle[1]
			}
		} else {
			if p.castleBK {
				p.castleBK = false
				p.hash ^= zobrist.castle[2]
			}
			if p.castleBQ {
				p.castleBQ = false
				p.hash ^= zobrist.castle[3]
			}
		}
	} else if m.piece() == ROOK {
		if file(m.from()) == 1 {
			if p.toMove == WHITE && p.castleWQ {
				p.castleWQ = false
				p.hash ^= zobrist.castle[1]
			} else if p.toMove == BLACK && p.castleBQ {
				p.castleBQ = false
				p.hash ^= zobrist.castle[3]
			}
		} else if file(m.from()) == 8 {
			if p.toMove == WHITE && p.castleWK {
				p.castleWK = false
				p.hash ^= zobrist.castle[0]
			} else if p.toMove == BLACK && p.castleBK {
				p.castleBK = false
				p.hash ^= zobrist.castle[2]
			}
		}
	}

	//checking for rook captures that break castling
	if m.capturePiece() == ROOK {
		switch m.to() {
		case 0:
			if p.castleBQ {
				p.castleBQ = false
				p.hash ^= zobrist.castle[3]
			}
		case 7:
			if p.castleBK {
				p.castleBK = false
				p.hash ^= zobrist.castle[2]
			}
		case 56:
			if p.castleWQ {
				p.castleWQ = false
				p.hash ^= zobrist.castle[1]
			}
		case 63:
			if p.castleWK {
				p.castleWK = false
				p.hash ^= zobrist.castle[0]
			}
		}
	}

	p.moveHistory = append(p.moveHistory, m)
	p.toMove = opponent(p.toMove)
	p.hash ^= zobrist.black
}

func (p position) perft(n int) (nodes int) {
	if entry, ok := table.Load(p.hash); ok {
		if entry.depth == n {
			return entry.score
		}
	}
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

	table.Store(p.hash, n, 0, nodes)
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
		p := multiThreadedPerft(&nextPosition, n-1)
		total += p
		fmt.Println(m.string()+":", p)
	}

	fmt.Println("Total: ", total)
	return
}

func (p *position) generateZobristHash() (hash uint64) {
	//pieces
	for colour := WHITE; colour <= BLACK; colour++ {
		for piece := PAWN; piece <= KING; piece++ {
			forEachBit(p.colours[colour]&p.pieces[piece], func(square int) {
				hash ^= zobrist.pieces[colour][piece][square]
			})
		}
	}

	//castling
	if p.castleWK {
		hash ^= zobrist.castle[0]
	}
	if p.castleWQ {
		hash ^= zobrist.castle[1]
	}
	if p.castleBK {
		hash ^= zobrist.castle[2]
	}
	if p.castleBQ {
		hash ^= zobrist.castle[3]
	}

	//enpassant
	if p.enpassant >= 0 {
		hash ^= zobrist.enpassant[file(p.enpassant)-1]
	}

	//turn
	if p.toMove == BLACK {
		hash ^= zobrist.black
	}

	return
}
