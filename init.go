package main

import (
	"math/rand"
)

// offBoard is the sentinel value used to mark off-board squares in the
// 120-to-64 mapping table (Sq120ToSq64).
const offBoard = 65

// SetMask contains precomputed bitmasks with a single bit set at each 64-square index.
var SetMask [64]Bitboard

// ClearMask contains precomputed bitmasks with all bits set except one at each index.
var ClearMask [64]Bitboard

// PieceKeys holds Zobrist random numbers for each piece type at each 120-square index.
var PieceKeys [13][120]Bitboard

// SideKey is the Zobrist random number XORed into the hash when it is black's turn to move.
var SideKey Bitboard

// CastleKeys holds Zobrist random numbers for each of the 16 possible castling-rights states.
var CastleKeys [16]Bitboard

// FilesBrd maps each 120-square index to its file (FileA–FileH), or FileNone for off-board squares.
var FilesBrd [120]File

// RanksBrd maps each 120-square index to its rank (Rank1–Rank8), or RankNone for off-board squares.
var RanksBrd [120]Rank

var FileBBMask [8]Bitboard
var RankBBMask [8]Bitboard

// AllInit initialises all lookup tables required by the engine: square mapping,
// bit masks, and Zobrist hash keys. Must be called once at startup.
func AllInit() {
	InitSq120ToSq64()
	InitBitMasks()
	InitHashKeys()
	InitFileRankBrd()
	InitEvalMasks()
	InitMvvLva()
}

func InitEvalMasks() {
	for sq := 0; sq < 8; sq++ {
		FileBBMask[sq] = 0
		RankBBMask[sq] = 0
	}

	for r := Rank8; r >= Rank1; r-- {
		for f := FileA; f <= FileH; f++ {
			sq := int(r)*8 + int(f)
			FileBBMask[f] |= 1 << Bitboard(sq)
			RankBBMask[r] |= 1 << Bitboard(sq)
		}
	}

	for r := Rank8; r >= Rank1; r-- {
		PrintBitBoard(RankBBMask[r])
	}

	for f := FileA; f <= FileH; f++ {
		PrintBitBoard(FileBBMask[f])
	}
}

// InitFileRankBrd precomputes the FilesBrd and RanksBrd lookup tables.
// For every 120-square index that corresponds to a real board square, the
// table stores its file and rank; all other entries are set to FileNone/RankNone.
func InitFileRankBrd() {
	for index := range 120 {
		FilesBrd[index] = FileNone
		RanksBrd[index] = RankNone
	}

	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FR2SQ(file, rank)
			FilesBrd[sq] = file
			RanksBrd[sq] = rank
		}
	}
}

// Rand64 generates a random 64-bit value for use as a Zobrist key component.
func Rand64() Bitboard {
	return Bitboard(rand.Uint64())
}

// InitHashKeys populates the Zobrist random number tables: piece-square keys,
// side-to-move key, and castling rights keys.
func InitHashKeys() {
	for index := range 13 {
		for index2 := range 120 {
			PieceKeys[index][index2] = Rand64()
		}
	}
	SideKey = Rand64()
	for index := range 16 {
		CastleKeys[index] = Rand64()
	}
}

// InitBitMasks precomputes the SetMask and ClearMask tables for fast bitboard
// bit manipulation. SetMask[i] has only bit i set; ClearMask[i] has all but bit i set.
func InitBitMasks() {
	for index := range 64 {
		SetMask[index] = 0
		ClearMask[index] = 0
	}

	for index := range 64 {
		SetMask[index] |= 1 << Bitboard(index)
		ClearMask[index] = ^SetMask[index]
	}
}

// InitSq120ToSq64 builds the lookup tables that translate between the 120-square
// mailbox board indices and the compact 64-square array indices.
func InitSq120ToSq64() {
	for index := range 120 {
		Sq120ToSq64[index] = offBoard
	}
	for index := range 64 {
		Sq64ToSq120[index] = 120
	}

	sq64 := 0
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FR2SQ(file, rank)
			Sq64ToSq120[sq64] = sq
			Sq120ToSq64[sq] = sq64
			sq64++
		}
	}
}
