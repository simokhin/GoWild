package main

import "math/rand"

const offBoard = 65 // Sentinel value used to mark off-board squares in the 120-to-64 mapping

var SetMask [64]Bitboard
var ClearMask [64]Bitboard

var PieceKeys [13][120]Bitboard
var SideKey Bitboard
var CastleKeys [16]Bitboard

func AllInit() {
	InitSq120ToSq64()
	InitBitMasks()
	InitHashKeys()
}

func Rand64() Bitboard {
	return Bitboard(rand.Uint64())
}

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
