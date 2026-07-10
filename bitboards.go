package main

import "fmt"

func PrintBitBoard(bb Bitboard) {
	var shiftMe Bitboard = 1

	fmt.Println()
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
			sq := FR2SQ(file, rank) // 120 based
			sq64 := SQ64(sq)        // 64 based

			if (shiftMe<<sq64)&bb != 0 {
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
