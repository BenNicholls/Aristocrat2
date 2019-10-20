package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type chessInterface interface {
	processCommand(cmd, params string) string //channel to let us know when program can end
	mode() string
}

type CLIinterface struct {
}

func (cli *CLIinterface) processCommand(cmd, params string) string {
	switch cmd {
	case "quit":
		return "quit"
	case "new":
		game = NewPosition("")
		game.Output()
	case "setboard":
		game = NewPosition(params)
		game.Output()
	case "display":
		game.Output()
	case "divide":
		if params != "" {
			plys, err := strconv.Atoi(params)
			if err == nil {
				game.divide(plys)
			} else {
				fmt.Println("Divide command argument must be integer")
			}
		} else {
			fmt.Println("Divide command must have single argument (# of plys)")
		}
	case "perft":
		if params != "" {
			plys, err := strconv.Atoi(params)
			if err == nil {
				startTime := time.Now()
				nodes := multiThreadedPerft(&game, plys)
				fmt.Printf("Perft %d: %d. (%s)\n", plys, nodes, nps(nodes, time.Since(startTime).Seconds()))
			} else {
				fmt.Println("Perft command argument must be integer")
			}
		} else {
			fmt.Println("Perft command must have single argument (# of plys)")
		}
	case "move":
		if params != "" {
			from, to := algebraicToSquare(params[:2]), algebraicToSquare(params[2:4])
			list := movegen(&game)
			for _, m := range list {
				if m.from() == from && m.to() == to {
					game.doMove(m)
					game.Output()
				}
			}
		}
	case "search":
		if params != "" {
			depth, err := strconv.Atoi(params)
			if err != nil {
				fmt.Println("search command argument must be integer")
			}
			go iterativeSearch(&game, depth)
		}
	case "uci":
		return "uci"
	case "stop":
		if calcController.calculating() {
			calcController.stopCalculators()
		} else {
			fmt.Println("Engine not calculating")
		}
	}

	fmt.Print(">>> ")
	return ""
}

func (cli *CLIinterface) mode() string {
	return "cli"
}

type UCIinterface struct {
}

func (uci *UCIinterface) processCommand(cmd, params string) string {
	switch cmd {
	case "uci":
		fmt.Println("id name Aristocrat")
		fmt.Println("id author Benjamin Nicholls")
		fmt.Println("option name Hash type spin default 2 min 0 max 512 ")
		fmt.Println("uciok")
	case "debug":
	case "isready":
		fmt.Println("readyok")
	case "setoption":
		option := strings.Split(params, " ")
		if option[0] != "name" {
			return ""
		}
		optionName := option[1]

		switch optionName {
		case "Hash":
			if len(option) == 4 {
				size, err := strconv.Atoi(option[3])
				if err != nil {
					return ""
				}
				hashSize = size
				initHashTable()
				fmt.Println("info string Hashtable set to ", size, "MB")
			}
		}
	case "ucinewgame":
	case "position":
		s := strings.Split(strings.TrimPrefix(params, "fen "), " moves ")
		game = NewPosition(s[0])
		if len(s) == 2 {
			moves := strings.Split(s[1], " ")
			for _, m := range moves {
				from, to := algebraicToSquare(m[:2]), algebraicToSquare(m[2:4])
				var p int
				if len(m) == 5 {
					p = displayLookup[rune(m[4])].piece
				}
				list := movegen(&game)
				for _, m := range list {
					if m.from() == from && m.to() == to {
						if !m.promote() || m.promotedPiece() == p {
							game.doMove(m)
						}
					}
				}
			}
		}
	case "go":
		paramScanner := bufio.NewScanner(strings.NewReader(params))
		paramScanner.Split(bufio.ScanWords)
		var infinite bool
		var depth, time int
		for scan := paramScanner.Scan(); scan; scan = paramScanner.Scan() {
			switch paramScanner.Text() {
			case "infinite":
				infinite = true
			case "depth":
				if paramScanner.Scan() {
					if d, err := strconv.Atoi(paramScanner.Text()); err == nil {
						depth = d
					} else {
						fmt.Println("info string ERROR: could not parse depth value. (not an int?)")
					}
				} else {
					fmt.Println("info string ERROR: could not parse depth. (no value?)")
				}
			case "movetime":
				if paramScanner.Scan() {
					if t, err := strconv.Atoi(paramScanner.Text()); err == nil {
						time = t
					} else {
						fmt.Println("info string ERROR: could not parse time value. (not an int?)")
					}
				} else {
					fmt.Println("info string ERROR: could not parse time. (no value?)")
				}
			}
		}
		if infinite {
			go iterativeSearch(&game, 100)
		} else if depth != 0 {
			go iterativeSearch(&game, depth)
		} else if time != 0 {
			calcController.timeForMove = time
			go iterativeSearch(&game, 100)
		} else {
			fmt.Println("info string Not sure what to do here. Search for 8 plys I guess?")
			go iterativeSearch(&game, 8)
		}
	case "stop":
		calcController.stopCalculators()
	case "ponderhit":
	case "quit":
		return "quit"
	case "cli": //return to command line mode
		return "cli"
	}

	return ""
}

func (uci UCIinterface) mode() string {
	return "uci"
}
