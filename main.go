// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

import "fmt"

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

	ParseFEN(START_FEN, board)

	move1 := ParseMove("e2e4", board)
	StorePvMove(board, move1)
	MakeMove(board, move1)

	move2 := ParseMove("e7e5", board)
	StorePvMove(board, move2)
	MakeMove(board, move2)

	move3 := ParseMove("g1f3", board)
	StorePvMove(board, move3)
	MakeMove(board, move3)

	TakeMove(board)
	TakeMove(board)
	TakeMove(board)

	PrintBoard(board)

	count := GetPvLine(5, board)

	fmt.Printf("\nPV line found, %d moves:\n", count)
	for i := range count {
		fmt.Printf("move %d: %s\n", i+1, PrMove(board.PvArray[i]))
	}

	PrintBoard(board)
}
