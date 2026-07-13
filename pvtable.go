package main

import (
	"fmt"
	"unsafe"
)

// PvSize is the size in bytes of the PV table.
const PvSize = 0x100000 * 2

// InitPvTable (re)allocates table's slots to fill PvSize and clears them.
func InitPvTable(table *PVTable) {
	numEntries := PvSize / int(unsafe.Sizeof(PVEntry{}))
	numEntries -= 2

	table.PTable = make([]PVEntry, numEntries)
	ClearPvTable(table)

	fmt.Printf("PvTable init complete with %d entries\n", table.NumEntries())
}

// ClearPvTable resets every slot in table to its empty state.
func ClearPvTable(table *PVTable) {
	for i := range table.PTable {
		table.PTable[i].PosKey = 0
		table.PTable[i].Move = 0
	}
}

// GetPvLine walks the PV table from pos's current position, following the
// stored best move at each step, making it on the board, and recording it in
// pos.PvArray. It stops after depth moves, when a stored move no longer
// exists in the current position (a hash collision or stale entry), or when
// no move is stored. All moves made while walking are unwound before
// returning, so pos is left unchanged. Returns the number of moves found.
func GetPvLine(depth int, pos *Board) int {
	Assert(depth < MaxDepth && depth >= 1, "invalid depth")

	move := ProbePvMove(pos)
	count := 0

	for move != NoMove && count < depth {
		Assert(count < MaxDepth, "count out of range")

		if MoveExist(pos, move) {
			MakeMove(pos, move)
			pos.PvArray[count] = move
			count++
		} else {
			break
		}
		move = ProbePvMove(pos)
	}

	for pos.Ply > 0 {
		TakeMove(pos)
	}

	return count
}

// StorePvMove records move as the best move found for pos's current
// position, keyed by its Zobrist hash. Must be called before move is made on
// the board, since it stores under the position the move is played from.
func StorePvMove(pos *Board, move int) {
	index := int(pos.PosKey % uint64(pos.PvTable.NumEntries()))
	Assert(index >= 0 && index <= pos.PvTable.NumEntries()-1, "pv table index out of range")

	pos.PvTable.PTable[index].Move = move
	pos.PvTable.PTable[index].PosKey = pos.PosKey
}

// ProbePvMove looks up the best move stored for pos's current position. It
// returns NoMove if no entry is found, or if the slot's key belongs to a
// different position (hash collision).
func ProbePvMove(pos *Board) int {
	index := int(pos.PosKey % uint64(pos.PvTable.NumEntries()))
	Assert(index >= 0 && index <= pos.PvTable.NumEntries()-1, "pv table index out of range")

	if pos.PvTable.PTable[index].PosKey == pos.PosKey {
		return pos.PvTable.PTable[index].Move
	}
	return NoMove
}
