package main

import "fmt"

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	var playBitBoard Bitboard = 0

	playBitBoard |= 1 << SQ64(D2)
	playBitBoard |= 1 << SQ64(D3)
	playBitBoard |= 1 << SQ64(D4)

	PrintBitBoard(playBitBoard)

	count := CNT(playBitBoard)

	fmt.Printf("Count: %d\n", count)

	index := POP(&playBitBoard)
	fmt.Printf("Index: %d\n", index)
	PrintBitBoard(playBitBoard)

	count = CNT(playBitBoard)
	fmt.Printf("Count: %d\n", count)
}
