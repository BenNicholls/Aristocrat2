package main

var scoreModifier [2]int = [2]int{1, -1}

func search(p *position, depth, alpha, beta int) (score, nodes int, bestMove move) {
	if entry, ok := table.Load(p.hash); ok {
		if entry.depth >= depth {
			return entry.score, 1, entry.bestMove
		}
	}
	if depth == 0 {
		return eval(p), 1, 0
	}

	moves := movegen(p)
	for _, m := range moves {
		next := *p
		next.doMove(m)
		e, n, _ := search(&next, depth-1, -beta, -alpha)
		e = e * scoreModifier[p.toMove]
		if e > alpha {
			alpha = e
			bestMove = m
		}
		nodes += n
		if alpha >= beta {
			break
		}
	}

	score = scoreModifier[p.toMove] * alpha
	table.Store(p.hash, depth, bestMove, score)
	return
}
