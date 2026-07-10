package main

import "fmt"

// Piece represents a chess piece type.
// The values are used internally in the 120-square mailbox board.
type Piece int8

const (
	Empty Piece = iota // Empty square
	WP                 // White Pawn
	WN                 // White Knight
	WB                 // White Bishop
	WR                 // White Rook
	WQ                 // White Queen
	WK                 // White King
	BP                 // Black Pawn
	BN                 // Black Knight
	BB                 // Black Bishop
	BR                 // Black Rook
	BQ                 // Black Queen
	BK                 // Black King
)

// File represents a file (column) on the chessboard, from A (left) to H (right).
type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	FileNone
)

// Rank represents a rank (row) on the chessboard, from 1 (white's home) to 8 (black's home).
type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	RankNone
)

// Color represents the side to move or the owner of a piece.
type Color int8

const (
	White Color = iota
	Black
	Both // Used for bitboards that track both colours together
)

// Square is a 0x88-style mailbox board index (10x12 = 120 squares).
// The border squares (ranks 0, 9 and files 0, 9) are off-board padding
// used to quickly detect moves that leave the board.
type Square int8

const (
	A1 Square = 21
	B1 Square = 22
	C1 Square = 23
	D1 Square = 24
	E1 Square = 25
	F1 Square = 26
	G1 Square = 27
	H1 Square = 28

	A2 Square = 31
	B2 Square = 32
	C2 Square = 33
	D2 Square = 34
	E2 Square = 35
	F2 Square = 36
	G2 Square = 37
	H2 Square = 38

	A3 Square = 41
	B3 Square = 42
	C3 Square = 43
	D3 Square = 44
	E3 Square = 45
	F3 Square = 46
	G3 Square = 47
	H3 Square = 48

	A4 Square = 51
	B4 Square = 52
	C4 Square = 53
	D4 Square = 54
	E4 Square = 55
	F4 Square = 56
	G4 Square = 57
	H4 Square = 58

	A5 Square = 61
	B5 Square = 62
	C5 Square = 63
	D5 Square = 64
	E5 Square = 65
	F5 Square = 66
	G5 Square = 67
	H5 Square = 68

	A6 Square = 71
	B6 Square = 72
	C6 Square = 73
	D6 Square = 74
	E6 Square = 75
	F6 Square = 76
	G6 Square = 77
	H6 Square = 78

	A7 Square = 81
	B7 Square = 82
	C7 Square = 83
	D7 Square = 84
	E7 Square = 85
	F7 Square = 86
	G7 Square = 87
	H7 Square = 88

	A8 Square = 91
	B8 Square = 92
	C8 Square = 93
	D8 Square = 94
	E8 Square = 95
	F8 Square = 96
	G8 Square = 97
	H8 Square = 98

	NoSquare Square = 99 // Sentinel value meaning "no square" (e.g., no en passant target)
)

// CastlePerm is a bitmask of which castling rights remain available.
type CastlePerm int8

const (
	WKCA CastlePerm = 1 // White kingside castle available
	WQCA CastlePerm = 2 // White queenside castle available
	BKCA CastlePerm = 4 // Black kingside castle available
	BQCA CastlePerm = 8 // Black queenside castle available
)

// MaxGameMoves is the maximum number of half-moves we store in the history.
const MaxGameMoves = 2048

// Undo stores the board state before a move, so we can take it back.
type Undo struct {
	Move       int        // The move that was played
	CastlePerm CastlePerm // Castling rights before the move
	EnPas      Square     // En passant square before the move
	FiftyMove  int        // 50-move rule counter before the move
	PosKey     uint64     // Zobrist hash key of the position before the move
}

// Bitboard is a 64-bit mask representing a set of squares on the board.
// Bit 0 = A1, bit 63 = H8.
type Bitboard uint64

// Board holds the complete state of a chess position.
// It uses a 120-square mailbox representation for fast move generation.
type Board struct {
	Pieces [120]Piece // Mailbox board (with border/padding), indexed by Square

	Pawns [3]Bitboard // White pawns, black pawns, both colours pawns

	KingSq [2]Square // White king, black king

	Side      Color  // Which side to move
	EnPas     Square // En passant square
	FiftyMove int    // 50 moves rule count

	Ply    int // Depth of current search
	HisPly int // Count of all half-moves (plies) made since game start

	CastlePerm CastlePerm // Bitmask of castling rights: which sides can castle which way

	PosKey uint64 // Zobrist hash key uniquely indentifying the current position

	PceNum [13]int // Count of each piece type on the board, indexed by Piece
	BigPce [3]int  // Count of non-pawn pieces (per side + both)
	MajPce [3]int  // Count of major pieces: rooks and queens (per side + both)
	MinPce [3]int  // Count of minor pieces: knights and bishops (per side + both)

	History []Undo

	PList [13][10]Square
}

// Translate file+rank (A1, B4 etc.) to square index
func FR2SQ(file File, rank Rank) Square {
	return Square(21 + int(file) + int(rank)*10)
}

var Sq120ToSq64 [120]int   // Maps a 120-square mailbox index to its 64-square index
var Sq64ToSq120 [64]Square // Maps a 64-square index back to its 120-square mailbox index

func SQ64(sq120 Square) int {
	return Sq120ToSq64[sq120]
}

const Debug = true

func Assert(condition bool, message string) {
	if !Debug {
		return
	}
	if !condition {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}
