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

// KnightsKings is a FEN test position with knights and kings only, both sides.
// Used to validate knight and king move generation (non-sliding pieces).
const KnightsKings = "5k2/1n6/4n3/6N1/8/3N4/8/5K2 b - - 0 1"

// Rooks is a FEN test position featuring rooks alongside knights and kings.
// Used to validate rook (orthogonal slider) and castling move generation.
const Rooks = "6k1/8/5r2/8/1nR5/5N2/8/6K1 w - - 0 1"

// Queens is a FEN test position with queens, knights, and kings on the board.
// Used to validate queen (combined slider) move generation.
const Queens = "6k1/8/4nq2/8/1nQ5/5N2/1N6/6K1 b - - 0 1"

// Bishops is a FEN test position with bishops, knights, and kings.
// Used to validate bishop (diagonal slider) move generation.
const Bishops = "6k1/1b6/4n3/8/1n4B1/1B3N2/1N6/2b3K1 w - - 0 1"

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	board := &Board{}
	list := &MoveList{}

	ParseFEN(Bishops, board)
	GenerateAllMoves(board, list)

	PrintMoveList(list)
}
