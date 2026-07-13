package main

import "fmt"

const Infinite = 30000
const Mate = 29000

func PickNextMove(moveNum int, list *MoveList) {
	bestScore := 0
	bestNum := moveNum

	for index := moveNum; index < list.Count; index++ {
		if list.Moves[index].Score > bestScore {
			bestScore = list.Moves[index].Score
			bestNum = index
		}
	}

	list.Moves[moveNum], list.Moves[bestNum] = list.Moves[bestNum], list.Moves[moveNum]
}

// IsRepetition reports whether the current position has occurred before
// since the last irreversible move (capture, pawn move, or loss of castling
// rights), which resets the fifty-move counter. Only that window of history
// can contain a repeated position.
func IsRepetition(pos *Board) bool {
	for index := pos.HisPly() - pos.FiftyMove; index < pos.HisPly(); index++ {
		Assert(index >= 0 && index < MaxGameMoves, "history index out of range")
		if pos.PosKey == pos.History[index].PosKey {
			return true
		}
	}
	return false
}

func ClearForSearch(pos *Board, info *SearchInfo) {
	for index := range 13 {
		for index2 := range 120 {
			pos.SearchHistory[index][index2] = 0
		}
	}

	for index := range 2 {
		for index2 := range MaxDepth {
			pos.SearchKillers[index][index2] = 0
		}
	}

	ClearPvTable(pos.PvTable)

	pos.Ply = 0

	info.StartTime = GetTimeMs()
	info.Stopped = false
	info.Nodes = 0

	info.Fh = 0
	info.Fhf = 0
}

func SearchPosition(pos *Board, info *SearchInfo) {
	bestMove := NoMove
	bestScore := -Infinite
	var pvMoves int

	ClearForSearch(pos, info)

	for currentDepth := 1; currentDepth <= info.Depth; currentDepth++ {
		bestScore = AlphaBeta(-Infinite, Infinite, currentDepth, pos, info, true)

		pvMoves = GetPvLine(currentDepth, pos)
		bestMove = pos.PvArray[0]

		fmt.Printf("Depth:%d score:%d move:%s nodes:%d ", currentDepth, bestScore, PrMove(bestMove), info.Nodes)

		pvMoves = GetPvLine(currentDepth, pos)
		fmt.Print("pv")
		for pvNum := 0; pvNum < pvMoves; pvNum++ {
			fmt.Printf(" %s", PrMove(pos.PvArray[pvNum]))
		}
		fmt.Println()
		fmt.Printf("Ordering:%.2f\n", info.Fhf/info.Fh)
	}
}

func AlphaBeta(alpha, beta, depth int, pos *Board, info *SearchInfo, doNull bool) int {
	Assert(CheckBoard(pos), "board check failed")

	if depth == 0 {
		info.Nodes++
		return EvalPosition(pos)
	}

	if IsRepetition(pos) || pos.FiftyMove >= 100 {
		return 0
	}

	if pos.Ply > MaxDepth-1 {
		return EvalPosition(pos)
	}

	info.Nodes++

	list := &MoveList{}
	GenerateAllMoves(pos, list)

	legal := 0
	oldAlpha := alpha
	bestMove := NoMove
	score := -Infinite

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		PickNextMove(moveNum, list)

		if !MakeMove(pos, list.Moves[moveNum].MoveInt) {
			continue
		}

		legal++
		score = -AlphaBeta(-beta, -alpha, depth-1, pos, info, true)
		TakeMove(pos)

		if score > alpha {
			if score >= beta {
				if legal == 1 {
					info.Fhf++
				}
				info.Fh++
				return beta
			}
			alpha = score
			bestMove = list.Moves[moveNum].MoveInt
		}
	}

	if legal == 0 {
		if SqAttacked(pos.KingSq[pos.Side], pos.Side^1, pos) {
			return -Mate + pos.Ply
		}
		return 0
	}

	if alpha != oldAlpha {
		StorePvMove(pos, bestMove)
	}

	return alpha
}
