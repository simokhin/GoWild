package main

// PceChar maps each internal Piece constant to its printable character.
// Index by Piece value: Empty('.') -> King('k').
// The order is: ., P, N, B, R, Q, K, p, n, b, r, q, k
var PceChar = ".PNBRQKpnbrqk"

// SideChar maps each Color constant to its single-character representation.
// Index by Color: White -> 'w', Black -> 'b', Both -> '-'.
var SideChar = "wb-"

// RankChar maps each Rank constant (0-7) to its chess notation digit (1-8).
var RankChar = "12345678"

// FileChar maps each File constant (0-7) to its chess notation letter (a-h).
var FileChar = "abcdefgh"

// PieceBig marks piece types that are non-pawn (minor or major pieces).
// Index by Piece: Empty(false), Pawn(false), Knight(true), Bishop(true),
// Rook(true), Queen(true), King(true), and the same for black pieces.
var PieceBig = [13]bool{
	false, false, true, true, true, true, true,
	false, true, true, true, true, true,
}

// PieceMaj marks major piece types: rooks and queens.
// Index by Piece: rooks and queens are true, all others false.
var PieceMaj = [13]bool{
	false, false, false, false, true, true, true,
	false, false, false, true, true, true,
}

// PieceMin marks minor piece types: knights and bishops.
// Index by Piece: knights and bishops are true, all others false.
var PieceMin = [13]bool{
	false, false, true, true, false, false, false,
	false, true, true, false, false, false,
}

// PieceVal maps each Piece constant to its material value in centipawns.
// Pawn=100, Knight=Bishop=325, Rook=550, Queen=1000, King=50000 (effectively infinite).
var PieceVal = [13]int{
	0, 100, 325, 325, 550, 1000, 50000,
	100, 325, 325, 550, 1000, 50000,
}

// PieceCol maps each Piece constant to the Color that owns it.
// White pieces map to White, black pieces to Black, Empty/OffBoard to Both.
var PieceCol = [13]Color{
	Both, White, White, White, White, White, White,
	Black, Black, Black, Black, Black, Black,
}
