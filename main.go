package main

import "fmt"

// main is the program entry point. It initialises the board representation
// lookup tables. Further game logic will be added here.
func main() {
	AllInit()

	var playBitBoard Bitboard = 0

	fmt.Println("Start:")
	PrintBitBoard(playBitBoard)

	playBitBoard |= 1 << SQ64(D2)
	fmt.Println("D2 Added:")
	PrintBitBoard(playBitBoard)

	playBitBoard |= 1 << SQ64(H7)
	fmt.Println("H7 Added:")
	PrintBitBoard(playBitBoard)
}
