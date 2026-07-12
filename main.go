// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	board := &Board{}
	ParseFEN(START_FEN, board)

	PerftTest(4, board)
}
