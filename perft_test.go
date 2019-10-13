package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type perftTest struct {
	fen    string
	depths []int
	vals   []int
}

func TestPerft(t *testing.T) {
	//compile perft tests
	tests := make([]perftTest, 0)

	testSuite, err := os.Open("test/perftSuite.epd")
	if err != nil {
		t.Error("Could not open perft test suite.")
		return
	}
	defer testSuite.Close()

	suiteReader := bufio.NewScanner(testSuite)
	for ok := suiteReader.Scan(); ok; ok = suiteReader.Scan() {
		line := strings.SplitN(strings.TrimSpace(suiteReader.Text()), ";", 2)
		valStrings := strings.Split(line[1], ";")
		vals := make([]int, 0)
		depths := make([]int, 0)
		for _, v := range valStrings {
			num, err := strconv.Atoi(strings.TrimSpace(v[3:]))
			if err != nil {
				t.Error("Bad conversion: " + v[3:])
			}
			vals = append(vals, num)
			depth, err := strconv.Atoi(strings.TrimSpace(v[1:2]))
			if err != nil {
				t.Error("Bad conversion: " + v[1:2])
			}
			depths = append(depths, depth)
		}
		tests = append(tests, perftTest{
			fen:    line[0],
			vals:   vals,
			depths: depths,
		})
	}

	for tNum, perft := range tests {
		pos := NewPosition(perft.fen)
		for i, val := range perft.vals {
			fmt.Printf("(%d/%d) %s | depth: %d, expecting %d ...", tNum+1, len(tests), perft.fen, perft.depths[i], val)
			n := multiThreadedPerft(&pos, perft.depths[i])
			if n != val {
				fmt.Println("NO!")
				t.Errorf("PERFT FAIL: %s, depth %d. Expected %d, got %d", perft.fen, perft.depths[i], val, n)
				break
			} else {
				fmt.Println("YES!")
			}
		}
	}
}
