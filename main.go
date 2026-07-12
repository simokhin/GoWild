// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

import "fmt"

// FEN1 is a test position: the Italian Game, Black to move after 1.e4 (en passant available on e3).
var FEN1 = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"

// FEN2 is a test position: Italian Game, White to move after 1.e4 c5 (en passant available on c6).
var FEN2 = "rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2"

// FEN3 is a test position: Italian Game with Nf3, Black to move after 1.e4 c5 2.Nf3.
var FEN3 = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"

// FEN4 is a test position: an asymmetric middlegame with both sides having castling rights,
// multiple pieces per side, and pawns on various files. Useful for testing CheckBoard assertions.
var FEN4 = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

// FEN5 is a test position: white queen on e4 vs. black queen on d7, no other pieces.
// Useful for testing attack detection between two queens on an otherwise empty board.
var FEN5 = "8/3q4/8/8/4Q3/8/8/8 w - - 0 2"

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	// Encode a test move by packing from-square, to-square, captured piece,
	// and promoted piece into a single 28-bit integer using bit shifts.
	// This validates the move encoding/decoding scheme used throughout the engine.
	move := 0
	from := 6
	to := 12
	cap := WR
	prom := BR

	move = from | (to << 7) | (int(cap) << 14) | (int(prom) << 20)

	fmt.Printf("\ndec:%d hex:%X\n", move, move)
	PrintBin(move)

	fmt.Printf("from:%d to:%d cap:%d prom:%d\n", FromSq(move), ToSq(move), Captured(move), Promoted(move))

}

// PrintBin prints the binary representation of a move integer, grouped
// in 4-bit nibbles for readability. Used for debugging move encoding.
func PrintBin(move int) {
	fmt.Println("As binary:")
	for index := 27; index >= 0; index-- {
		if (1<<index)&move != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		if index != 28 && index%4 == 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}
