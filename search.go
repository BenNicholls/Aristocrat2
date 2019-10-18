package main

var scoreModifier [2]int = [2]int{1, -1}

//special scores
const (
	MATE int = 1000000
)

type result int

const (
	none result = iota
	stalemate
	checkmate
	//insufficient material/tablebase losses/draws/whatever can be added here later
)

func search(p *position, depth, alpha, beta int, currentVariation moveList) (score, nodes int, result result, bestVariation moveList) {
	if entry, ok := table.Load(p.hash); ok {
		if entry.depth >= depth {
			currentVariation = append(currentVariation, entry.bestMove)
			return entry.score, 1, entry.result, currentVariation
		}
	}
	if depth == 0 {
		return eval(p), 1, none, currentVariation
	}

	moves := movegen(p)
	if len(moves) == 0 {
		if p.isSquareAttacked(p.getKingSquare(p.toMove), opponent(p.toMove)) {
			return -MATE * scoreModifier[p.toMove], 1, checkmate, currentVariation
		} else {
			return 0, 1, stalemate, currentVariation
		}
	}
	var next position
	bestVariation = make(moveList, len(currentVariation))
	copy(bestVariation, currentVariation)
	for _, m := range moves {
		next = *p
		next.doMove(m)
		currentVariation = append(currentVariation, m)
		e, n, r, bv := search(&next, depth-1, alpha, beta, currentVariation)
		currentVariation = currentVariation[:len(currentVariation)-1]
		if p.toMove == WHITE && e > alpha {
			alpha = e
			result = r
			bestVariation = append(bestVariation[:len(currentVariation)], bv[len(currentVariation):]...)
		} else if p.toMove == BLACK && e < beta {
			beta = e
			result = r
			bestVariation = append(bestVariation[:len(currentVariation)], bv[len(currentVariation):]...)
		}
		nodes += n
		if alpha >= beta {
			break
		}
	}

	if p.toMove == WHITE {
		score = alpha
	} else {
		score = beta
	}
	if len(bestVariation) > len(currentVariation) { //if we found a followup move, store it
		table.Store(p.hash, depth, bestVariation[len(currentVariation)], score, result)
	}

	return
}
