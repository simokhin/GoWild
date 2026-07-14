package engine

import (
	"fmt"
	"math/bits"
)

// CountBits returns the number of set bits (population count) in a bitboard.
// Uses the hardware-accelerated popcount instruction via math/bits.
func CountBits(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}

// PopBit removes and returns the index of the lowest set bit in a bitboard.
// After calling, bb will have that bit cleared via the classic x & (x-1) trick.
func PopBit(bb *Bitboard) int {
	index := bits.TrailingZeros64(uint64(*bb))
	*bb &= *bb - 1 // Clear the lowest set bit
	return index
}

// PrintBitBoard prints a visual representation of a bitboard to stdout.
// Set bits are shown as 'X', empty squares as '-', arranged in the standard
// chessboard orientation (rank 8 at top, rank 1 at bottom).
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
