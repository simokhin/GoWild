// Package gofish is a chess engine implementing a 120-square mailbox board
// representation with bitboard support, Zobrist hashing, and move generation.
package main

// PawnMovesW is a FEN test position featuring white pawns on various ranks with
// promotion, double-push, and en passant opportunities. Used to validate white
// pawn move generation (single pushes, double pushes, captures, promotions, and
// en passant captures).
const PawnMovesW = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"

// PawnMovesB is a FEN test position featuring black pawns on various ranks with
// promotion, double-push, and en passant opportunities. Used to validate black
// pawn move generation (single pushes, double pushes, captures, promotions, and
// en passant captures).
const PawnMovesB = "rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/RNBQKB1R b KQkq e3 0 1"

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	board := &Board{}

	ParseFEN(PawnMovesB, board)
	PrintBoard(board)

	list := &MoveList{}

	GenerateAllMoves(board, list)

	PrintMoveList(list)
}
