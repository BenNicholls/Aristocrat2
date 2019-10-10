package main

import "fmt"
import "strings"

type position struct {
	toMove int

	castleWK bool
	castleWQ bool
	castleBK bool
	castleBQ bool

	enpassant int

	fiftyMoveCounter int
	fullMoveCounter  int

	//bitboards! so many bitboards!
	white  uint64
	black  uint64
	pieces [6]uint64 //one for each kind of piece
}

func (p position) Print() {
	boardString := make([]string, 64)
	for piece, board := range p.pieces {
		if board == 0 {
			continue
		}

		for i := 0; i < 64; i++ {
			if checkBit(board, i) {
				if checkBit(p.white, i) {
					boardString[i] = pieceNamesDisplay[WHITE][piece]
				} else if checkBit(p.black, i) {
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

}

func NewPosition(fen string) (pos position) {
	pos = position{}

	if fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	fenPieces := strings.Split(fen, " ")

	posString := strings.NewReader(fenPieces[0])
	square := 0
	for ch, _, err := posString.ReadRune(); err == nil; ch, _, err = posString.ReadRune() {
		p, ok := displayLookup[ch]
		if ok {
			pos.pieces[p.piece] = setBit(pos.pieces[p.piece], square)
			if p.colour == WHITE {
				pos.white = setBit(pos.white, square)
			} else {
				pos.black = setBit(pos.black, square)
			}
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

	if len(fenPieces) >= 5 {
		pos.fiftyMoveCounter = int(fenPieces[4][0] - '0')
	}

	if len(fenPieces) >= 6 {
		pos.fullMoveCounter = int(fenPieces[5][0] - '0')
	}

	return
}
