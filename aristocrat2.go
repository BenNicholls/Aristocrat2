package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//import "github.com/pkg/profile"

var game position
var table *hashTable
var usingHashtable bool

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	game = NewPosition("")
	game.Output()

	//command line mode
	cli := CLI{}
	cli.gameLoop()
}

type CLI struct {
	command string
}

func (cli *CLI) gameLoop() {
	commandReader := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Aristocrat 2")
	for quit := false; quit != true; commandReader.Scan() {
		cmd := strings.SplitN(strings.TrimSpace(commandReader.Text()), " ", 2)
		switch strings.ToLower(cmd[0]) {
		case "quit":
			quit = true
		case "new":
			game = NewPosition("")
		case "setboard":
			game = NewPosition(cmd[1])
		case "display":
			game.Output()
		case "divide":
			if len(cmd) == 2 {
				plys, err := strconv.Atoi(cmd[1])
				if err == nil {
					game.divide(plys)
				} else {
					fmt.Println("Divide command argument must be integer")
				}
			} else {
				fmt.Println("Divide command must have single argument (# of plys)")
			}
		case "perft":
			if len(cmd) == 2 {
				plys, err := strconv.Atoi(cmd[1])
				if err == nil {
					startTime := time.Now()
					nodes := multiThreadedPerft(&game, plys)
					fmt.Printf("Perft %d: %d. (%.0f n/s)\n", plys, nodes, float64(nodes)/time.Since(startTime).Seconds())
				} else {
					fmt.Println("Perft command argument must be integer")
				}
			} else {
				fmt.Println("Perft command must have single argument (# of plys)")
			}
		case "move":
			if len(cmd) == 2 {
				from, to := algebraicToSquare(cmd[1][:2]), algebraicToSquare(cmd[1][2:4])
				list := movegen(&game)
				for _, m := range list {
					if m.from() == from && m.to() == to {
						game.doMove(m)
						game.Output()
					}
				}
					}
		case "search":
			if len(cmd) == 2 {
				depth, err := strconv.Atoi(cmd[1])
				if err != nil {
					fmt.Println("search command arhument must be integer")
				}
				startTime := time.Now()
				score, nodes, best := search(&game, depth, -10000000, 10000000)
				fmt.Printf("Eval: %.2f | Move: %s\n", float64(score)/100, best.string())
				dur := time.Since(startTime).Seconds()
				fmt.Printf("searched %d nodes in %.3fs (%.0f nps)\n", nodes, dur, float64(nodes)/dur)
			}
		}

		if quit {
			break
		} else {
			fmt.Print(">>> ")
		}
	}
}

func multiThreadedPerft(pos *position, plys int) (nodes int) {
	list := movegen(pos)
	if plys <= 0 {
		return 0
	} else if plys == 1 {
		return len(list)
	}
	results := make(chan int, len(list))
	for _, m := range list {
		nextPosition := *pos
		nextPosition.doMove(m)
		go func() {
			results <- nextPosition.perft(plys - 1)
		}()
	}
	for i := 0; i < len(list); i++ {
		nodes += <-results
	}

	return
}
