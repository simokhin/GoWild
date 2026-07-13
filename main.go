// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// main is the program entry point. It initialises the board representation
// lookup tables, sets up the starting position, and enters an interactive
// move-input loop where the user can enter moves in algebraic notation
// (e.g., "e2e4") or type "quit" to exit.
func main() {
	AllInit()

	board := &Board{}
	ParseFEN(START_FEN, board)
	PrintBoard(board)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter move (or 'quit'): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" {
			break
		}

		move := ParseMove(input, board)
		if move == NoMove {
			fmt.Println("Invalid move, try again")
			continue
		}

		if !MakeMove(board, move) {
			fmt.Println("Illegal move (king would be in check), try again")
			continue
		}

		PrintBoard(board)
		fmt.Println("Poskey:", board.PosKey)
	}
}
