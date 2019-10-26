package main

import (
	"fmt"
	"sync"
	"time"
)

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

type calculationController struct {
	sync.RWMutex
	calculators int  //number of threads currently calculating
	stop        bool //set to true to stop calculators
	timeForMove int  //time limit for move (in msec)
	timer       time.Time
}

func (cc *calculationController) calculating() bool {
	cc.RLock()
	defer cc.RUnlock()
	if cc.calculators > 0 {
		return true
	}
	return false
}

//Hi Future Ben, past Ben here. don't put the RUnlcok in a defer, it's really slow for some reason.
func (cc *calculationController) needToStop() bool {
	cc.RLock()
	if cc.stop || (cc.timeForMove != 0 && int64(cc.timeForMove)-time.Since(cc.timer).Milliseconds() < 10) {
		cc.RUnlock()
		return true
	}
	cc.RUnlock()
	return false
}

func (cc *calculationController) stopCalculators() {
	cc.Lock()
	cc.stop = true
	cc.Unlock()
}

func (cc *calculationController) beginCalculating() {
	cc.Lock()
	if cc.calculators == 0 && cc.timeForMove != 0 {
		cc.timer = time.Now()
	}
	cc.calculators++
	cc.Unlock()
}

func (cc *calculationController) doneCalculating() {
	cc.Lock()
	cc.calculators -= 1
	if cc.calculators == 0 {
		cc.stop = false
		cc.timeForMove = 0
	}
	cc.Unlock()
}

func search(p *position, depth, alpha, beta int, currentVariation moveList) (score, nodes int, result result, bestVariation moveList) {
	var candidateMove move
	if entry, ok := table.Load(p.hash); ok {
		candidateMove = entry.bestMove
		if entry.depth >= depth {
			if entry.node == EXACT {
				currentVariation = append(currentVariation, entry.bestMove)
				return entry.score, 1, entry.result, currentVariation
			} else if entry.node == LOWER {
				alpha = entry.score
				if alpha >= beta { //beta cutoff.
					currentVariation = append(currentVariation, entry.bestMove)
					return entry.score, 1, entry.result, currentVariation
				}
			}
		}
	}

	moves, numCaptures := movegen(p)
	if len(moves) == 0 {
		if p.isSquareAttacked(p.getKingSquare(p.toMove), opponent(p.toMove)) {
			return -MATE, 1, checkmate, currentVariation
		} else {
			return 0, 1, stalemate, currentVariation
		}
	}

	var quiesce bool
	if depth <= 0 {
		if numCaptures == 0 {
			return eval(p), 1, none, currentVariation
		}

		stand := eval(p)
		if stand > beta {
			return eval(p), 1, none, currentVariation
		}
		if stand > alpha {
			alpha = stand
		}
		quiesce = true
	}

	//find candidate move and place in front
	if candidateMove != 0 {
		for i, m := range moves {
			if m == candidateMove {
				if i != 0 {
					moves[i] = moves[0]
					moves[0] = candidateMove
					if moves[i].capture() && i > numCaptures {
						//if the replaced move is a capture, make sure it's with the captures at the start of the movelist
						swap := moves[numCaptures]
						moves[numCaptures] = moves[i]
						moves[i] = swap
					}
				}
				break
			}
		}
	}

	var next position
	score = -MATE * 2
	continuation := make(moveList, 0, 10)
	for _, m := range moves {
		if quiesce && (!m.capture() && m != candidateMove) {
			break
		}
		next = *p
		next.doMove(m)
		currentVariation = append(currentVariation, m)
		e, n, r, bv := search(&next, depth-1, -beta, -alpha, currentVariation)
		currentVariation = currentVariation[:len(currentVariation)-1]
		e = -e
		nodes += n
		if e > score {
			score = e
			result = r
		}
		if e > alpha {
			alpha = e
			continuation = continuation[:0]
			continuation = append(continuation, bv[len(currentVariation):]...)
		}
		if alpha >= beta || calcController.needToStop() {
			break
		}
	}

	bestVariation = make(moveList, len(currentVariation))
	copy(bestVariation, currentVariation)
	if len(continuation) != 0 { //if we found a followup move, add it to the PV and store it
		bestVariation = append(bestVariation, continuation...)
		if alpha >= beta { //beta cutoff node
			table.Store(p.hash, depth, continuation[0], score, result, LOWER)
		} else {
			table.Store(p.hash, depth, continuation[0], score, result, EXACT)
		}
	} else { //no improving move found.
		table.Store(p.hash, depth, move(0), score, result, UPPER)
	}

	return
}

func iterativeSearch(p *position, targetDepth int) (score, nodes int, result result, bestVariation moveList) {
	totalNodes := 0
	startTime := time.Now()
	calcController.beginCalculating()
	for depth := 1; depth <= targetDepth; depth++ {
		score, nodes, result, bestVariation = search(p, depth, -MATE*2, MATE*2, make(moveList, 0, 10))
		score *= scoreModifier[p.toMove]
		totalNodes += nodes
		if engineMode.mode() == "uci" {
			fmt.Printf("info depth %d score cp %d nodes %d nps %.0f pv %s\n", depth, score, totalNodes, float64(totalNodes)/time.Since(startTime).Seconds(), bestVariation.variation())
		} else if engineMode.mode() == "cli" {
			fmt.Print(depth, " | ")
			switch result {
			case checkmate:
				if score == MATE {
					fmt.Print("White is mating")
				} else {
					fmt.Print("Black is mating")
				}
			case stalemate:
				fmt.Print("Stalemate")
			default:
				fmt.Printf("Eval: %.2f", float64(score)/100)
			}
			fmt.Printf(" | Variation: %s\n", bestVariation.variation())
			dur := time.Since(startTime).Seconds()
			fmt.Printf("searched %d nodes in %.3fs (%s)\n", totalNodes, dur, nps(totalNodes, time.Since(startTime).Seconds()))
		}
		if calcController.needToStop() {
			break
		}
	}
	if engineMode.mode() == "uci" {
		fmt.Println("bestmove", bestVariation[0].UCIstring())
	}
	calcController.doneCalculating()

	return
}
