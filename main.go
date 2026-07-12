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

	board := &Board{}

	ParseFEN(FEN3, board)
	PrintBoard(board)

	// CheckBoard(board)

	fmt.Println("\n\nWhite Attacking:")
	ShowSqAtBySide(White, board)

	fmt.Println("\n\nBlack Attacking")
	ShowSqAtBySide(Black, board)
}

// ShowSqAtBySide prints an 8x8 board of X and - characters showing which squares
// are attacked by the given side. Used for visually debugging attack detection.
func ShowSqAtBySide(side Color, pos *Board) {
	var rank Rank
	var file File
	var sq Square

	fmt.Printf("\n\nSquares attacked by: %c\n", SideChar[side])

	for rank = Rank8; rank >= Rank1; rank-- {
		for file = FileA; file <= FileH; file++ {
			sq = FR2SQ(file, rank)
			if SqAttacked(sq, side, pos) {
				fmt.Print("X")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
}
