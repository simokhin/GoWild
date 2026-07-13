package main

import "fmt"

// ---- Constants ----

// START_FEN is the Forsyth–Edwards Notation string for the standard initial chess position.
const START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// MaxGameMoves is the maximum number of half-moves we store in the history.
const MaxGameMoves = 2048

// MaxPositionMoves is the maximum number of moves per position in the move list.
const MaxPositionMoves = 256

// Debug controls whether assertions are compiled in. Set to false for release builds.
const Debug = true

// Move encoding bit flags and masks.
const (
	MFlagEP   = 0x40000   // En passant capture flag (bit 18)
	MFlagPS   = 0x80000   // Pawn start / double-push flag (bit 19)
	MFlagCA   = 0x1000000 // Castle flag (bit 24)
	MFlagCap  = 0x7C000   // Mask for the captured piece field (bits 14–17)
	MFlagProm = 0xF00000  // Mask for the promoted piece field (bits 20–23)
)

// MaxDepth is the maximum search ply, and thus the maximum length of a PV line.
const MaxDepth = 64

// ---- Types ----

// Piece represents a chess piece type.
// The values are used internally in the 120-square mailbox board.
type Piece int8

// PVEntry is a single slot in the principal variation table, mapping a
// position's Zobrist key to the best move found from that position.
type PVEntry struct {
	PosKey uint64 // Zobrist hash key of the position this entry belongs to
	Move   int    // Best move found for that position
}

// PVTable is a hash table of PVEntry slots keyed by position hash, used to
// recover the principal variation after a search completes.
type PVTable struct {
	PTable []PVEntry
}

// NumEntries returns the number of slots in the PV table.
func (t *PVTable) NumEntries() int {
	return len(t.PTable)
}

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
	OffBoard
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

// Bitboard is a 64-bit mask representing a set of squares on the board.
// Bit 0 = A1, bit 63 = H8.
type Bitboard uint64

// Move represents a chess move encoded as a 28-bit integer (MoveInt) paired with a
// Score used for move ordering during search (e.g., MVV-LVA or history heuristic).
type Move struct {
	MoveInt int
	Score   int
}

// Undo stores the board state before a move, so we can take it back.
type Undo struct {
	Move       int        // The move that was played
	CastlePerm CastlePerm // Castling rights before the move
	EnPas      Square     // En passant square before the move
	FiftyMove  int        // 50-move rule counter before the move
	PosKey     uint64     // Zobrist hash key of the position before the move
}

// Board holds the complete state of a chess position.
// It uses a 120-square mailbox representation for fast move generation.
type Board struct {
	Pieces [120]Piece // Mailbox board (with border/padding), indexed by Square

	Pawns [3]Bitboard // White pawns, black pawns, both colours pawns

	KingSq [2]Square // White king, black king

	Side      Color  // Which side to move
	EnPas     Square // En passant square
	FiftyMove int    // 50 moves rule count

	Ply        int        // Depth of current search
	CastlePerm CastlePerm // Bitmask of castling rights: which sides can castle which way

	PosKey uint64 // Zobrist hash key uniquely identifying the current position

	PceNum   [13]int // Count of each piece type on the board, indexed by Piece
	BigPce   [2]int  // Count of non-pawn pieces (per side)
	MajPce   [2]int  // Count of major pieces: rooks and queens (per side)
	MinPce   [2]int  // Count of minor pieces: knights and bishops (per side)
	Material [2]int  // Total material value of pieces in centipawns (per side)

	History []Undo // Move history stack for undoing moves (stores Undo snapshots)

	PList [13][10]Square // Piece list: for each piece type (13), up to 10 squares where that piece sits

	PvTable *PVTable      // Hash table of best moves found per position, used to recover the PV line
	PvArray [MaxDepth]int // Principal variation moves, filled in by GetPvLine

	SearchHistory [13][120]int
	SearchKillers [2][MaxDepth]int
}

type SearchInfo struct {
	StartTime int64
	StopTime  int64
	Depth     int
	DepthSet  int
	TimeSet   bool
	MovesToGo int
	Infinite  bool

	Nodes int64

	Quit    bool
	Stopped bool

	Fh  float64
	Fhf float64
}

// MoveList holds a list of legal moves for a position, used during search.
type MoveList struct {
	Moves [MaxPositionMoves]Move
	Count int
}

// ---- Variables ----

// Sq120ToSq64 maps a 120-square mailbox index to its corresponding 64-square index.
// Off-board squares map to sentinel value offBoard (65).
var Sq120ToSq64 [120]int

// Sq64ToSq120 maps a 64-square index back to its corresponding 120-square mailbox index.
var Sq64ToSq120 [64]Square

// ---- Board Methods ----

// HisPly returns the count of all half-moves (plies) made since game start.
func (b *Board) HisPly() int {
	return len(b.History)
}

// ---- Square Conversion Functions ----

// FR2SQ converts a file and rank into a 120-square mailbox board index.
// The formula (21 + file + rank*10) places (0,0) at A1 = 21.
func FR2SQ(file File, rank Rank) Square {
	return Square(21 + int(file) + int(rank)*10)
}

// SQ64 converts a 120-square mailbox index to a 64-square board index.
func SQ64(sq120 Square) int {
	return Sq120ToSq64[sq120]
}

// SQ120 converts a 64-square board index back to a 120-square mailbox index.
func SQ120(sq64 int) Square {
	return Sq64ToSq120[sq64]
}

// ---- Bitboard Utility Functions ----

// POP removes and returns the index of the lowest set bit in the given bitboard.
// This is a convenience wrapper around PopBit.
func POP(bb *Bitboard) int {
	return PopBit(bb)
}

// CNT returns the number of set bits (population count) in a bitboard.
func CNT(b Bitboard) int {
	return CountBits(b)
}

// SETBIT sets the bit corresponding to square sq in the bitboard.
func SETBIT(bb *Bitboard, sq int) {
	*bb |= SetMask[sq]
}

// CLRBIT clears the bit corresponding to square sq in the bitboard.
func CLRBIT(bb *Bitboard, sq int) {
	*bb &= ClearMask[sq]
}

// ---- Debug / Assert ----

// Assert panics with the given message if condition is false, but only when
// Debug is true. Intended for development-time invariant checking.
func Assert(condition bool, message string) {
	if !Debug {
		return
	}
	if !condition {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}

// ---- Piece Type Queries ----

// IsBQ returns true if the given piece is a bishop or queen (i.e., a diagonal slider).
func IsBQ(p Piece) bool {
	return PieceBishopQueen[p]
}

// IsRQ returns true if the given piece is a rook or queen (i.e., an orthogonal slider).
func IsRQ(p Piece) bool {
	return PieceRookQueen[p]
}

// IsKn returns true if the given piece is a knight.
func IsKn(p Piece) bool {
	return PieceKnight[p]
}

// IsKi returns true if the given piece is a king.
func IsKi(p Piece) bool {
	return PieceKing[p]
}

// ---- Move Encoding / Decoding ----

/*
	0000 0000 0000 0000 0111 1111 -> From 0x7F
	0000 0000 0011 1111 1000 0000 -> To >> 7, 0x7F
	0000 0000 0011 1100 0000 0000 -> Captured >> 14, 0xF
	0000 0000 0100 0000 0000 0000 -> EP 0x40000
	0000 0000 1000 0000 0000 0000 -> Pawn Start 0x80000
	0000 1111 0000 0000 0000 0000 -> Promoted Piece >> 20, 0xF
	0001 0000 0000 0000 0000 0000 -> Castle 0x1000000
*/

// FromSq extracts the origin square (0–63) from an encoded move integer.
// Bits 0–6 hold the from-square index, masked with 0x7F.
func FromSq(m int) int {
	return m & 0x7F
}

// ToSq extracts the destination square (0–63) from an encoded move integer.
// Bits 7–13 (shifted right by 7) hold the to-square index, masked with 0x7F.
func ToSq(m int) int {
	return (m >> 7) & 0x7F
}

// Captured extracts the captured piece type from an encoded move integer.
// Bits 14–17 (shifted right by 14) hold the captured piece, masked with 0xF.
// Returns Empty if the move is not a capture.
func Captured(m int) Piece {
	return Piece((m >> 14) & 0xF)
}

// Promoted extracts the promotion piece type from an encoded move integer.
// Bits 20–23 (shifted right by 20) hold the promoted piece, masked with 0xF.
// Returns Empty if the move is not a promotion.
func Promoted(m int) Piece {
	return Piece((m >> 20) & 0xF)
}
