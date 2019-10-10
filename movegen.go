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

func movegen(pos position) (list moveList) {
	list = make([]move, 0, 20)

	var pieces uint64
	occupied := pos.white | pos.black
	//PAWNS
	if pos.toMove == WHITE {
		pieces = pos.white
	} else {
		pieces = pos.black
	}

	//generate pawn moves
	pawns := pieces & pos.pieces[PAWN]
	forEachBit(pawns, func(pawnSquare int) {
		if pos.toMove == WHITE {
			moves := whitePawnMoves[pawnSquare] &^ occupied
			if rank(pawnSquare) == 2 && !checkBit(moves, pawnSquare-8) {
				moves = clearBit(moves, pawnSquare-16)
			}
			forEachBit(moves, func(square int) {
				if rank(square) == 8 {
					list = appendPromotionMoves(list, pawnSquare, square, WHITE, false)
				} else {
					list = append(list, packMove(pawnSquare, square, PAWN, 0, WHITE, false))
				}
			})

			captures := whitePawnAttacks[pawnSquare] & (setBit(pos.black, pos.enpassant))
			forEachBit(captures, func(square int) {
				if rank(square) == 8 {
					list = appendPromotionMoves(list, pawnSquare, square, WHITE, true)
				} else {
					list = append(list, packMove(pawnSquare, square, PAWN, 0, WHITE, true))
				}
			})
		} else {
			moves := blackPawnMoves[pawnSquare] &^ occupied
			if rank(pawnSquare) == 7 && !checkBit(moves, pawnSquare+8) {
				moves = clearBit(moves, pawnSquare+16)
			}
			forEachBit(moves, func(square int) {
				if rank(square) == 1 {
					list = appendPromotionMoves(list, pawnSquare, square, BLACK, false)
				} else {
					list = append(list, packMove(pawnSquare, square, PAWN, 0, BLACK, false))
				}
			})

			captures := blackPawnAttacks[pawnSquare] & (setBit(pos.white, pos.enpassant))
			forEachBit(captures, func(square int) {
				if rank(square) == 1 {
					list = appendPromotionMoves(list, pawnSquare, square, BLACK, true)
				} else {
					list = append(list, packMove(pawnSquare, square, PAWN, 0, BLACK, true))
				}
			})
		}
	})

	return
}

//hi future ben. You're wondering why this isn't a loop. Well apparently having it as a loop
//breaks the go compiler for reasons that are ungooglable. something about it trying to use
//invalid asm instructions under certain inputs. yeah. for some reason manually unrolling
//the loop works so don't touch it.
func appendPromotionMoves(list moveList, from, to, turn int, capture bool) moveList {
	list = append(list, packMove(from, to, PAWN, KNIGHT, turn, capture))
	list = append(list, packMove(from, to, PAWN, BISHOP, turn, capture))
	list = append(list, packMove(from, to, PAWN, ROOK, turn, capture))
	list = append(list, packMove(from, to, PAWN, QUEEN, turn, capture))

	return list
}
