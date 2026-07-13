package main

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
