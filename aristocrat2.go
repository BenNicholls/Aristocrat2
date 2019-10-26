package main

import (
	"bufio"
	"os"
	"strings"
)

//import "github.com/pkg/profile"

var game position

var table *hashTable
var usingHashtable bool

var engineMode chessInterface
var calcController calculationController

//engine options
var hashSize int = 2

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	game = NewPosition("")
	initHashTable()

	//command line mode
	engineMode = &UCIinterface{}

	commandReader := bufio.NewScanner(os.Stdin)

	for engineResult := ""; engineResult != "quit"; commandReader.Scan() {
		cmd := strings.SplitN(strings.TrimSpace(commandReader.Text()), " ", 2)
		params := ""
		if len(cmd) == 2 {
			params = cmd[1]
		}
		engineResult = engineMode.processCommand(strings.ToLower(cmd[0]), params)
		switch engineResult {
		case "quit":
			return
		case "uci":
			engineMode = &UCIinterface{}
			engineMode.processCommand("uci", "")
		case "cli":
			engineMode = &CLIinterface{}
			engineMode.processCommand("new", "")
		}
	}
}

func multiThreadedPerft(pos *position, plys int) (nodes int) {
	list, _ := movegen(pos)
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
