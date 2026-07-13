// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const TEST_FEN = "r1b1k2r/ppppnppp/2n2q2/2b5/3NP3/2P1B3/PP3PPP/RN1QKB1R w KQkq - 0 1"

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

	ParseFEN(TEST_FEN, board)

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
			info.Depth = 5
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
