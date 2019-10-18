package main

import (
	"fmt"
)

type moveList []move

func (ml moveList) output() {
	for _, m := range ml {
		fmt.Println(m.string())
	}
	fmt.Println(len(ml), "total moves.")
}

//returns a string assuming the mvovelist is a variation
func (ml moveList) variation() (v string) {
	if len(ml) == 0 {
		return "No Moves."
	}
	for _, m := range ml {
		v += m.string() + " "
	}
	return
}

func movegen(pos *position) (list moveList) {
	list = make([]move, 0, 32)

	var pieces uint64
	var opponentPieces uint64
	occupied := pos.colours[WHITE] | pos.colours[BLACK]
	pieces = pos.colours[pos.toMove]
	opponentPieces = pos.colours[opponent(pos.toMove)]

	//generate pawn moves
	forEachBit(pieces&pos.pieces[PAWN], func(fromSquare int) {
		moves := pawnMoves[pos.toMove][fromSquare] &^ occupied
		if pos.toMove == WHITE && rank(fromSquare) == 2 && !checkBit(moves, fromSquare-8) {
			moves = clearBit(moves, fromSquare-16)
		} else if pos.toMove == BLACK && rank(fromSquare) == 7 && !checkBit(moves, fromSquare+8) {
			moves = clearBit(moves, fromSquare+16)
		}
		forEachBit(moves, func(toSquare int) {
			if (pos.toMove == WHITE && rank(toSquare) == 8) || (pos.toMove == BLACK && rank(toSquare) == 1) {
				addPromosToMovelist(pos, &list, fromSquare, toSquare, pos.toMove, false)
			} else {
				addToMovelist(pos, &list, packMove(fromSquare, toSquare, PAWN, 0, 0, pos.toMove, false))
			}
		})
		captures := pawnAttacks[pos.toMove][fromSquare] & (setBit(opponentPieces, pos.enpassant))
		forEachBit(captures, func(toSquare int) {
			if (pos.toMove == WHITE && rank(toSquare) == 8) || (pos.toMove == BLACK && rank(toSquare) == 1) {
				addPromosToMovelist(pos, &list, fromSquare, toSquare, pos.toMove, true)
			} else {
				if toSquare == pos.enpassant {
					addToMovelist(pos, &list, packMove(fromSquare, toSquare, PAWN, 0, PAWN, pos.toMove, true))
				} else {
					addToMovelist(pos, &list, packMove(fromSquare, toSquare, PAWN, 0, pos.getPieceOnSquare(toSquare), pos.toMove, true))
				}
			}
		})
	})

	//workin on them knight moves
	forEachBit(pieces&pos.pieces[KNIGHT], func(fromSquare int) {
		moves := knightMoves[fromSquare] &^ pieces
		forEachBit(moves, func(toSquare int) {
			if checkBit(opponentPieces, toSquare) {
				addToMovelist(pos, &list, packMove(fromSquare, toSquare, KNIGHT, 0, pos.getPieceOnSquare(toSquare), pos.toMove, true))
			} else {
				addToMovelist(pos, &list, packMove(fromSquare, toSquare, KNIGHT, 0, 0, pos.toMove, false))
			}
		})
	})

	//king moves
	fromSquare := leftBit(pieces & pos.pieces[KING])
	moves := kingMoves[fromSquare] &^ pieces
	forEachBit(moves, func(toSquare int) {
		if checkBit(opponentPieces, toSquare) {
			addToMovelist(pos, &list, packMove(fromSquare, toSquare, KING, 0, pos.getPieceOnSquare(toSquare), pos.toMove, true))
		} else {
			addToMovelist(pos, &list, packMove(fromSquare, toSquare, KING, 0, 0, pos.toMove, false))
		}
	})
	if !pos.isSquareAttacked(pos.getKingSquare(pos.toMove), opponent(pos.toMove)) { //can't castle out of check
		if pos.toMove == WHITE {
			if pos.castleWK {
				if occupied&0b110 == 0 {
					if !pos.isSquareAttacked(61, BLACK) && !pos.isSquareAttacked(62, BLACK) { //can't castle through or into check either
						list = append(list, packMove(fromSquare, 62, KING, 0, 0, pos.toMove, false))
					}
				}
			}
			if pos.castleWQ {
				if occupied&0b1110000 == 0 {
					if !pos.isSquareAttacked(58, BLACK) && !pos.isSquareAttacked(59, BLACK) { //can't castle through or into check either
						list = append(list, packMove(fromSquare, 58, KING, 0, 0, pos.toMove, false))
					}
				}
			}
		} else {
			if pos.castleBK {
				if occupied&(0b11<<57) == 0 {
					if !pos.isSquareAttacked(5, WHITE) && !pos.isSquareAttacked(6, WHITE) { //can't castle through or into check either
						list = append(list, packMove(fromSquare, 6, KING, 0, 0, pos.toMove, false))
					}
				}
			}
			if pos.castleBQ {
				if occupied&(0b111<<60) == 0 {
					if !pos.isSquareAttacked(2, WHITE) && !pos.isSquareAttacked(3, WHITE) { //can't castle through or into check either
						list = append(list, packMove(fromSquare, 2, KING, 0, 0, pos.toMove, false))
					}
				}
			}
		}
	}

	//bishop
	forEachBit(pieces&pos.pieces[BISHOP], func(fromSquare int) {
		var moves uint64
		for dir := UPLEFT; dir <= DOWNLEFT; dir += 2 {
			rayMoves := slidingMoves[dir][fromSquare]
			if rayMoves&occupied != 0 {
				var endSquare int
				if dir <= UPRIGHT {
					endSquare = rightBit(rayMoves & occupied)
				} else {
					endSquare = leftBit(rayMoves & occupied)
				}

				moves = moves | (slidingMoves[dir][fromSquare] &^ slidingMoves[dir][endSquare])
				moves = clearBit(moves, endSquare)

				if checkBit(opponentPieces, endSquare) {
					addToMovelist(pos, &list, packMove(fromSquare, endSquare, BISHOP, 0, pos.getPieceOnSquare(endSquare), pos.toMove, true))
				}
			} else {
				moves = moves | slidingMoves[dir][fromSquare]
			}
		}
		forEachBit(moves, func(toSquare int) {
			addToMovelist(pos, &list, packMove(fromSquare, toSquare, BISHOP, 0, 0, pos.toMove, false))
		})
	})

	//rook
	forEachBit(pieces&pos.pieces[ROOK], func(fromSquare int) {
		var moves uint64
		for dir := LEFT; dir <= DOWN; dir += 2 {
			rayMoves := slidingMoves[dir][fromSquare]
			if rayMoves&occupied != 0 {
				var endSquare int
				if dir <= UPRIGHT {
					endSquare = rightBit(rayMoves & occupied)
				} else {
					endSquare = leftBit(rayMoves & occupied)
				}

				moves = moves | (slidingMoves[dir][fromSquare] &^ slidingMoves[dir][endSquare])
				moves = clearBit(moves, endSquare)

				if checkBit(opponentPieces, endSquare) {
					addToMovelist(pos, &list, packMove(fromSquare, endSquare, ROOK, 0, pos.getPieceOnSquare(endSquare), pos.toMove, true))
				}
			} else {
				moves = moves | slidingMoves[dir][fromSquare]
			}
		}
		forEachBit(moves, func(toSquare int) {
			addToMovelist(pos, &list, packMove(fromSquare, toSquare, ROOK, 0, 0, pos.toMove, false))
		})
	})

	//queen
	forEachBit(pieces&pos.pieces[QUEEN], func(fromSquare int) {
		var moves uint64
		for dir := LEFT; dir <= DOWNLEFT; dir++ {
			rayMoves := slidingMoves[dir][fromSquare]
			if rayMoves&occupied != 0 {
				var endSquare int
				if dir <= UPRIGHT {
					endSquare = rightBit(rayMoves & occupied)
				} else {
					endSquare = leftBit(rayMoves & occupied)
				}

				moves = moves | (slidingMoves[dir][fromSquare] &^ slidingMoves[dir][endSquare])
				moves = clearBit(moves, endSquare)

				if checkBit(opponentPieces, endSquare) {
					addToMovelist(pos, &list, packMove(fromSquare, endSquare, QUEEN, 0, pos.getPieceOnSquare(endSquare), pos.toMove, true))
				}
			} else {
				moves = moves | slidingMoves[dir][fromSquare]
			}
		}
		forEachBit(moves, func(toSquare int) {
			addToMovelist(pos, &list, packMove(fromSquare, toSquare, QUEEN, 0, 0, pos.toMove, false))
		})
	})

	return
}

func addToMovelist(pos *position, ml *moveList, m move) {
	//do a temporary shuffle so we can check move legality
	pos.colours[pos.toMove] = clearBit(pos.colours[pos.toMove], m.from())
	pos.colours[pos.toMove] = setBit(pos.colours[pos.toMove], m.to())

	//special case for captures, need to temporarily remove the captured piece as well
	if m.capture() {
		if m.to() == pos.enpassant {
			var captureSquare int
			if pos.toMove == WHITE {
				captureSquare = pos.enpassant + 8
			} else {
				captureSquare = pos.enpassant - 8
			}
			pos.colours[opponent(pos.toMove)] = clearBit(pos.colours[opponent(pos.toMove)], captureSquare)
			defer func() { pos.colours[opponent(pos.toMove)] = setBit(pos.colours[opponent(pos.toMove)], captureSquare) }()
		} else {
			pos.colours[opponent(pos.toMove)] = clearBit(pos.colours[opponent(pos.toMove)], m.to())
			defer func() { pos.colours[opponent(pos.toMove)] = setBit(pos.colours[opponent(pos.toMove)], m.to()) }()
		}
	}

	if m.piece() == KING {
		if !pos.isSquareAttacked(m.to(), opponent(pos.toMove)) {
			*ml = append(*ml, m)
		}
	} else if !pos.isSquareAttacked(pos.getKingSquare(pos.toMove), opponent(pos.toMove)) {
		*ml = append(*ml, m)
	}

	//put everything back now that we're done :)
	pos.colours[pos.toMove] = clearBit(pos.colours[pos.toMove], m.to())
	pos.colours[pos.toMove] = setBit(pos.colours[pos.toMove], m.from())
}

//hi future ben. You're wondering why this isn't a loop. Well apparently having it as a loop
//breaks the go compiler for reasons that are ungooglable. something about it trying to use
//invalid asm instructions under certain inputs. yeah. for some reason manually unrolling
//the loop works so don't touch it.
func addPromosToMovelist(pos *position, list *moveList, from, to, turn int, capture bool) {
	capturePiece := 0
	if capture {
		capturePiece = pos.getPieceOnSquare(to)
	}
	addToMovelist(pos, list, packMove(from, to, PAWN, KNIGHT, capturePiece, turn, capture))
	addToMovelist(pos, list, packMove(from, to, PAWN, BISHOP, capturePiece, turn, capture))
	addToMovelist(pos, list, packMove(from, to, PAWN, ROOK, capturePiece, turn, capture))
	addToMovelist(pos, list, packMove(from, to, PAWN, QUEEN, capturePiece, turn, capture))
}
