package engine

import (
	"fmt"
	"time"
)

var leafNodes int64

func Perft(depth int, pos *Board) {
	Assert(CheckBoard(pos), "board check failed")

	if depth == 0 {
		leafNodes++
		return
	}

	list := &MoveList{}
	GenerateAllMoves(pos, list)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		if !MakeMove(pos, list.Moves[moveNum].MoveInt) {
			continue
		}
		Perft(depth-1, pos)
		TakeMove(pos)
	}
}

func PerftTest(depth int, pos *Board) {
	Assert(CheckBoard(pos), "board check failed")

	PrintBoard(pos)
	fmt.Printf("\nStarting Test To Depth:%d\n", depth)

	leafNodes = 0
	start := time.Now()

	list := &MoveList{}
	GenerateAllMoves(pos, list)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		move := list.Moves[moveNum].MoveInt
		if !MakeMove(pos, move) {
			continue
		}

		cumNodes := leafNodes
		Perft(depth-1, pos)
		TakeMove(pos)

		oldNodes := leafNodes - cumNodes
		fmt.Printf("move %d : %s : %d\n", moveNum+1, PrMove(move), oldNodes)
	}

	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("\nTest Complete : %d nodes visited in %dms\n", leafNodes, elapsed)
}
