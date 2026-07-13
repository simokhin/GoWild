// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const PUZZLE_FEN = "1B4k1/P4rpp/q4p2/8/1p6/1Q2P3/3PKPPP/2r3R1 w - - 2 27"

// main is the program entry point. It initialises the board representation
// lookup tables, sets up the starting position, and exercises the PV table:
// it plays a short sequence of moves in algebraic notation (e.g., "e2e4"),
// storing each as the PV move for the position it was played from, then
// unwinds all of them back to the starting position and uses GetPvLine to
// recover the line from the PV table alone.
func main() {
	AllInit()

	board := &Board{}
	board.PvTable = &PVTable{}
	InitPvTable(board.PvTable)

	info := &SearchInfo{}

	ParseFEN(PUZZLE_FEN, board)

	reader := bufio.NewReader(os.Stdin)

	for {
		PrintBoard(board)
		fmt.Print("Please enter a move > ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if len(input) == 0 {
			continue
		}

		if input[0] == 'q' {
			break
		} else if input[0] == 't' {
			TakeMove(board)
		} else if input[0] == 's' {
			info.Depth = 4
			SearchPosition(board, info)
		} else {
			move := ParseMove(input, board)
			if move != NoMove {
				StorePvMove(board, move)
				MakeMove(board, move)

			} else {
				fmt.Printf("Move Not Parsed: %s\n", input)
			}
		}
	}
}
