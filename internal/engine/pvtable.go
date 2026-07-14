package engine

import (
	"unsafe"
)

// PvSize is the size in bytes of the PV table.
const PvSize = 0x100000 * 128

const IsMate = Mate - MaxDepth

// InitHashTable (re)allocates table's slots to fill PvSize and clears them.
func InitHashTable(table *HashTable) {
	numEntries := PvSize / int(unsafe.Sizeof(HashEntry{}))
	numEntries -= 2

	table.PTable = make([]HashEntry, numEntries)
	ClearHashTable(table)

	//fmt.Printf("HashTable init complete with %d entries\n", table.NumEntries())
}

// ClearHashTable resets every slot in table to its empty state.
func ClearHashTable(table *HashTable) {
	for i := range table.PTable {
		table.PTable[i].PosKey = 0
		table.PTable[i].Move = 0
		table.PTable[i].Depth = 0
		table.PTable[i].Score = 0
		table.PTable[i].Flags = 0
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

// StoreHashEntry records move as the best move found for pos's current
// position, keyed by its Zobrist hash. Must be called before move is made on
// the board, since it stores under the position the move is played from.
func StoreHashEntry(pos *Board, move, score, flags, depth int) {
	index := int(pos.PosKey % uint64(pos.HashTable.NumEntries()))

	Assert(index >= 0 && index <= pos.HashTable.NumEntries()-1, "pv table index out of range")
	Assert(depth >= 1 && depth < MaxDepth, "invalid depth")
	Assert(flags >= HFAlpha && flags <= HFExact, "invalid flags")
	Assert(score >= -Infinite && score <= Infinite, "score out of bounds")
	Assert(pos.Ply >= 0 && pos.Ply < MaxDepth, "invalid ply")

	if pos.HashTable.PTable[index].PosKey == 0 {
		pos.HashTable.NewWrite++
	} else {
		pos.HashTable.OverWrite++
	}

	if score > IsMate {
		score += pos.Ply
	} else if score < -IsMate {
		score -= pos.Ply
	}

	pos.HashTable.PTable[index].Move = move
	pos.HashTable.PTable[index].PosKey = pos.PosKey
	pos.HashTable.PTable[index].Flags = flags
	pos.HashTable.PTable[index].Score = score
	pos.HashTable.PTable[index].Depth = depth
}

// ProbePvMove looks up the best move stored for pos's current position,
// regardless of the stored depth or bound type. Used only to reconstruct the
// PV line after a search, where depth/alpha/beta gating is not relevant.
// Returns NoMove if no entry is found, or if the slot's key belongs to a
// different position (hash collision).
func ProbePvMove(pos *Board) int {
	index := int(pos.PosKey % uint64(pos.HashTable.NumEntries()))
	Assert(index >= 0 && index <= pos.HashTable.NumEntries()-1, "pv table index out of range")

	if pos.HashTable.PTable[index].PosKey == pos.PosKey {
		return pos.HashTable.PTable[index].Move
	}
	return NoMove
}

// ProbeHashEntry looks up the best move stored for pos's current position. It
// returns NoMove if no entry is found, or if the slot's key belongs to a
// different position (hash collision).
func ProbeHashEntry(pos *Board, alpha, beta, depth int) (move int, score int, found bool) {
	index := int(pos.PosKey % uint64(pos.HashTable.NumEntries()))
	Assert(index >= 0 && index <= pos.HashTable.NumEntries()-1, "pv table index out of range")
	Assert(depth >= 1 && depth < MaxDepth, "invalid depth")
	Assert(alpha < beta, "alpha must be less than beta")
	Assert(alpha >= -Infinite && alpha <= Infinite, "alpha out of bounds")
	Assert(beta >= -Infinite && beta <= Infinite, "beta out of bounds")
	Assert(pos.Ply >= 0 && pos.Ply < MaxDepth, "invalid ply")

	if pos.HashTable.PTable[index].PosKey == pos.PosKey {
		move = pos.HashTable.PTable[index].Move

		if pos.HashTable.PTable[index].Depth >= depth {
			pos.HashTable.Hit++

			score = pos.HashTable.PTable[index].Score

			if score > IsMate {
				score -= pos.Ply
			} else if score < -IsMate {
				score += pos.Ply
			}

			switch pos.HashTable.PTable[index].Flags {
			case HFAlpha:
				if score <= alpha {
					score = alpha
					return move, score, true
				}
			case HFBeta:
				if score >= beta {
					score = beta
					return move, score, true
				}
			case HFExact:
				return move, score, true
			}
		}
	}

	return move, score, false
}
