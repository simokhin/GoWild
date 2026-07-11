package main

import (
	"fmt"
	"math/bits"
)

// Counts the number of set bits in a bitboard.
func CountBits(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

// PopBit removes and returns index of lowest set bit in a bitboard.
// After calling, bb will have that bit cleared.
func PopBit(bb *Bitboard) int {
	index := bits.TrailingZeros64(uint64(*bb))
	*bb &= *bb - 1 // снимаем этот бит
	return index
}

// PrintBitBoard prints the board representation with X for set bits, - for empty squares.
func PrintBitBoard(bb Bitboard) {
	var shiftMe Bitboard = 1

	fmt.Println()
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
			sq := FR2SQ(file, rank) // 120 based
			sq64 := SQ64(sq)        // 64 based

			if (shiftMe<<sq64)&bb != 0 { // Check if bit is set
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
